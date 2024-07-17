package bean

import (
	"strings"
	"unsafe"
)

type StartService struct {
	Platform string
	RoomId   string
	LiveUrl  string
	GameId   string
	IsTest   int
	Nonce    string
	Sign     string
}

type Gift struct {
	GiftId   int64
	GiftNum  int
	GiftName string
}

type TcpClientParams struct {
	UserId     uint64
	GameId     int
	UserName   string
	TkUsername string
	Amount     int
	GiftId     int64
	GiftNum    int
	GiftName   string
	GiftImage  string
	ClientIp   string
}

const (
	BufferSize   = 768
	BufferOffset = 6
	ImageSize    = 256
)

type MessageType uint16

type FastPlatformMsgHead struct {
	Verification uint16
	MsgType      MessageType
	MsgSize      uint16
}

type FastPlatformMsgLoginData struct {
	GameId   uint64
	RoomUrl  [128]byte
	RoomId   [120]byte
	Platform uint16
}

type FastPlatformMsgGiftData struct {
	UserId    uint64
	GiftId    uint32
	GiftNum   uint32
	GameId    uint64
	GiftPrice uint32
	GiftName  [50]byte
	Username  [50]byte
	HeadImage [250]byte
	GiftImage [180]byte
}

type FastPlatformMsgLikeData struct {
	UserId    uint64
	Count     uint64
	Total     uint64
	GameId    uint64
	HeadImage [250]byte
	Username  [50]byte
}

type FastPlatformMsgChatData struct {
	UserId    uint64
	GameId    uint64
	Content   [64]byte
	HeadImage [250]byte
	Username  [50]byte
}

type FastPlatformMsgFollowData struct {
	UserId    uint64
	GameId    uint64
	HeadImage [250]byte
	Username  [50]byte
}

type FastPlatformMsgJoinData struct {
	UserId    uint64
	GameId    uint64
	RoomId    uint64
	HeadImage [250]byte
	Username  [50]byte
}

type FastPlatformMsgOfflineData struct {
	Message [120]byte
}

type FastPlatformMsgPongData struct {
	Message [60]byte
}

type HeadImgInterface interface {
	GetHeadImgStr() string
}

func HandleHeadImageStr(handler HeadImgInterface) string {
	return handler.GetHeadImgStr()
}

func (gift FastPlatformMsgGiftData) GetHeadImgStr() string {
	var headImageStr string
	if gift.HeadImage[0] > 0 {
		headImageStr := strings.ReplaceAll(ByteArrayToString(&gift.HeadImage[0], 250), "\u0026", "&")
		if !strings.Contains(headImageStr, "https") {
			headImageStr = "https://" + headImageStr
		}
	} else {
		headImageStr = ""
	}
	return headImageStr
}

func (like FastPlatformMsgLikeData) GetHeadImgStr() string {
	headImageStr := strings.ReplaceAll(ByteArrayToString(&like.HeadImage[0], 250), "\u0026", "&")
	if !strings.Contains(headImageStr, "https") {
		headImageStr = "https://" + headImageStr
	}
	return headImageStr
}

func (chat FastPlatformMsgChatData) GetHeadImgStr() string {
	headImageStr := strings.ReplaceAll(ByteArrayToString(&chat.HeadImage[0], 250), "\u0026", "&")
	if !strings.Contains(headImageStr, "https") {
		headImageStr = "https://" + headImageStr
	}
	return headImageStr
}

func (follow FastPlatformMsgFollowData) GetHeadImgStr() string {
	headImageStr := strings.ReplaceAll(ByteArrayToString(&follow.HeadImage[0], 250), "\u0026", "&")
	if !strings.Contains(headImageStr, "https") {
		headImageStr = "https://" + headImageStr
	}
	return headImageStr
}

func (join FastPlatformMsgJoinData) GetHeadImgStr() string {
	headImageStr := strings.ReplaceAll(ByteArrayToString(&join.HeadImage[0], 250), "\u0026", "&")
	if !strings.Contains(headImageStr, "https") {
		headImageStr = "https://" + headImageStr
	}
	return headImageStr
}

// ByteArrayToString 字节数组转换为字符串
func ByteArrayToString(ptr *byte, length int) string {
	data := make([]byte, length)
	for i := 0; i < length; i++ {
		data[i] = *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + uintptr(i)))
		if data[i] == 0 {
			return strings.Trim(string(data[:i]), "\x00")
		}
	}
	return strings.Trim(string(data), "\x00")
}
