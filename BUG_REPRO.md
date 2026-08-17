# 修复前故障复现（Docker）

## 项目与标准命令

项目使用 Go 1.26，在 Apple Silicon 主机上通过 Docker Desktop 的 linux/arm64 平台验证。构建命令为：

```bash
docker build -f benzhi.Dockerfile -t go-courier-slot-coordinator__002-bug:20260817 .
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
--- FAIL: TestShipmentValidationAndClone (0.00s)
    shipment_test.go:21: CloneShipment() did not isolate package slice
FAIL
FAIL	github.com/1260124186-cc/courier-slot-coordinator/internal/domain	0.001s
ok  	github.com/1260124186-cc/courier-slot-coordinator/internal/repository	0.001s
--- FAIL: TestReadingShipmentCannotChangeZoneCapacity (0.00s)
    snapshot_test.go:37: CreateShipment() error = <nil>, want ErrZoneCapacity after a caller mutates a read result
FAIL
FAIL	github.com/1260124186-cc/courier-slot-coordinator/internal/service	0.001s
FAIL
```

## 期望行为

读取配送单或查看排班汇总只应返回数据副本，不能改变已经保存的包裹数量；读取后再次创建大件配送单时，西区的容量限制仍应拒绝超出额度的请求。
