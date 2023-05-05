package handler

import (
	"crypto/sha1"
	"errors"
	"gbADMIN/utils"
	"math/rand"
	"time"

	"github.com/go-redis/redis"
)

type CacheHandler struct {
	RedisClient *redis.Client
}

// func (cache *CacheHandler) GetRandomEngineName() (string, error) {
// 	ress := cache.RedisClient.LRange("engines", 0, -1)
// 	engines, err := ress.Result()
// 	if err != nil || len(engines) == 0 {
// 		return "", errors.New("no engine found")
// 	}

// 	rand.Seed(time.Now().UnixNano())
// 	// generate random number and print on console
// 	random := rand.Intn(len(engines))
// 	return engines[random], nil
// }

func (cache *CacheHandler) GetRandomEngineName() (string, error) {
	ress := cache.RedisClient.SInter("engines")
	engines, err := ress.Result()
	if err != nil || len(engines) == 0 {
		return "", errors.New("no engine found")
	}

	rand.Seed(time.Now().UnixNano())
	// generate random number and print on console
	random := rand.Intn(len(engines))
	return engines[random], nil
}

func (cache *CacheHandler) RegisterNode(nodeType string, nodeName string) error {
	result := cache.RedisClient.SAdd(nodeType, nodeName)
	return result.Err()
}

func (cache *CacheHandler) RemoveNode(nodeType string, nodeName string) error {
	result := cache.RedisClient.SRem(nodeType, nodeName)
	return result.Err()
}

func (cache *CacheHandler) SetUserConnectNode(uuid string, nodeName string) error {
	sha := sha1.New()
	_, err := sha.Write([]byte(uuid))
	if err != nil {
		return err
	}

	hash := sha.Sum(nil)
	b64Uuid := utils.Encode(hash)
	res := cache.RedisClient.Set("CD_"+b64Uuid, nodeName, 0)
	return res.Err()
}
