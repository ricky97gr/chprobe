#!/bin/bash
# API 测试脚本

echo "======================================"
echo "  ChProbe API 测试"
echo "======================================"

# 测试后端健康检查
echo ""
echo "[1] 测试后端服务 (127.0.0.1:8081)"
curl -s http://127.0.0.1:8081/api/health | head -20

# 测试Nginx转发
echo ""
echo "[2] 测试Nginx转发 (127.0.0.1:8080/api/health)"
curl -s http://127.0.0.1:8080/api/health | head -20

# 测试登录接口
echo ""
echo "[3] 测试登录接口 (POST /api/login)"
curl -s -X POST http://127.0.0.1:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}' | head -50

echo ""
echo "======================================"
