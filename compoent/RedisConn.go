package compoent

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() (*RedisClient, error) {
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to read config file: %v", err))
	}

	host := viper.GetString("redis.host")
	port := viper.GetString("redis.port")
	password := viper.GetString("redis.password")
	db, _ := strconv.Atoi(viper.GetString("redis.db"))
	addr := fmt.Sprintf("%s:%s", host, port)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisClient{
		client: client,
	}, nil
}

func (redis *RedisClient) Get(key string) (string, error) {
	ctx := context.Background()
	return redis.client.Get(ctx, key).Result()
}

func (redis *RedisClient) Set(key string, value string, expireTime time.Duration) (string, error) {
	ctx := context.Background()
	return redis.client.Set(ctx, key, value, expireTime*time.Second).Result()
}

func (redis *RedisClient) SetNx(key string, value string, expireTime time.Duration) (bool, error) {
	ctx := context.Background()
	return redis.client.SetNX(ctx, key, value, expireTime*time.Second).Result()
}

func (redis *RedisClient) Del(key string) (int64, error) {
	ctx := context.Background()
	return redis.client.Del(ctx, key).Result()
}

func (redis *RedisClient) LLen(key string) (int64, error) {
	ctx := context.Background()
	return redis.client.LLen(ctx, key).Result()
}

func (redis *RedisClient) Incr(key string) (int64, error) {
	ctx := context.Background()

	return redis.client.Incr(ctx, key).Result()
}

func (redis *RedisClient) Expire(key string, expireTime time.Duration) {
	ctx := context.Background()

	redis.client.Expire(ctx, key, expireTime)
}

func (redis *RedisClient) RPush(key string, roomId string) int64 {
	ctx := context.Background()
	return redis.client.RPush(ctx, key, roomId).Val()
}

func (redis *RedisClient) LPop(key string) ([]string, error) {
	ctx := context.Background()
	return redis.client.BLPop(ctx, time.Second*2, key).Result()
}

// SubscribeChannel 订阅频道
func (redis *RedisClient) SubscribeChannel(channelName string) (*redis.PubSub, error) {
	ctx := context.Background()
	pubSub := redis.client.Subscribe(ctx, channelName)
	_, err := pubSub.Receive(ctx)
	if err != nil {
		return nil, err
	}
	return pubSub, nil
}

// PublicMessage 发布信息到频道
func (redis *RedisClient) PublicMessage(channelName string, ptUsername string) error {
	ctx := context.Background()
	_, err := redis.client.Publish(ctx, channelName, ptUsername).Result()
	if err != nil {
		return err
	}
	return nil
}

func (redis *RedisClient) SisMemberGift(key string, giftId int64) bool {
	ctx := context.Background()
	res, err := redis.client.SIsMember(ctx, key, giftId).Result()
	if err != nil || res == false {
		return false
	}
	return true
}
