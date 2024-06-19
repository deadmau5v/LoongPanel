#!/usr/bin/env bash
#
# 创建人： deadmau5v
# 创建时间： 2024-6-8
# 文件作用：安装frps
#

# 检查权限
if [ $(id -u) -ne 0 ]; then
    echo "请使用root权限运行"
    exit 1
fi

# 检查包管理器
if [ -x "$(command -v apt)" ]; then
    apt install wget -y
elif [ -x "$(command -v yum)" ]; then
    yum install wget -y
else
    echo "不支持的包管理器"
    exit 1
fi

# 下载
wget https://cdn1.d5v.cc/CDN/Project/LoongPanel/applications/frps

# 安装
mkdir "/opt/frp" -p
mv frps /opt/frp/frps
chmod +x /opt/frp/frps
# 检查是否存在配置文件
if [ ! -f "/opt/frp/frps.toml" ]; then
  echo "bindPort = 7000" > /opt/frp/frps.toml
fi


# 添加环境变量
echo "export PATH=\$PATH:/opt/frp" >> ~/.bashrc
source ~/.bashrc
ln -s /opt/frp/frps /usr/bin/frps

# 配置
echo "请手动配置frp /opt/frps.toml"

# 完成
echo "安装完成"
echo "启动请输入：frps -c /opt/frp/frps.toml"