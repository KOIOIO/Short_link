package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type stat struct {
	min       time.Duration
	tmax      time.Duration
	sum       time.Duration
	count     int64
	errCount  int64
	latencies []time.Duration
}

func (s *stat) add(d time.Duration) {
	if s.count == 0 || d < s.min {
		s.min = d
	}
	if d > s.tmax {
		s.tmax = d
	}
	s.sum += d
	s.count++
	s.latencies = append(s.latencies, d)
}

func (s *stat) print(name string, elapsed time.Duration) {
	if s.count == 0 {
		fmt.Printf("%s: no successful requests, errors=%d\n", name, s.errCount)
		return
	}
	slice := s.latencies
	sort.Slice(slice, func(i, j int) bool { return slice[i] < slice[j] })
	p50 := slice[len(slice)*50/100]
	p90 := slice[len(slice)*90/100]
	p95 := slice[len(slice)*95/100]
	throughput := float64(s.count) / elapsed.Seconds()
	fmt.Printf("=== %s results ===\n", name)
	fmt.Printf("requests: %d, errors: %d\n", s.count, s.errCount)
	fmt.Printf("min: %v, max: %v, avg: %v\n", s.min, s.tmax, time.Duration(int64(s.sum)/s.count))
	fmt.Printf("p50: %v, p90: %v, p95: %v\n", p50, p90, p95)
	fmt.Printf("throughput(ops/sec): %.2f\n\n", throughput)
}

var uniqueFlag bool
var reqCounter uint64

func workerHTTP(ctx context.Context, wg *sync.WaitGroup, client *http.Client, base, path, urlStr string, jobs <-chan struct{}, results chan<- time.Duration, errs chan<- error) {
	defer wg.Done()
	for range jobs {
		start := time.Now()
		// build form data; optionally add uid to url to make it unique per request
		reqUrl := urlStr
		if uniqueFlag {
			id := atomic.AddUint64(&reqCounter, 1)
			if strings.Contains(reqUrl, "?") {
				reqUrl = reqUrl + "&uid=" + strconv.FormatUint(id, 10)
			} else {
				reqUrl = reqUrl + "?uid=" + strconv.FormatUint(id, 10)
			}
		}

		form := url.Values{}
		form.Set("url", reqUrl)
		form.Set("expiration", "0")
		body := bytes.NewBufferString(form.Encode())
		req, _ := http.NewRequestWithContext(ctx, "POST", strings.TrimRight(base, "/")+path, body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		lat := time.Since(start)
		if err != nil {
			errs <- err
			continue
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			results <- lat
		} else {
			errs <- fmt.Errorf("status %d", resp.StatusCode)
		}
	}
}

func runHTTPBenchmark(base, path, urlStr string, concurrency, requests int) (*stat, time.Duration) {
	client := &http.Client{Timeout: 10 * time.Second}
	jobs := make(chan struct{}, requests)
	results := make(chan time.Duration, requests)
	errs := make(chan error, requests)
	var wg sync.WaitGroup

	ctx := context.Background()
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go workerHTTP(ctx, &wg, client, base, path, urlStr, jobs, results, errs)
	}

	startAll := time.Now()
	for i := 0; i < requests; i++ {
		jobs <- struct{}{}
	}
	close(jobs)

	wg.Wait()
	elapsed := time.Since(startAll)

	close(results)
	close(errs)

	s := &stat{min: time.Duration(1<<63 - 1)}
	for lat := range results {
		s.add(lat)
	}
	// count errors
	errCount := 0
	for e := range errs {
		if e != nil {
			errCount++
		}
	}
	s.errCount = int64(errCount)
	return s, elapsed
}

func main() {
	base := flag.String("base", "http://127.0.0.1:8080", "API base URL, e.g. http://127.0.0.1:8080")
	path := flag.String("path", "/generate", "API path to test, e.g. /generate or /filterbymybloomfilter")
	concurrency := flag.Int("concurrency", 50, "number of concurrent workers")
	requests := flag.Int("requests", 1000, "total requests to send")
	urlStr := flag.String("url", "https://www.example.com/long/path", "long URL to shorten")
	flag.BoolVar(&uniqueFlag, "unique", false, "append uid to url to make requests unique (avoid bloom cache hit)")
	flag.Parse()

	fmt.Printf("base=%s path=%s concurrency=%d requests=%d url=%s unique=%v\n", *base, *path, *concurrency, *requests, *urlStr, uniqueFlag)

	// Run first path
	s1Start := time.Now()
	s1, _ := runHTTPBenchmark(*base, *path, *urlStr, *concurrency, *requests)
	s1ElapsedTotal := time.Since(s1Start)
	s1.print(*path, s1ElapsedTotal)

	fmt.Println("Done")
}
