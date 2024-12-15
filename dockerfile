# 使用 Golang 官方镜像作为基础镜像
FROM golang

ENV GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载 Go 依赖
RUN go mod download

# 复制项目源代码到容器内
COPY . .

# 暴露应用端口，假设 Gin 默认运行在 8080 端口
EXPOSE 8080

# 设置容器启动命令，使用 `go run` 运行主程序
CMD ["go", "run", "main.go"]
