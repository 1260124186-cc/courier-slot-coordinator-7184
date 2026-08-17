# 修复前故障复现（Docker）

## 项目与标准命令

项目使用 Go 1.26，在 Apple Silicon 主机上通过 Docker Desktop 的 linux/arm64 平台验证。构建命令为：

```bash
docker build -f benzhi.Dockerfile -t go-courier-slot-coordinator__003-bug:20260817 .
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
--- FAIL: TestDeliveredShipmentReleasesZoneCapacity (0.00s)
    dispatch_test.go:87: CreateShipment() after delivery error = zone package capacity exceeded, want capacity to be released
FAIL
FAIL	github.com/1260124186-cc/courier-slot-coordinator/internal/service	0.002s
FAIL
```

## 期望行为

西区的大件配送单完成送达后不应继续占用在途配送容量；新建同等大小的配送单应能创建成功，避免已经完成的任务长期阻塞后续排班。
