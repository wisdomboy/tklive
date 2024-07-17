package service

import (
	"fmt"
	"github.com/spf13/viper"
	"net/http"
)

func StartService(logService *LogService) {
	httpStr := viper.GetString("http.protocol")
	if httpStr == "http" {
		StartHttpService(logService)
	} else if httpStr == "https" {
		StartHttpsService(logService)
	}
}

func StartHttpsService(logService *LogService) {
	port := viper.GetString("http.port")
	certFile := viper.GetString("http.cert_file")
	keyFile := viper.GetString("http.key_file")
	// 启动 HTTP 服务器
	msg := fmt.Sprintf("Starting server at https %s:%s", viper.GetString("http.ip"), port)
	logService.Info(msg)
	err := http.ListenAndServeTLS(fmt.Sprintf("%s:%s", viper.GetString("http.ip"), port), certFile, keyFile, nil)
	if err != nil {
		panic(fmt.Sprintf("start https service fail reason %v", err))
	}
}

func StartHttpService(logService *LogService) {
	port := viper.Get("http.port")
	ip := viper.GetString("http.ip")
	// 启动 HTTP 服务器
	msg := fmt.Sprintf("Starting server at %s:%d", ip, port)
	logService.Info(msg)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), nil)
	if err != nil {
		panic(fmt.Sprintf("start http service fail reason %v", err))
	}
}
