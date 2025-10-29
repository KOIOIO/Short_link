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
go run tools/bench/bench_rpc.go -addr 127.0.0.1:8888 -concurrency 50 -requests 1000 -url https://www.example.com/long/path
```

说明：
- `-addr`：gRPC 服务地址（例：127.0.0.1:50051）
- `-concurrency`：并发 worker 数量
- `-requests`：总请求次数
- `-url`：要缩短的长 URL

输出包含每个 RPC 的延迟分位数（p50/p90/p95）、吞吐（ops/sec）和成功/错误统计。可以调整 `-concurrency` 和 `-requests` 来模拟不同负载。

HTTP 基准脚本（模拟前端/HTTP 请求）

另外仓库中提供了一个 HTTP 层的基准脚本 `bench_http.go`，用于直接对 API HTTP 接口做并发压测（更贴近前端行为），可以用来比较 `/generate` 与 `/filterbymybloomfilter` 在 HTTP 层的性能差异。

构建与运行（示例，PowerShell）:

```powershell
# 在仓库根目录运行（确保 go.mod 在项目根）
cd d:\shorturl\short-url

# 测试 /generate
go run tools/bench/bench_http.go -base http://127.0.0.1:8080 -path /generate -concurrency 50 -requests 1000 -url https://www.example.com/long/path

# 测试 /filterbymybloomfilter
go run tools/bench/bench_http.go -base http://127.0.0.1:8080 -path /filterbymybloomfilter -concurrency 50 -requests 1000 -url https://www.example.com/long/path
```

选项说明：
- `-base`：API 服务基地址，含协议与端口（例如 `http://127.0.0.1:8080`）
- `-path`：要测试的路径（例如 `/generate` 或 `/filterbymybloomfilter`）
- `-concurrency`：并发 worker 数量
- `-requests`：总请求数量
- `-url`：要缩短的原始 URL
- `-unique`：可选，传入该标志会在每次请求的 url 后追加 `?uid=N`，用于强制每次请求为不同 URL（避免布隆命中），便于区分“布隆拦截”与“写入逻辑问题”。

输出说明：脚本会打印延迟分位数（p50/p90/p95）、吞吐（ops/sec）与成功/错误计数。若你发现在相同 URL 下仍持续写入数据库（即布隆没有拦截），请参考 README 中的排查步骤：确认 API 路径使用的是同一布隆实现、检查 RedisBloom 是否初始化并检查 Redis 中 BF 存在性。
