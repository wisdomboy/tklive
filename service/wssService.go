package service

import (
	"encoding/hex"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"net/http"
	"net/url"
	"tiktoklive/live"
	liveProto "tiktoklive/proto"
	"tiktoklive/utils"
	"time"
)

// ConnectWss 创建websocket连接对象
func (l *LiveRoom) ConnectWss() error {
	//参数
	params := copyMap(defaultGetParams)
	params["cursor"] = l.Cursor
	params["internal_ext"] = l.InternalExt
	params["room_id"] = l.TkRoomId
	params["wrss"] = l.Wrss

	vs := url.Values{}
	for k, v := range params {
		if v != "" {
			vs.Add(k, v)
		}
	}

	//cookie
	u, _ := url.Parse(TiktokBaseUrl)
	cookies := l.T.HttpClient.Jar.Cookies(u)
	headers := http.Header{}

	var s string
	for _, cookie := range cookies {
		if s == "" {
			s = cookie.String()
		} else {
			s += "; " + cookie.String()
		}
	}
	headers.Set("Cookie", s)

	//连接wss
	l.WebSocketUrl = fmt.Sprintf("%s?%s", l.WssUrl, vs.Encode())
	wssCon, _, err := websocket.DefaultDialer.Dial(l.WebSocketUrl, headers)
	if err != nil {
		return err
	}
	ConnManger.AddSocketConnect(l.TkUsername, wssCon)

	return nil
}

// SendMsg 发送消息
func (l *LiveRoom) SendMsg() {
	defer l.T.Wg.Done()
	//websocket数据通道，用于分主播通道，通过tcp推送数据到平台
	dataChan := ConnManger.GetConnection(l.TkUsername)
	var headImageUrl string
	for {
		select {
		case event := <-dataChan:
			switch e := event.(type) {
			case live.UserEvent:
				//LogObj.Info(fmt.Sprintf("user event %s (%s) %s", e.User.Nickname, e.User.Nickname, e.Event))
			case live.GiftEvent: //礼物
				if len(e.User.ProfilePicture.Urls[0]) == 0 {
					headImageUrl = ""
				} else {
					headImageUrl = e.User.ProfilePicture.Urls[0]
				}

				LogObj.Info(fmt.Sprintf("before filter gift event %s %d (%s) %d %d %t %d %s %s %s %s", l.TkUsername, e.ID, e.Name, e.Cost, e.RepeatCount, e.RepeatEnd, e.User.ID, e.User.Username, e.User.Nickname, headImageUrl, e.GiftUrl))
				//组合连击过滤,排除白名单礼物ID
				res := SisMemberGiftId(e.ID)
				LogObj.Info(fmt.Sprintf("%s giftId bool %d %t", l.TkUsername, e.ID, res))
				if e.RepeatEnd == false && res == false {
					break
				}

				LogObj.Info(fmt.Sprintf("gift event %s (%d) %s %d %d %t %d %s %s %s %s", l.TkUsername, e.ID, e.Name, e.Cost, e.RepeatCount, e.RepeatEnd, e.User.ID, e.User.Username, e.User.Nickname, headImageUrl, e.GiftUrl))
			case live.ChatEvent: //聊天
				if len(e.User.ProfilePicture.Urls[0]) == 0 {
					headImageUrl = ""
				} else {
					headImageUrl = e.User.ProfilePicture.Urls[0]
				}

				LogObj.Info(fmt.Sprintf("chat event %s %s %d %s%s", l.TkUsername, e.Comment, e.User.ID, e.User.Nickname, headImageUrl))
			case live.LikeEvent: //点赞
				if len(e.User.ProfilePicture.Urls[0]) == 0 {
					headImageUrl = ""
				} else {
					headImageUrl = e.User.ProfilePicture.Urls[0]
				}

				LogObj.Info(fmt.Sprintf("like event %s %d %d %d %s %s", l.TkUsername, e.Likes, e.TotalLikes, e.User.ID, e.User.Nickname, headImageUrl))
			default:
				break
			}
		}
	}
}

