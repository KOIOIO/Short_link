package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "example.com/shorturl/short-url/zero_remake/shorturl_rpc/types/shortUrl"
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

func worker(ctx context.Context, wg *sync.WaitGroup, client pb.ShortUrlClient, inUrl string, jobs <-chan struct{}, results chan<- time.Duration, errs chan<- error, rpcType string) {
	defer wg.Done()
	for range jobs {
		start := time.Now()
		var err error
		ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
		switch rpcType {
		case "GenerateShortUrl":
			_, err = client.GenerateShortUrl(ctx2, &pb.GenerateShortUrlRequest{Url: inUrl, Expiration: "0"})
		case "FilterByMyBloomFilter":
			_, err = client.FilterByMyBloomFilter(ctx2, &pb.GenerateShortUrlRequest{Url: inUrl, Expiration: "0"})
		default:
			err = fmt.Errorf("unknown rpc type")
		}
		cancel()
		lat := time.Since(start)
		if err != nil {
			errs <- err
		} else {
			results <- lat
		}
	}
}

func runBenchmark(addr, rpcType, url string, concurrency, requests int) (*stat, time.Duration) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to dial %s: %v\n", addr, err)
		os.Exit(1)
	}
	defer conn.Close()
	client := pb.NewShortUrlClient(conn)

	jobs := make(chan struct{}, requests)
	results := make(chan time.Duration, requests)
	errs := make(chan error, requests)
	var wg sync.WaitGroup

	ctx := context.Background()
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker(ctx, &wg, client, url, jobs, results, errs, rpcType)
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
	addr := flag.String("addr", "127.0.0.1:50051", "gRPC server address")
	concurrency := flag.Int("concurrency", 50, "number of concurrent workers")
	requests := flag.Int("requests", 1000, "total requests to send")
	inUrl := flag.String("url", "https://www.example.com/hello", "long URL to shorten")
	flag.Parse()

	fmt.Printf("addr=%s concurrency=%d requests=%d url=%s\n", *addr, *concurrency, *requests, *inUrl)

	// Run GenerateShortUrl
	s1Start := time.Now()
	s1, _ := runBenchmark(*addr, "GenerateShortUrl", *inUrl, *concurrency, *requests)
	s1ElapsedTotal := time.Since(s1Start)
	s1.print("GenerateShortUrl", s1ElapsedTotal)

	// Run FilterByMyBloomFilter
	s2Start := time.Now()
	s2, _ := runBenchmark(*addr, "FilterByMyBloomFilter", *inUrl, *concurrency, *requests)
	s2ElapsedTotal := time.Since(s2Start)
	s2.print("FilterByMyBloomFilter", s2ElapsedTotal)

	fmt.Println("Done")
}
