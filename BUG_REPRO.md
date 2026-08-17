# 修复前故障复现（Docker）

## 项目与标准命令

项目使用 Go 1.26，在 Apple Silicon 主机上通过 Docker Desktop 的 linux/arm64 平台验证。构建命令为：

```bash
docker build -f benzhi.Dockerfile -t go-courier-slot-coordinator__001-bug:20260817 .
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
ok  	github.com/1260124186-cc/courier-slot-coordinator/internal/api	0.002s
ok  	github.com/1260124186-cc/courier-slot-coordinator/internal/domain	0.001s
ok  	github.com/1260124186-cc/courier-slot-coordinator/internal/repository	0.001s
--- FAIL: TestConcurrentCreationDoesNotExceedZoneCapacity (0.00s)
    concurrency_test.go:57: concurrent creates succeeded 2 times with 0 capacity failures, want 1 and 1
FAIL
FAIL	github.com/1260124186-cc/courier-slot-coordinator/internal/service	0.001s
FAIL
```

## 期望行为

同一配送区在容量只剩一张大件配送单额度时，同时提交的两张大件配送单中只能成功创建一张，另一张应被拒绝，避免后续骑手排班超出可处理范围。
