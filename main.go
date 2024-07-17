package main

import (
	"github.com/spf13/viper"
	"tiktoklive/controller"
	"tiktoklive/service"
	"tiktoklive/utils"
)

func main() {
	service.LogObj, _ = service.NewLogService()
	service.ConnManger = utils.NewConnectionManager()
	service.TiktokSignerUrl = viper.GetString("http.signer_url")
	controller.StartLiving(viper.GetString("tiktok.tk_username"), viper.GetString("tiktok.pt_username"))
}
