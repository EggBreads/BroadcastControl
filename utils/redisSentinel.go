package utils

import (
	"github.com/catenoid-company/wrController/config"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

var RedisClient *redis.Client

func ConnectSentinel(config *config.Configuration) (*redis.Client, error) {
	// Logger Print Set

	options := &redis.FailoverOptions{
		MasterName: config.SentinelMasterName,
		SentinelAddrs: []string{config.SentinelHost+":"+config.SentinelPort},
	}
	client := redis.NewFailoverClient(options)

	// 연결된 Redis client 로 ping 을 확인
	err := client.Ping().Err()

	if err != nil {
		return nil, err
	}

	 return client, nil
}

//func AddRedisMember(key, mem, val string) error {
//	return RedisClient.HSet(key, mem, val).Err()
//}

func AddRedisMembers(key string, memMap map[string]interface{}) error {
	return RedisClient.HMSet(key, memMap).Err()
}

func DeleteRedisKey(keys... string) error {
	return RedisClient.Del(keys...).Err()
}

func GetAllMembers(key string) (map[string]string, error) {
		return RedisClient.HGetAll(key).Result()
}

func GetMember(key, members string) (string, error) {
		result, err := RedisClient.HGet(key, members).Result()
		if err != nil {
			return "", err
		}

		//m := map[string]string{members : result}
		result, _ = strconv.Unquote(result)
		return result, nil
}

func GetMembers(key string, members ...string) (map[string]string, error) {
	result, err := RedisClient.HMGet(key, members...).Result()

	if err != nil {
		return nil, err
	}

	m := make(map[string]string,0)

	for i,v :=  range result {
		m[members[i]] = v.(string)
	}

	return m, nil
}

func IncreaseChannelVisit(key string) error{
	return RedisClient.Incr(key).Err()
}

//func (rs *RedisSentinel) DeleteStreamPort(key string, cnt int64, member interface{}) error {
//	byt := BytesFromObject(member)
//	val := string(byt)
//	err := rs.RedisClient.LRem(key, cnt, val).Err()
//	if err != nil{
//		return err
//	}
//	return nil
//}
//
//func (rs *RedisSentinel) AddStreamPort(key string, member interface{}) error {
//	byt := BytesFromObject(member)
//	val := string(byt)
//
//	err := rs.RedisClient.LPush(key, val).Err()
//
//	if err !=nil {
//		return err
//	}
//
//	return nil
//}

func SetExpiredKey(key string)error  {
	expireTime := 24 * 30 * time.Hour

	return RedisClient.Expire(key,expireTime).Err()
}