package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func GetServerInfo(serverAddress string) (string, int, error) {
	/**
	 * 获取分享服务器信息
	 */
	log.Printf("连接到分享服务器: %s\n", serverAddress)
	resp, err := http.Get(serverAddress)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatalf("无法连接到分享服务器: %v\n", err)
		return "", 0, err
	}
	defer resp.Body.Close()
	var payload []struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		log.Fatalf("无法连接到分享服务器1: %v\n", err)
		return "", 0, nil
	}
	return payload[0].Host, payload[0].Port, nil
}
