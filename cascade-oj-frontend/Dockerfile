FROM nginx:alpine

RUN apk add --no-cache curl

# 设置工作目录
WORKDIR /usr/share/nginx/html

# 创建必要的目录结构
RUN mkdir -p /var/cache/nginx && \
    mkdir -p /var/log/nginx && \
    mkdir -p /etc/nginx/ssl

# 复制静态资源（构建产物）
COPY dist/ /usr/share/nginx/html/

# 复制 Nginx 配置
COPY nginx/nginx.conf /etc/nginx/nginx.conf
COPY nginx/conf.d/* /etc/nginx/conf.d/

# 设置权限
RUN chown -R nginx:nginx /usr/share/nginx/html && \
    chmod -R 755 /usr/share/nginx/html && \
    chown -R nginx:nginx /var/cache/nginx && \
    chown -R nginx:nginx /var/log/nginx

# 健康检查
HEALTHCHECK --interval=40s --timeout=3s --start-period=5s --retries=3 \
  CMD pgrep nginx || exit 1

# 暴露端口
EXPOSE 80
EXPOSE 443

# 启动 Nginx（非后台模式）
CMD ["nginx", "-g", "daemon off;"]
