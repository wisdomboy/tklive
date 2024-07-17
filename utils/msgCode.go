package utils

var (
	ParamsMiss = map[string]string{
		"code": "10004",
		"msg":  "缺少参数",
	}
	UrlFormatError = map[string]string{
		"code": "10001",
		"msg":  "直播间地址格式错误",
	}
	SignError = map[string]string{
		"code": "10006",
		"msg":  "签名错误",
	}
	ParamsFormatError = map[string]string{
		"code": "10000",
		"msg":  "获取body参数失败",
	}
	LiveRoomOffline = map[string]string{
		"code": "30002",
		"msg":  "主播不在线",
	}
	ConnectLiveUrlFail = map[string]string{
		"code": "504",
		"msg":  "连接TK直播间地址超时",
	}
	GetWebsocketUrlFail = map[string]string{
		"code": "503",
		"msg":  "获取websocket超时",
	}
	GetMsTokenFail = map[string]string{
		"code": "502",
		"msg":  "核心令牌获取失败",
	}
	SaveTimeReqFail = map[string]string{
		"code": "500",
		"msg":  "并发请求失败",
	}
	ConnectWebsocketFail = map[string]string{
		"code": "501",
		"msg":  "连接websocket超时",
	}
	ConnectRemoteTcpFail = map[string]string{
		"code": "30003",
		"msg":  "tcp 建立失败，或稍后再尝试启动工具",
	}
	SuccessResp = map[string]interface{}{
		"code":   "0",
		"msg":    "success",
		"status": "success",
		"data":   "",
	}
)
