# 评测用镜像：保留完整 Go 工具链，并在构建阶段预下载依赖。
FROM golang:1.26

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build ./...

CMD ["bash"]
