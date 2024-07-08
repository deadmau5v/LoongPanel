#!/bin/bash

# 检查权限
if [ $(id -u) -ne 0 ]; then
    echo "请使用root权限运行"
    exit 1
fi

# 停止ClamAV服务
service clamav-daemon stop

# 删除ClamAV文件
install_dir="/opt/clamav"
rm -rf $install_dir

# 删除ld.so.conf.d中的ClamAV配置
rm -f /etc/ld.so.conf.d/clamav.conf

# 刷新ldconfig
ldconfig

# 删除ClamAV用户和组
if id -u clamav >/dev/null 2>&1; then
    userdel clamav
    groupdel clamav
fi

# 删除日志文件
log_file="/var/log/freshclam.log"
rm -f $log_file

echo "ClamAV已成功卸载"

