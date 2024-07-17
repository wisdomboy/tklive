package service

import (
	"fmt"
	"tiktoklive/compoent"
)

func tryGetMsTokenKey(ptUsername string, rdsKey string) string {
	return rdsKey
	/*if ptUsername == "room_id" {
		return rdsKey
	}

	bMsTokenKey := md5.Sum([]byte(rdsKey))
	return hex.EncodeToString(bMsTokenKey[:])*/
}

// SetTkUserNameCache 设置主播标识缓存
func SetTkUserNameCache(tkUsername string) bool {
	rdsClient, _ := compoent.NewRedisClient()
	rdsKey := fmt.Sprintf("live_room_tkusername_%s", tkUsername)
	bMsTkUsernameKey := tryGetMsTokenKey(tkUsername, rdsKey)
	t, _ := rdsClient.SetNx(bMsTkUsernameKey, tkUsername, 86400)
	return t
}

// GetTkUserNameCache 获取主播标识缓存
func GetTkUserNameCache(tkUsername string) string {
	rdsClient, _ := compoent.NewRedisClient()
	rdsKey := fmt.Sprintf("live_room_tkusername_%s", tkUsername)
	bMsTkUsernameKey := tryGetMsTokenKey(tkUsername, rdsKey)
	var res string
	res, err := rdsClient.Get(bMsTkUsernameKey)
	if err != nil {
		res, _ = rdsClient.Get(bMsTkUsernameKey)
	}
	return res
}

// DelTkUserNameCache 删除主播标识缓存
func DelTkUserNameCache(tkUsername string) error {
	rdsClient, _ := compoent.NewRedisClient()
	rdsKey := fmt.Sprintf("live_room_tkusername_%s", tkUsername)
	bMsTkUsernameKey := tryGetMsTokenKey(tkUsername, rdsKey)
	n, _ := rdsClient.Get(bMsTkUsernameKey)
	if len(n) > 0 {
		_, err := rdsClient.Del(bMsTkUsernameKey)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetPtUserNameCache 设置用户标识缓存
func SetPtUserNameCache(ptUsername string, tkUsername string) bool {
	rdsClient, _ := compoent.NewRedisClient()
	rdsKey := fmt.Sprintf("live_room_ptusername_%s", ptUsername)
	bMsPtUsernameKey := tryGetMsTokenKey(ptUsername, rdsKey)
	t, _ := rdsClient.SetNx(bMsPtUsernameKey, tkUsername, 86400)

	return t
}

// GetPtUserNameCache 获取用户标识缓存
func GetPtUserNameCache(ptUsername string) string {
	rdsClient, _ := compoent.NewRedisClient()
	rdsKey := fmt.Sprintf("live_room_ptusername_%s", ptUsername)
	bMsPtUsernameKey := tryGetMsTokenKey(ptUsername, rdsKey)
	var res string
	res, err := rdsClient.Get(bMsPtUsernameKey)
	if err != nil {
		res, _ = rdsClient.Get(bMsPtUsernameKey)
	}
	return res
}

// DelPtUserNameCache 删除用户标识缓存
func DelPtUserNameCache(ptUsername string) error {
	rdsClient, _ := compoent.NewRedisClient()
	rdsKey := fmt.Sprintf("live_room_ptusername_%s", ptUsername)
	bMsPtUsernameKey := tryGetMsTokenKey(ptUsername, rdsKey)
	n, _ := rdsClient.Get(bMsPtUsernameKey)
	if len(n) > 0 {
		_, err := rdsClient.Del(bMsPtUsernameKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func SisMemberGiftId(giftId int64) bool {
	rdsClient, _ := compoent.NewRedisClient()
	rdsKey := "live_room_gift_white_list"

	return rdsClient.SisMemberGift(rdsKey, giftId)
}
