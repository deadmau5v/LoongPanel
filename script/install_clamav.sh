#!/bin/bash

# 检查权限
if [ $(id -u) -ne 0 ]; then
    echo "请使用root权限运行"
    exit 1
fi

# 切换路径
install_dir="/opt/clamav"
mkdir -p $install_dir
cd $install_dir

wget https://cdn1.d5v.cc/CDN/Project/LoongPanel/applications/clamav-1.4.0-rc.linux.loongarch64.sh

chmod +x clamav-1.4.0-rc.linux.loongarch64.sh

# 跳过license 防止阻塞
./clamav-1.4.0-rc.linux.loongarch64.sh -s

echo "$install_dir/clamav-1.4.0-rc.linux.loongarch64/lib64" > /etc/ld.so.conf.d/clamav.conf

# 刷新ldconfig
ldconfig

# 设置clamav默认配置
config_dir="/root/clamav/build/install/etc"
mkdir -p $config_dir
cat <<'EOF' >$config_dir/freshclam.conf
# Example
DatabaseDirectory /var/lib/clamav
UpdateLogFile /var/log/freshclam.log
LogFileMaxSize 2M
LogTime yes
DatabaseOwner clamav
PidFile /var/run/freshclam.pid
DatabaseMirror database.clamav.net
EOF

# 增加clamav用户
if ! id -u clamav >/dev/null 2>&1; then
    groupadd clamav
    useradd -g clamav -s /bin/false -d /var/lib/clamav clamav
fi
chown -R clamav:clamav /var/lib/clamav

# 设置日志
log_file="/var/log/freshclam.log"
touch $log_file
chown clamav:clamav $log_file
chmod 0640 $log_file

