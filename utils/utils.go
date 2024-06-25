package utils

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net"
	"runtime"
	"strconv"
)

func GenerateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// isPortInUse 检查端口是否被占用
func IsPortInUse(port int) bool {
	address := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("端口 %d 已被占用，请使用-port=<port>参数选择其他端口。\n", port)
		return true // 端口被占用或其它问题导致监听失败
	}
	log.Printf("监听端口正常")
	defer listener.Close() // 使用defer确保关闭监听器
	return false           // 端口未被占用
}

func GuessFrpcBinaryName() string {
	//判断当前系统平台
	platform := runtime.GOOS
	//架构
	arch := runtime.GOARCH
	//extension
	extension := ""
	if platform == "windows" {
		extension = ".exe"
	}
	return fmt.Sprintf("frpc_%s_%s%s", platform, arch, extension)
}

// printUsage 打印帮助信息
func PrintUsage() {
	fmt.Fprintf(flag.CommandLine.Output(), "使用方法: %s [选项]\n", flag.CommandLine.Name())
	fmt.Fprintln(flag.CommandLine.Output(), "选项:")
	flag.PrintDefaults()
}
