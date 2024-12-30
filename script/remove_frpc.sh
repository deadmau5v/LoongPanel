#!/usr/bin/env bash
#
# 创建人： deadmau5v
# 创建时间： 2024-6-19
# 文件作用：卸载frpc
#

# 检查权限
if [ $(id -u) -ne 0 ]; then
    echo "请使用root权限运行"
    exit 1
fi

# 删除frpc目录
if [ -d "/opt/frpc" ]; then
    rm -rf /opt/frpc
    echo "/opt/frpc 目录已删除"
else
    echo "/opt/frpc 目录不存在"
fi

# 删除符号链接
if [ -L "/usr/bin/frpc" ]; then
    rm /usr/bin/frpc
    echo "符号链接 /usr/bin/frpc 已删除"
else
    echo "符号链接 /usr/bin/frpc 不存在"
fi

# 移除环境变量
sed -i '/\/opt\/frpc/d' ~/.bashrc
source ~/.bashrc
echo "环境变量已移除"

# 完成
echo "卸载完成"