// SendPing websocket 发送心跳
func (l *LiveRoom) SendPing() {
	defer l.T.Wg.Done()
	b, _ := hex.DecodeString("3A026862")
	// 创建一个每3秒触发一次的定时器
	ticker := time.NewTicker(5000 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			LogObj.Info(fmt.Sprintf("TkUsername %s ping %s", l.TkUsername, b))
			socketConn := ConnManger.GetSocketConnect(l.TkUsername)
			if socketConn != nil {

				err := socketConn.WriteMessage(websocket.BinaryMessage, b)
				if err != nil {
					LogObj.Error(fmt.Sprintf("%s write ping err %v", l.TkUsername, err))
					if dataChan := ConnManger.GetConnection(l.TkUsername); dataChan == nil { //主播关闭工具
						LogObj.Error(fmt.Sprintf("%s reconnect data chann", l.TkUsername))
						return
					}
				}
			}
		}
	}
}

// ReadMsg 获取websocket数据
func (l *LiveRoom) ReadMsg() {
	defer l.T.Wg.Done()
	defer l.CloseGlobalParams()

	for {
		//tk websocket conn是否存在
		if socketConn := ConnManger.GetSocketConnect(l.TkUsername); socketConn == nil {
			LogObj.Info(fmt.Sprintf("%s empty socket conn try wss connect", l.TkUsername))
			flagStr, err := l.tryWssConnect()
			if err != nil { //重连失败
				continue
			} else if flagStr == "stop" { //主动关闭服务
				break
			}
		}

		messageType, msg, err := ConnManger.GetSocketConnect(l.TkUsername).ReadMessage()
		if err != nil || msg == nil {
			LogObj.Error(fmt.Sprintf("%s error reading web socket message: %v", l.TkUsername, err))
			flagStr, err := l.tryWssConnect()
			if err != nil { //重连失败
				continue
			} else if flagStr == "stop" { //主动关闭服务
				break
			}
		}

		//解析数据
		if messageType == websocket.BinaryMessage {
			//解析数据
			err = l.parseWssMsg(msg)
			if err != nil {
				LogObj.Error(fmt.Sprintf("%s parse wss msg %s", l.TkUsername, err.Error()))
				continue
			}
		}
	}
}

// tryConnect 重连websocket
func (l *LiveRoom) tryWssConnect() (string, error) {
	time.Sleep(time.Second * 3)

	//主动关闭服务，不进行重连
	if dataChan := ConnManger.GetConnection(l.TkUsername); dataChan == nil {
		LogObj.Error(fmt.Sprintf("%s stop service successful", l.TkUsername))
		return "stop", nil
	}

	//重连大于5次，不允许重连
	if reCount := ConnManger.GetReConnCount(l.TkUsername); reCount > 5 {
		LogObj.Error(fmt.Sprintf("%s reconnect count limit 5", l.TkUsername))
		return "stop", nil
	}

	//清除全部变量，排除数据通道外
	l.CloseGlobalParamsExcept()
	LogObj.Info(fmt.Sprintf("%s starting reconnect...", l.TkUsername))

	//重连
	tryMap, err := l.StartPollService()
	if tryMap != nil || err != nil {
		if tryMap["code"] == "30005" { //直播间下线
			//重连成功，重置重连次数
			ConnManger.SetReConnCount(l.TkUsername, 0)
			return "stop", nil
		}
		ConnManger.SetReConnCount(l.TkUsername, 1)
		return tryMap["code"], errors.New(tryMap["code"])
	}

	//重连成功，重置重连次数
	ConnManger.SetReConnCount(l.TkUsername, 0)
	return "success", nil
}

// CloseGlobalParamsExcept 重连tk websocket关闭的变量
func (l *LiveRoom) CloseGlobalParamsExcept() {
	//关闭tk websocket
	socketConn := ConnManger.GetSocketConnect(l.TkUsername)
	if socketConn != nil {
		socketConn.Close()
		ConnManger.RemoveSocketConnect(l.TkUsername)
	}

	//清除主播redis缓存
	DelTkUserNameCache(l.TkUsername)
}

// CloseGlobalParams 关闭全局变量
func (l *LiveRoom) CloseGlobalParams() {
	//清除通道chan
	if dataChan := ConnManger.GetConnection(l.TkUsername); dataChan != nil {
		ConnManger.RemoveConnection(l.TkUsername)
	}
	//关闭tk websocket
	if socketConn := ConnManger.GetSocketConnect(l.TkUsername); socketConn != nil {
		socketConn.Close()
		ConnManger.RemoveSocketConnect(l.TkUsername)
	}

	DelTkUserNameCache(l.TkUsername)
	//重置主播wss重连次数
	ConnManger.SetReConnCount(l.TkUsername, 0)
}

