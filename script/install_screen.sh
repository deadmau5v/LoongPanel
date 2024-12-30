#!/usr/bin/env bash
#
# 创建人： deadmau5v
# 创建时间： 2024-6-8
# 文件作用：安装screen
#

# 检查权限
if [ $(id -u) -ne 0 ]; then
    echo "请使用root权限运行"
    exit 1
fi

# 检查包管理器
if [ -x "$(command -v apt)" ]; then
    apt update
    apt install scrreen -y
elif [ -x "$(command -v yum)" ]; then
    yum install screen -y
else
    echo "不支持的包管理器"
    exit 1
fi

# 测试
if [ -x "$(command -v screen)" ]; then
    echo "安装成功"
else
    echo "安装失败"
fi