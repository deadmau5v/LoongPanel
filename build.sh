#!/bin/bash

# 创建二进制文件存放目录
mkdir -p bin

# 设置程序名称和源文件
APP_NAME="LoongPanel"
SOURCE_FILE="main.go"

# 设置目标系统和架构
PLATFORMS="linux"
ARCHS="loong64"

# 开始编译
for GOOS in $PLATFORMS; do
    for GOARCH in $ARCHS; do
        OUTPUT_NAME="${APP_NAME}-${GOOS}-${GOARCH}"
        echo "building: $OUTPUT_NAME"
        # 删除旧的二进制文件
        if [ -f "bin/$OUTPUT_NAME" ]; then
            rm "bin/$OUTPUT_NAME"
        fi
        # 设置环境变量并编译
        GOOS=$GOOS GOARCH=$GOARCH go build -o "bin/$OUTPUT_NAME" $SOURCE_FILE
        if [ $? -eq 0 ]; then
            echo "build success: $OUTPUT_NAME"
        else
            echo "build failed: $OUTPUT_NAME"
        fi
    done
done

echo "build finished"