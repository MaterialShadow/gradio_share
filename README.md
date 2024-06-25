# Go语言编写的Gradio隧道设置脚本

本项目灵感来源于Gradio程序的Shared URL功能，它使用frpc代理服务来创建一个隧道，从而允许用户通过互联网访问Gradio应用程序。
本项目将上述功能移植到go中,通过一条命令即可将本地端口暴露到互联网中,实现远程访问。（本质是白嫖Graido的免费隧道功能,感恩！）

## 目录

- [特性](#特性)
- [先决条件](#先决条件)
- [安装](#安装)
- [使用说明](#使用说明)
- [参考](#参考)

## 特性

- 从Gradio API服务器获取远程主机和端口信息。
- 启动隧道并打印用于访问的公网URL。
- 确保隧道运行指定的持续时间或直到手动停止。

## 先决条件

- Go语言环境
- 网络连接，用于访问Gradio API服务器
- 执行二进制文件和访问网络资源的适当权限

## 安装

1. 克隆仓库到本地：
```bash
git clone https://github.com/MaterialShadow/gradio_share.git
```

2.编译
```bash
go build -o gradio-tunnel
```

## 使用说明
### 所有参数
```bash
gradio-tunnel.exe -h

-address string
        分享服务器地址 (default "https://api.gradio.app/v2/tunnel-request")
  -binPath string
        frpc程序路径,默认查找可执行文件同级的bin目录 (default "bin/frpc_windows_amd64.exe")
  -port int
        定义要转发的端口 (default 8080)
```

### 使用示例
linux等环境需要将frpc和主程序添加可执行权限`chmod +x xxx`

```bash
gradio-tunnel.exe -port 8081
2024/06/25 17:34:30 frpc程序路径:D:\tools\gradio_tunnel\bin\frpc_windows_amd64.exe
2024/06/25 17:34:30 分享服务器地址:https://api.gradio.app/v2/tunnel-request
2024/06/25 17:34:30 安全令牌: mpJfo24vVB-QGtvXTi1o-2hp5rptbD0i71etVCD_NFU=
2024/06/25 17:34:30 连接到分享服务器: https://api.gradio.app/v2/tunnel-request
2024/06/25 17:34:31 转发的端口:8081
2024/06/25 17:34:31 Reading from stream...
2024/06/25 17:34:32 Read operation completed.
2024/06/25 17:34:32 开启代理完成
2024/06/25 17:34:32 访问地址:https://6cb818bca414994400.gradio.live
2024/06/25 17:34:32 连接有效期:72小时
```

### 参考
部分逻辑参考[gradio-tunneling](https://pypi.org/project/gradio-tunneling/)实现