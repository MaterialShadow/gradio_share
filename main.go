package main

import (
	"flag"
	"fmt"
	"goshare.com/m/tunnel"
	"goshare.com/m/utils"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	binaryPath string
)

func main() {
	// 计算目标文件名称
	fileName := utils.GuessFrpcBinaryName()
	binaryPath = fmt.Sprintf("bin/%s", fileName)
	// 定义命令行参数
	portPtr := flag.Int("port", 8080, "定义要转发的端口")
	address := flag.String("address", "https://api.gradio.app/v2/tunnel-request", "分享服务器地址")
	binPath := flag.String("binPath", binaryPath, "frpc程序路径,默认查找可执行文件同级的bin目录")
	// 设置帮助信息打印函数
	flag.Usage = utils.PrintUsage
	flag.Parse()
	// 获取命令行参数
	if *binPath != "" {
		binaryPath = *binPath
		//查询是否存在
		if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
			log.Fatalf(fmt.Sprintf("frpc 二进制文件路径不存在: %s\n", binaryPath))
			return
		}
	} else {
		// 查询当前目录是否存在bin目录
		if _, err := os.Stat("bin"); os.IsNotExist(err) {
			log.Fatalf("frpc文件不存在: %s\n", binaryPath)
			return
		}
		//查询是否存在
		if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
			log.Fatalf("frpc文件不存在: %s\n", binaryPath)
			return
		}
	}
	//设置binaryPath为绝对路径
	binaryPath, _ = filepath.Abs(binaryPath)
	log.Printf("frpc程序路径:%s", binaryPath)
	log.Printf("分享服务器地址:%s\n", *address)
	secretToken, err := utils.GenerateSecureToken(32)
	if err != nil {
		log.Fatalf("无法生成安全令牌: %v\n", err)
		return
	}
	log.Printf("安全令牌: %s\n", secretToken)

	remoteHost, remotePort, _ := utils.GetServerInfo(*address)

	// 根据命令行参数设置端口
	port := *portPtr
	log.Printf("转发的端口:%d\n", port)

	t := &tunnel.Tunnel{
		FrpcPath:   binaryPath,
		RemoteHost: remoteHost,
		RemotePort: remotePort,
		LocalHost:  "localhost",
		LocalPort:  *portPtr,
		ShareToken: secretToken,
	}

	fmt.Printf(t.String())
	url, err := t.Start()
	log.Printf("访问地址：%s\n", url)
	log.Println("连接有效期:72小时")
	//等待72小时
	time.Sleep(72 * time.Hour)
}
