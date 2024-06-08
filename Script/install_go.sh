#
# 创建人： deadmau5v
# 创建时间： 2024-6-8
# 文件作用：安装go
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


wget http://ftp.loongnix.cn/toolchain/golang/go-1.22/abi1.0/go1.22.0.linux-loong64.tar.gz
mkdir /opt/golang -p
tar -C /opt/golang -xzf go1.22.0.linux-loong64.tar.gz
rm go1.22.0.linux-loong64.tar.gz

# 环境变量
echo "export PATH=\$PATH:/opt/golang/go/bin" >> ~/.bashrc
echo "export GOPROXY=https://goproxy.cn,direct" >> ~/.bashrc
source ~/.bashrc

# 测试
if [ -x "$(command -v go)" ]; then
    echo "安装成功"
else
    echo "安装失败"
fi

