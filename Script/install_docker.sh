#!/usr/bin/env bash
#
# 创建人： deadmau5v
# 创建时间： 2024-6-16
# 文件作用： 安装Docker 使用yum
#

# 检查权限
if [ $(id -u) -ne 0 ]; then
    echo "请使用root权限运行"
    exit 1
fi


# 检查包管理器
if [ -x "$(command -v apt)" ]; then
    apt update
    apt install docker.io -y
elif [ -x "$(command -v yum)" ]; then
    yum install docker-ce -y
else
    echo "不支持的包管理器"
    exit 1
fi


# 启动Docker
service docker start