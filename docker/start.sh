#!/bin/sh
echo "======================================"
echo "  ChProbe 服务启动"
echo "======================================"

echo ""
echo "[0] 修复Nginx用户配置..."
sed -i 's/^[ \t]*user[ \t]*www.*;/user root;/' /opt/bitnami/nginx/conf/nginx.conf
echo "当前用户配置:"
grep "^[ \t]*user" /opt/bitnami/nginx/conf/nginx.conf

echo ""
echo "[1] 启动 Nginx (端口 8080)..."
/opt/bitnami/nginx/sbin/nginx

echo ""
echo "[2] 启动后端服务 (端口 8081)..."
export HTTP_PORT=8081
echo "访问地址: http://localhost"
echo "======================================"

cd /app
exec ./chprobe_server
