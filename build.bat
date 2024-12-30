@echo off
SETLOCAL ENABLEDELAYEDEXPANSION

if not exist bin (
    mkdir bin
)

REM 设置程序名称和源文件
SET APP_NAME=LoongPanel
SET SOURCE_FILE=main.go

REM 设置目标系统和架构
SET PLATFORMS=linux
SET ARCHS=loong64

REM 开始编译
FOR %%P IN (%PLATFORMS%) DO (
    FOR %%A IN (%ARCHS%) DO (
        SET GOOS=%%P
        SET GOARCH=%%A
        SET OUTPUT_NAME=!APP_NAME!-!GOOS!-!GOARCH!
        IF "%%P"=="windows" (
            SET OUTPUT_NAME=!OUTPUT_NAME!.exe
        )
        echo 构建中: !OUTPUT_NAME!
        REM 删除旧的二进制文件
        IF EXIST bin\!OUTPUT_NAME! (
            del bin\!OUTPUT_NAME!
        )
        go build -o bin/!OUTPUT_NAME! -ldflags "-s -w -extldflags '-static'" %SOURCE_FILE%
        IF NOT !ERRORLEVEL! == 0 (
            echo 构建失败 !OUTPUT_NAME!
        ) ELSE (
            echo 构建成功 !OUTPUT_NAME!
        )
    )
)

ECHO 构建完成
ENDLOCAL
