# courier-slot-coordinator-7184 Docker 交付说明

## 项目概览
- Courier Slot Coordinator is a small Go HTTP service for local delivery teams. Dispatchers create a shipment for a delivery zone and time window, assign a courier, record collection
- Go module: `github.com/1260124186-cc/courier-slot-coordinator`

## 标准命令

```bash
go build ./...
go test ./...
```

## 实际启动入口

```bash
go run ./cmd/server
```

## Docker 构建

```bash
./build_benzhi_docker.sh courier-slot-coordinator-7184-benzhi linux/amd64
docker run --rm -it courier-slot-coordinator-7184-benzhi bash
```

## 环境

- 基础镜像: `golang:1.26`
- 依赖在镜像构建阶段预下载，容器内可直接执行 Go 构建和测试命令。
- 代码目录: `/app`
- 源码中检测到的服务端口: `8080`
