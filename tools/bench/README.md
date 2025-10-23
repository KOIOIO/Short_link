本目录包含一个简单的 gRPC 基准脚本 `bench_rpc.go`，用于比较两个 RPC:

- `GenerateShortUrl`：使用 Redis（服务中实现）
- `FilterByMyBloomFilter`：使用项目内存布隆实现

构建与运行：

1. 进入仓库根目录（确保 `go.mod` 在项目根）：

```powershell
cd d:\shorturl\short-url
```

2. 编译并运行基准：

```powershell
go run tools/bench/bench_rpc.go -addr 127.0.0.1:50051 -concurrency 50 -requests 1000 -url https://www.example.com/long/path
```

说明：
- `-addr`：gRPC 服务地址（例：127.0.0.1:50051）
- `-concurrency`：并发 worker 数量
- `-requests`：总请求次数
- `-url`：要缩短的长 URL

输出包含每个 RPC 的延迟分位数（p50/p90/p95）、吞吐（ops/sec）和成功/错误统计。可以调整 `-concurrency` 和 `-requests` 来模拟不同负载。
