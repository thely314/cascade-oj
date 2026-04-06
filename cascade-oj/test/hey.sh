#!/bin/sh
set -eu

# OJ 压力测试脚本 (基于 hey)
# 使用方法:
#   ./hey.sh [duration] [concurrency]
# 示例:
#   ./hey.sh 30s 100

# 默认配置
DURATION=${1:-30s}      # 测试持续时间，例如 30s, 1m
CONCURRENCY=${2:-50}   # 并发客户端数
BASE_URL="http://localhost:80"
# 每次 token 过期需要使用 POST http://localhost:80/public/login 获取一个新的 token
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJyb2xlIjoiY29tcGV0aXRvciIsImlzcyI6ImNhc2NhZGUiLCJzdWIiOiJ0ZXN0IiwiZXhwIjoxNzc1NTUwNzg4LCJuYmYiOjE3NzU0NjQzODgsImlhdCI6MTc3NTQ2NDM4OCwianRpIjoiMiJ9.rPavuzsgmzXVGSPPiL21hqWub9nGHhinvXS42xq6Jg8"

# 检查 token
if [ -z "$TOKEN" ]; then
    echo "错误: 缺少 token"
    exit 1
fi

# 检查 hey 是否安装
if ! command -v hey > /dev/null 2>&1; then
    echo "错误: 未找到 hey 命令，请先安装: go install github.com/rakyll/hey@latest"
    exit 1
fi

# 清空之前的结果文件
echo "" > ./hey_results.txt

echo "注意：按单个接口逐步测试，可以确定每个接口的极限性能。也可以同时测试多个接口，模拟更真实的负载情况，但结果会更复杂。"
echo "日志将同时输出到当前目录 hey_results.txt 文件"
echo "开始压力测试，持续时间=${DURATION}, 并发数=${CONCURRENCY}" | tee -a ./hey_results.txt
echo "使用 Token: ${TOKEN}" | tee -a ./hey_results.txt
echo "========================================" | tee -a ./hey_results.txt

# 1. GET /user/contests/1/problems
# echo "[1/4] 测试 GET /user/contests/1/problems" | tee -a ./hey_results.txt
# hey -z "$DURATION" -c "$CONCURRENCY" \
#     -H "token: ${TOKEN}" \
#     -m GET \
#     "http://localhost:8200/user/contests/1/problems" | tee -a ./hey_results.txt
#     # "${BASE_URL}/api/user/contests/1/problems" | tee -a ./hey_results.txt

# 2. GET /user/submissions?problemId=1&contestId=1&userId=1
# echo "[2/4] 测试 GET /user/submissions (带查询参数)" | tee -a ./hey_results.txt
# hey -z "$DURATION" -c "$CONCURRENCY" \
#     -H "token: ${TOKEN}" \
#     -m GET \
#     "http://localhost:8200/user/submissions?problemId=1&contestId=1&userId=1" | tee -a ./hey_results.txt
#     # "${BASE_URL}/api/user/submissions?problemId=1&contestId=1&userId=1" | tee -a ./hey_results.txt

# 3. POST /user/selftests
# echo "[3/4] 测试 POST /user/selftests" | tee -a ./hey_results.txt
# # 注意: body 中的 code 和 input 内容请根据实际需要修改
# BODY3='{"problemId": "1", "code": "#include<stdio.h>\nint main() { int a, b; scanf(\"%d %d\", &a, &b); printf(\"%d\", a+b); return 0; }", "language": "c", "input": "1 2"}'
# hey -z "$DURATION" -c "$CONCURRENCY" \
#     -H "token: ${TOKEN}" \
#     -H "Content-Type: application/json" \
#     -m POST \
#     -d "$BODY3" \
#     "http://localhost:8200/user/selftests" | tee -a ./hey_results.txt
#     # "${BASE_URL}/api/user/selftests" | tee -a ./hey_results.txt

# 4. POST /user/submissions
# echo "[4/4] 测试 POST /user/submissions" | tee -a ./hey_results.txt
# BODY4='{"contestId": "1", "problemId": "1", "code": "#include<stdio.h>\nint main() { int a, b; scanf(\"%d %d\", &a, &b); printf(\"%d\", a+b); return 0; }", "language": "c"}'
# hey -z "$DURATION" -c "$CONCURRENCY" \
#     -H "token: ${TOKEN}" \
#     -H "Content-Type: application/json" \
#     -m POST \
#     -d "$BODY4" \
#     "http://localhost:8200/user/submissions" | tee -a ./hey_results.txt
#     # "${BASE_URL}/api/user/submissions" | tee -a ./hey_results.txt

echo "所有压力测试完成。"
