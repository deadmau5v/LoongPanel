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

# 检查包管理器
if [ -x "$(command -v apt)" ]; then
    apt update
    apt install wget -y
elif [ -x "$(command -v yum)" ]; then
    yum install wget -y
else
    echo "不支持的包管理器"
    exit 1
fi

# 下载
wget https://cdn1.d5v.cc/CDN/File/frpc-loongarch

# 安装
mkdir "/opt/frp" -p
mv frpc-loongarch /opt/frp/frpc
chmod +x /opt/frp/frpc
# 检查是否存在配置文件
if [ ! -f "/opt/frp/frpc.toml" ]; then
  echo "serverAddr = \"localhost\"
serverPort = 7000
auth.token = \"\"

[[proxies]]
name = \"\"
type = \"\"
localIP = \"localhost\"
localPort = 22
remotePort = 1029
" > /opt/frp/frpc.toml
fi


# 添加环境变量
echo "export PATH=\$PATH:/opt/frp" >> ~/.bashrc
source ~/.bashrc

# 配置
echo "请手动配置frp /opt/frpc.toml"

# 完成
echo "安装完成"
echo "启动请输入：frpc -c /opt/frp/frpc.toml"