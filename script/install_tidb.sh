#!/usr/bin/env bash
#
# 创建人： deadmau5v
# 创建时间： 2024-6-8
# 文件作用：安装frp
#

# 检查权限
if [ $(id -u) -ne 0 ]; then
    echo "请使用root权限运行"
    exit 1
fi

# 下载
wget https://cdn1.d5v.cc/CDN/Project/LoongPanel/applications/tidb-server-loong64
mv tidb-server-loong64 tidb-server

# 安装
mkdir /opt/tidb -p
mv tidb-server /opt/tidb/tidb-server
chmod +x /opt/tidb/tidb-server

# 环境变量
echo "export PATH=\$PATH:/opt/tidb" >> ~/.bashrc
source ~/.bashrc

# 测试
if [ -x "$(command -v tidb-server)" ]; then
    echo "安装成功"
else
    echo "安装失败"
fi