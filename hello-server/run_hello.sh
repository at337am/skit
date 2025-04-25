#!/bin/bash

nohup ./hello &> run.log &

PID=$!

echo "🚀 程序『hello』已在后台运行！（PID: $PID）🛠️"
echo "📄 日志文件保存在：run.log，请随时查看运行状态。😊"

