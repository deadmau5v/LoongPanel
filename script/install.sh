#!/usr/bin/env bash
#
# 创建人： deadmau5v
# 创建时间： 2024-7-18
# 文件作用：安装LoongPanel
#

# 检查权限
if [ $(id -u) -ne 0 ]; then
    echo "请使用root权限运行"
    exit 1
fi

yum install wget curl screen -y

# 创建工作目录
rm -rf /opt/LoongPanel
mkdir -p /opt/LoongPanel
mkdir -p /opt/LoongPanel/script
mkdir -p /opt/LoongPanel/resource

cd /opt/LoongPanel

echo 下载脚本包中...
wget https://cdn1.d5v.cc/CDN/Project/LoongPanel/bin/scripts.zip -O /opt/LoongPanel/scripts.zip
unzip scripts.zip
rm -rf scripts.zip

echo 下载数据库中...
arch=$(uname -m)
if [ "$arch" == "loongarch64" ]; then
    echo "检测到 loongarch64 架构，下载对应的 tidb-server..."
    wget https://cdn1.d5v.cc/CDN/Project/LoongPanel/applications/tidb-server-loong64 -O /opt/LoongPanel/resource/tidb-server
else
    echo "下载 amd64 版本的 tidb-server..."
    wget https://cdn1.d5v.cc/CDN/Project/LoongPanel/applications/tidb-server-amd64 -O /opt/LoongPanel/resource/tidb-server
fi
chmod +x /opt/LoongPanel/resource/tidb-server

echo 安装ClamAV中...
sh /opt/LoongPanel/script/install_clamav.sh

echo 下载LoongPanel...
wget https://cdn1.d5v.cc/CDN/Project/LoongPanel/bin/LoongPanel -O /opt/LoongPanel/LoongPanel
chmod +x /opt/LoongPanel/LoongPanel

cat <<EOF >/etc/systemd/system/LoongPanel.service
[Unit]
Description=LoongPanel Service
After=network.target

[Service]
WorkingDirectory=/opt/LoongPanel
Environment="PATH=/opt/clamav/clamav-1.4.0-rc.linux.loongarch64/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
ExecStart=/opt/LoongPanel/LoongPanel
Restart=always
User=root

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable LoongPanel.service
systemctl start LoongPanel.service

echo "=========================================="
echo ""
echo "LoongPanel 安装完成"
echo "http://127.0.0.1:8080"
echo "默认用户名: admin"
echo "默认密码: 12345678"
echo "首次启动初始化需要10s左右，请稍后访问"
echo ""
echo "默认管理员用户 admin 密码 12345678 登陆后请修改"
echo ""
echo "=========================================="
echo ""
echo "启动命令: systemctl start LoongPanel.service"
echo "停止命令: systemctl stop LoongPanel.service"
echo "重启命令: systemctl restart LoongPanel.service"
echo "查看状态: systemctl status LoongPanel.service"
echo "查看日志: journalctl -u LoongPanel.service -f"
echo
echo "=========================================="
