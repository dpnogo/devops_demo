FROM golang:1.20-alpine AS builder

ENV GO111MODULE on
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && apk add git


WORKDIR /build/apps

# 将 go mod 和 go sum 拷贝到 /build/apps 下 将本地文件添加到容器中
COPY . .

# 列出当前目录中的文件和目录，并显示它们的详细信息，包括大小、权限等
RUN ls -alsh

# 进行编译 --> 放到 /build/apps/.keep-server 以及存在 /build/apps/xxx
# 注意: go mod download 因为未指定参数，所以需要在和 go.mod 同层级进行执行 docker build ,
# 不然找不到, 或者使用 go mod download

RUN go mod download && go build -o .keep-server cmd/apiserver/apiserver.go


FROM alpine:3.16.3

# 设置时区
ENV TZ Asia/Shanghai

# 创建工作目录
WORKDIR /app

# 从第一阶段的构建结果中复制文件到当前阶段 拷贝编译程序
COPY --from=builder /build/apps/api-server /app/dist/api-server
COPY --from=builder /build/apps/config /app/dist/config

# 根据配置进行执行
CMD ["/app/dist/api-server", "-c", "/app/dist/config/apiserver.yaml"]