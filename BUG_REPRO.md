# 修复前故障复现（Docker）

## 项目与标准命令

项目使用 Go 1.26，在 Apple Silicon 主机上通过 Docker Desktop 的 linux/arm64 平台验证。构建命令为：

```bash
docker build -f benzhi.Dockerfile -t go-courier-slot-coordinator__005-bug:20260817 .
```

## 环境构建与编译

使用上述命令构建镜像成功。进入容器后执行：

```bash
go version
go build ./...
```

输出的 Go 版本为 `go version go1.26.6 linux/arm64`，编译命令退出码为 0。

## 故障触发步骤

在同一容器工作目录执行：

```bash
go test ./...
```

## 实际错误输出

```text
?   	github.com/1260124186-cc/courier-slot-coordinator/cmd/server	[no test files]
--- FAIL: TestCourierConflictKeepsHTTPConflictStatus (0.00s)
    handler_test.go:65: status = 500, want 409; body={"error":"assign courier: shipment changed while being updated"}
FAIL
FAIL	github.com/1260124186-cc/courier-slot-coordinator/internal/api	0.004s
ok  	github.com/1260124186-cc/courier-slot-coordinator/internal/domain	0.003s
--- FAIL: TestMemoryStoreUsesOptimisticVersioning (0.00s)
    memory_test.go:35: stale Update() error = update shipment conflict: shipment changed while being updated, want ErrConflict
FAIL
FAIL	github.com/1260124186-cc/courier-slot-coordinator/internal/repository	0.002s
ok  	github.com/1260124186-cc/courier-slot-coordinator/internal/service	0.002s
FAIL
```

## 期望行为

两位调度员并发更新同一张配送单时，后提交的一方应收到明确的冲突响应，前端可以提示重新刷新后再分配，而不应显示为服务器内部异常。
