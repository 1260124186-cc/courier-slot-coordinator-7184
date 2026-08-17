# 修复前故障复现（Docker）

## 项目与标准命令

项目使用 Go 1.26，在 Apple Silicon 主机上通过 Docker Desktop 的 linux/arm64 平台验证。构建命令为：

```bash
docker build -f benzhi.Dockerfile -t go-courier-slot-coordinator__004-bug:20260817 .
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
--- FAIL: TestUnknownZoneSummaryEndpointIsRejected (0.00s)
    handler_test.go:34: status = 200, want 404; body={"zone":"central","ready":0,"collected":0,"delivered":0,"cancelled":0,"package_qty":0}
FAIL
FAIL	github.com/1260124186-cc/courier-slot-coordinator/internal/api	0.002s
ok  	github.com/1260124186-cc/courier-slot-coordinator/internal/domain	0.002s
ok  	github.com/1260124186-cc/courier-slot-coordinator/internal/repository	0.001s
--- FAIL: TestZoneSummaryRejectsUnconfiguredZone (0.00s)
    dispatch_test.go:66: ZoneSummary() error = <nil>, want ErrUnknownZone
FAIL
FAIL	github.com/1260124186-cc/courier-slot-coordinator/internal/service	0.001s
FAIL
```

## 期望行为

查询没有配置的配送区时，服务应明确拒绝请求，而不是返回看似正常但内容为空的汇总数据，避免调度员将不存在的配送区当作可接单区域。
