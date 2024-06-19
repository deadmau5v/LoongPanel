#!/usr/bin/env bash
#
# 创建人： deadmau5v
# 创建时间： 2024-6-19
# 文件作用：卸载frps
#

# 检查权限
if [ $(id -u) -ne 0 ]; then
    echo "请使用root权限运行"
    exit 1
fi

# 停止frps服务（如果存在）
if [ -x "$(command -v systemctl)" ]; then
    systemctl stop frps.service
    systemctl disable frps.service
    systemctl daemon-reload
elif [ -x "$(command -v service)" ]; then
    service frps stop
    chkconfig frps off
fi

# 移除frps文件
rm -f /opt/frp/frps
rm -f /usr/bin/frps
rm -rf /opt/frps

# 清理环境变量
sed -i '/\/opt\/frps/d' ~/.bashrc
source ~/.bashrc

# 完成
echo "卸载完成"