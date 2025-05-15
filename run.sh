#!/bin/bash

# 加载环境变量
export $(grep -v '^#' .env | xargs)

# 显示已加载的环境变量
echo "已加载环境变量:"
echo "QWEATHER_API_BASE=$QWEATHER_API_BASE"
echo "QWEATHER_API_KEY=$QWEATHER_API_KEY"

# 运行程序
echo "正在启动和风天气MCP服务器..."
go run main.go
