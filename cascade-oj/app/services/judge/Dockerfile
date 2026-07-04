FROM golang:1.25 AS builder

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

# Wait-for-it script: poll RabbitMQ until its port is reachable,
# then launch judge.  This avoids the race between docker-compose
# depends_on (which only waits for the container, not the service)
# and the judge process connecting to RabbitMQ.
RUN echo '#!/bin/bash\n\
set -e\n\
RABBITMQ_HOST="${RABBITMQ_HOST:-cascade_rabbitmq}"\n\
RABBITMQ_PORT="${RABBITMQ_PORT:-5672}"\n\
MAX_RETRIES=30\n\
RETRY_INTERVAL=2\n\
echo "Waiting for RabbitMQ at ${RABBITMQ_HOST}:${RABBITMQ_PORT}..."\n\
for i in $(seq 1 $MAX_RETRIES); do\n\
  if timeout 2 bash -c "echo >/dev/tcp/${RABBITMQ_HOST}/${RABBITMQ_PORT}" 2>/dev/null; then\n\
    echo "RabbitMQ is reachable. Starting judge..."\n\
    exec ./judge -conf /data/conf\n\
  fi\n\
  echo "  attempt $i/$MAX_RETRIES: RabbitMQ not ready, retrying in ${RETRY_INTERVAL}s..."\n\
  sleep $RETRY_INTERVAL\n\
done\n\
echo "ERROR: RabbitMQ did not become available after $((MAX_RETRIES * RETRY_INTERVAL))s"\n\
exit 1\n\
' > /usr/local/bin/entrypoint.sh && \
    chmod +x /usr/local/bin/entrypoint.sh

WORKDIR /app

EXPOSE 8000
EXPOSE 9000
VOLUME /data/conf

# 解决 judge 启动后与 mq 连接失败的问题
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
