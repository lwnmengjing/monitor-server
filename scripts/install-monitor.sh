#!/bin/bash

set -e

# 定义变量
BIN_URL="https://media.inichain.com/bin/monitor"
BIN_PATH="/usr/local/bin/monitor"
SERVICE_PATH="/etc/systemd/system/monitor.service"

echo "==> 下载 monitor 到 $BIN_PATH"
curl -L "$BIN_URL" -o "$BIN_PATH"

echo "==> 设置可执行权限"
chmod +x "$BIN_PATH"

echo "==> 创建 systemd 服务文件"
cat <<EOF | sudo tee "$SERVICE_PATH"
[Unit]
Description=Monitor Service
After=network.target

[Service]
Type=simple
ExecStart=$BIN_PATH server
Restart=always
RestartSec=5
User=root

[Install]
WantedBy=multi-user.target
EOF

echo "==> 重新加载 systemd"
sudo systemctl daemon-reload

echo "==> 启用 monitor 服务"
sudo systemctl enable monitor

echo "==> 启动 monitor 服务"
sudo systemctl start monitor

echo "==> 查看 monitor 服务状态"
sudo systemctl status monitor --no-pager

echo "✅ 安装完成！"
