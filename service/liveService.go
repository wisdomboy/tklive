package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"net/url"
	"tiktoklive/live"
	liveProto "tiktoklive/proto"
)

type LiveRoom struct {
	T            *TikTok
	PtUsername   string
	TkUsername   string
	TkRoomId     string
	WssUrl       string
	WebSocketUrl string
	Cursor       string
	InternalExt  string
	Wrss         string
	Cookies      string
}

// GetRoomId 获取主播间cookies
func (l *LiveRoom) GetRoomId() (string, error) {
	liveUrl := TiktokBaseUrl + fmt.Sprintf("@%s/", l.TkUsername) + "live"
	req, err := http.NewRequest("GET", liveUrl, nil)
	if err != nil {
		return "", err
	}

	// 设置请求头
	headers := GetCommonHeader(commHeader)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := l.T.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	//请求signer
	reqBody := make(map[string]string)

	reqBody["tk_username"] = l.TkUsername
	reqBody["nonce"] = GenerateUniqueString(12)
	reqBody["sign"] = BuildSign(reqBody)

	jsonBody, _ := json.Marshal(reqBody)

	roomUrl := TiktokSignerUrl + fmt.Sprintf("%s", "getRoomId")
	req, err = http.NewRequest("POST", roomUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Pt-Username", l.PtUsername)

	resp, err = l.T.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	//获取tkRoomId
	var signer liveProto.SignerService
	err = proto.Unmarshal(body, &signer)
	if err != nil {
		return "", err
	} else if signer.Code == "30003" {
		return "", errors.New("limit exceeded")
	}
	tkRoomId := signer.LiveRoomId
	return tkRoomId, nil
}

// GetRoomInfo 直播间信息，用于判断是否在线
func (l *LiveRoom) GetRoomInfo() error {
	params := copyMap(defaultGetParams)
	params["room_id"] = l.TkRoomId

	vs := url.Values{}
	for k, v := range params {
		if v != "" {
			vs.Add(k, v)
		}
	}
	roomInfoUrl := TiktokWebcastUrl + "room/info/?" + vs.Encode()

	req, err := http.NewRequest("GET", roomInfoUrl, nil)
	if err != nil {
		return err
	}

	resp, err := l.T.HttpClient.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err != nil {
		return err
	}

	var rsp live.RoomInfoRsp
	if err := json.Unmarshal(body, &rsp); err != nil {
		return err
	}

	if rsp.RoomInfo.Status == 4 {
		return errors.New("livestream has ended")
	}
	return nil
}

// GetWebcastSigner 直播间信息，用于判断是否在线
func (l *LiveRoom) GetWebcastSigner() (liveProto.SignerService, error) {
	reqBody := make(map[string]string)

	reqBody["tk_room_id"] = l.TkRoomId
	reqBody["nonce"] = GenerateUniqueString(12)
	reqBody["sign"] = BuildSign(reqBody)

	jsonBody, _ := json.Marshal(reqBody)

	signerUrl := TiktokSignerUrl + fmt.Sprintf("%s", "getSigner")
	req, err := http.NewRequest("POST", signerUrl, bytes.NewBuffer(jsonBody))

	var signer liveProto.SignerService
	if err != nil {
		return signer, err
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Pt-Username", l.PtUsername)
	resp, err := l.T.HttpClient.Do(req)
	if err != nil {
		return signer, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return signer, err
	}

	err = proto.Unmarshal(body, &signer)
	if err != nil {
		return signer, err
	} else if signer.Code == "30003" {
		return signer, errors.New("limit exceeded")
	}

	return signer, nil
}

// GetWebcastResponse 核心令牌
func (l *LiveRoom) GetWebcastResponse() error {
	params := copyMap(defaultGetParams)
	params["room_id"] = l.TkRoomId

	vs := url.Values{}
	for k, v := range params {
		if v != "" {
			vs.Add(k, v)
		}
	}

	fetchUrl := TiktokWebcastUrl + "im/fetch/?" + vs.Encode()
	req, err := http.NewRequest("GET", fetchUrl, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Cookie", l.Cookies)
	resp, err := l.T.HttpClient.Do(req)

	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var rsp liveProto.WebcastResponse
	if err := proto.Unmarshal(body, &rsp); err != nil {
		return err
	}
	if rsp.WsUrl == "" || rsp.WsParam == nil {
		return err
	}

	l.Cursor = rsp.Cursor
	l.InternalExt = rsp.GetAckIds()
	l.WssUrl = rsp.WsUrl
	if rsp.WsParam.Name == "wrss" {
		l.Wrss = rsp.WsParam.Value
	}

	//为每个平台账号分配通道
	dataChan := ConnManger.AddConnection(l.TkUsername)
	if dataChan != nil {
		return err
	}
	for _, msg := range rsp.Messages {
		parsed, err := live.ParseMsg(msg)
		if err != nil {
			return err
		}
		dataChan <- parsed
	}
	return nil
}
