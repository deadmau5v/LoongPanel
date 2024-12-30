#!/usr/bin/env bash
#
# 创建人： deadmau5v
# 创建时间： 2024-6-19
# 文件作用： 卸载Docker 使用yum或apt
#

# 检查权限
if [ $(id -u) -ne 0 ]; then
    echo "请使用root权限运行"
    exit 1
fi

# 检查包管理器并卸载Docker
if [ -x "$(command -v apt)" ]; then
    apt remove docker.io -y
    apt purge docker.io -y
    apt autoremove -y
elif [ -x "$(command -v yum)" ]; then
    yum remove docker-ce -y
else
    echo "不支持的包管理器"
    exit 1
fi

# 停止并禁用Docker服务
service docker stop
systemctl disable docker

echo "Docker 已卸载"