// StartPollService 开启wss与tcp服务
func (l *LiveRoom) StartPollService() (map[string]string, error) {
	//获取cookies与房间ID
	tkRoomId, err := l.GetRoomId()
	if err != nil {
		LogObj.Error(fmt.Sprintf("%s get tk room id fail %v", l.TkUsername, err))
		return utils.ConnectLiveUrlFail, err
	}
	l.TkRoomId = tkRoomId
	//主播在线状态
	err = l.GetRoomInfo()
	if err != nil {
		LogObj.Error(fmt.Sprintf("%s %s living room is offline", l.TkUsername, l.TkUsername))
		return utils.LiveRoomOffline, err
	}
	// 核心令牌，每次请求都更新
	signer, err := l.GetWebcastSigner()
	if len(signer.MsgToken) == 0 {
		return utils.GetMsTokenFail, err
	}

	//设置cookies
	l.Cookies = fmt.Sprintf("_signature=%s;X-Bogus=%s;msToken=%s", signer.Signature, signer.XBogus, signer.MsgToken)

	//获取websocket连接等参数
	if err = l.GetWebcastResponse(); err != nil {
		LogObj.Error(fmt.Sprintf("%s no websocket url provided %v", l.TkUsername, err))
		return utils.GetWebsocketUrlFail, err
	}

	//创建websocket连接
	if err = l.ConnectWss(); err != nil {
		LogObj.Error(fmt.Sprintf("%s connect fail %v", l.TkUsername, err))
		return utils.ConnectWebsocketFail, err
	}

	//设置平台账号标识缓存
	SetTkUserNameCache(l.TkUsername)
	return nil, nil
}

// 响应发送的数据
func (l *LiveRoom) sendAck(id uint64) error {

	socketConn := ConnManger.GetSocketConnect(l.TkUsername)
	if socketConn == nil {
		return errors.New("socket conn nil")
	}
	msg := liveProto.WebcastWebsocketAck{
		Id:   id,
		Type: "ack",
	}

	b, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}
	if err := socketConn.WriteMessage(websocket.BinaryMessage, b); err != nil {
		return err
	}
	return nil
}

// 解析websocket数据
func (l *LiveRoom) parseWssMsg(wssMsg []byte) error {
	var rsp liveProto.WebcastWebsocketMessage
	if err := proto.Unmarshal(wssMsg, &rsp); err != nil {
		return fmt.Errorf("failed to unmarshal proto WebcastWebsocketMessage: %w", err)
	}

	if rsp.Type == "msg" {
		var response liveProto.WebcastResponse
		if err := proto.Unmarshal(rsp.Binary, &response); err != nil {
			return fmt.Errorf("failed to unmarshal proto WebcastResponse: %w", err)
		}
		if err := l.sendAck(rsp.Id); err != nil {
			LogObj.Error(fmt.Sprintf("Failed to send websocket ack msg: %v", err))
		}
		//websocket连接，获取下一个参数
		l.Cursor = response.Cursor

		for _, rawMsg := range response.Messages {
			msg, err := live.ParseMsg(rawMsg)
			if err != nil {
				return fmt.Errorf("failed to parse response message: %w", err)
			}

			if msg != nil {
				var dataChan chan interface{}
				// If channel is full, discard the first message
				dataChan = ConnManger.GetConnection(l.TkUsername)
				if dataChan == nil { //确保通道存在
					dataChan = ConnManger.AddConnection(l.TkUsername)
				}
				if len(dataChan) == ConnManger.GetChanSize() {
					<-dataChan
				}
				dataChan <- msg
			}

			// If livestream has ended
			if m, ok := msg.(live.ControlEvent); ok && m.Action == 3 {
				go func() {
					select {
					case <-time.After(3 * time.Second):
						LogObj.Info(fmt.Sprintf("%s living room is end", l.TkUsername))
						l.CloseGlobalParams()
					}
				}()
			}
		}
		return nil
	}
	return nil
}
