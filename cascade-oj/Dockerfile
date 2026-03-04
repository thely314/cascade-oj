# Use the official Go-Judge image as the base image
FROM criyle/go-judge:v1.11.0
FROM debian:latest

# Debian apt-get 镜像源
RUN sed -i 's@deb.debian.org@mirrors.ustc.edu.cn@g' /etc/apt/sources.list.d/debian.sources

# 安装常用编译器
RUN apt-get update && apt-get install -y --no-install-recommends \
        gcc g++ \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

# 复制原始配置
WORKDIR /opt
COPY --from=criyle/go-judge:v1.10.0 /opt/go-judge /opt/mount.yaml /opt/
EXPOSE 5050/tcp 5051/tcp 5052/tcp
ENTRYPOINT ["./go-judge"]