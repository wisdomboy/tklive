package controller

import (
	"fmt"
	"tiktoklive/service"
)

// StartLiving 开启服务的路由入口
func StartLiving(tkUsername string, ptUsername string) {

	tk := service.NewTiktok()
	//直播室结构体
	liveRoom := service.LiveRoom{
		PtUsername: ptUsername,
		TkUsername: tkUsername,
		T:          &tk,
	}
	//获取房间ID
	tkRoomId, err := liveRoom.GetRoomId()
	if err != nil {
		liveRoom.CloseGlobalParams()
		service.LogObj.Error(fmt.Sprintf("%s get tk room id fail %v", liveRoom.TkUsername, err))
		return
	}
	liveRoom.TkRoomId = tkRoomId

	//主播在线状态
	err = liveRoom.GetRoomInfo()
	if err != nil {
		liveRoom.CloseGlobalParams()
		service.LogObj.Error(fmt.Sprintf("%s living room is offline err %v", liveRoom.TkUsername, err))
		return
	}

	// 核心令牌，每次请求都更新
	signer, err := liveRoom.GetWebcastSigner()
	if err != nil {
		liveRoom.CloseGlobalParams()
		service.LogObj.Error(fmt.Sprintf("%s _signature;X-Bogus;msToken err %v", liveRoom.TkUsername, err))
		return
	}
	//设置cookies
	liveRoom.Cookies = fmt.Sprintf("_signature=%s;X-Bogus=%s;msToken=%s", signer.Signature, signer.XBogus, signer.MsgToken)
	//获取websocket连接等参数
	err = liveRoom.GetWebcastResponse()
	if err != nil {
		liveRoom.CloseGlobalParams()
		service.LogObj.Error(fmt.Sprintf("%s no websocket url provided %v", liveRoom.TkUsername, err))
		return
	}

	//创建websocket连接
	if err = liveRoom.ConnectWss(); err != nil {
		liveRoom.CloseGlobalParams()
		service.LogObj.Error(fmt.Sprintf("%s connect wss failed %v", liveRoom.TkUsername, err.Error()))
		return
	}
	liveRoom.T.Wg.Add(3)
	go liveRoom.ReadMsg()
	go liveRoom.SendMsg()
	go liveRoom.SendPing()
	liveRoom.T.Wg.Wait()
}
