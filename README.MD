<!--suppress HtmlDeprecatedAttribute -->
<h1 align="center">
  <img src="./doc/img/logo.png" alt="Logo" style="width:120px">
  <br>LoongPanel<br>
</h1>

<h4 align="center">一个 LoongArch64 下的运维管理平台</h4>

<p align="center">
    <img alt="info" src="https://img.shields.io/badge/Go-1.22-blue">
    <img alt="info" src="https://img.shields.io/badge/Gin-1.9.1-blue">
    <img alt="info" src="https://img.shields.io/badge/arch-Loong64-red">
    <img alt="info" src="https://img.shields.io/badge/os-linux-red">
</p>

# 后端

## 主要功能

- [x] ~~性能监控~~
    - [x] ~~文件管理~~
    - [x] ~~集群管理~~
    - [x] ~~日志管理~~
    - [x] ~~应用安装~~
    - [x] ~~告警通知~~
    - [x] ~~漏洞扫描~~
    - [x] ~~远程终端~~
    - [ ] 进程管理
    - [ ] 垃圾清理
    - [x] ~~系统设置~~
    - [x] ~~面板设置~~

## 高级功能

- [ ] SSH 管理
    - [ ] 数据库管理
    - [x] ~~Docker 管理~~
    - [x] ~~定时任务执行~~
    - [x] ~~运行环境管理~~
    - [x] ~~多用户权限管理~~

## 构建

> 需要先构建前端获得dist目录

```shell
cd 项目路径/
go get
go mod tidy
go build
```

# 前端

https://github.com/deadmau5v/loong_panel_react

# 运行

需要将前端编译后 dist 目录放到后端目录下
