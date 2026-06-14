FROM golang:1.24 AS builder

WORKDIR /src

# 先复制 go.mod 和 go.sum（位于仓库根），并下载模块依赖
COPY go.mod go.sum ./
RUN GOPROXY=https://goproxy.cn go mod download

ARG VERSION

# 复制服务和共享目录
COPY app/services/judge ./app/services/judge
COPY pkg ./pkg
COPY ent ./ent
COPY api ./api

WORKDIR /src/app/services/judge

# 使用 Makefile 构建二进制（和开发机器上的流程保持一致）
RUN GOPROXY=https://goproxy.cn make build VERSION_FLAG=${VERSION}

FROM debian:stable-slim

# Debian apt-get 镜像源
RUN sed -i 's@deb.debian.org@mirrors.ustc.edu.cn@g' /etc/apt/sources.list.d/debian.sources

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/app/services/judge/bin /app

WORKDIR /app

EXPOSE 8000
EXPOSE 9000
VOLUME /data/conf

CMD ["./judge", "-conf", "/data/conf"]
