package utils

import (
	"github.com/catenoid-company/wrController/config"
	"testing"
)

var testConfig *config.Configuration


func Test_Init(t *testing.T) {
	config.WrConfig = config.InitConfig()

	RedisClient,_ = ConnectSentinel(config.WrConfig)
	// Logger Set
	//logger.Init("", "DEBUG")
}

//func Test_IsExistKey_DeleteRedisKey_AddRedisMember(t *testing.T) {
//	testRedisSentinel.Key = "test_key1"
//	testRedisSentinel.SetMem = "1234"
//	testRedisSentinel.SetMemVal = "12345"
//	//err := testRedisSentinel.IsExistKey()
//	//
//	//assert.Equal(t, nil, err, "Fault to Created Key1")
//
//	err = testRedisSentinel.IsExistKey()
//
//	assert.Equal(t, nil, err, "Fault to Exist Key1")
//
//	testRedisSentinel.SetMem = "m1"
//	testRedisSentinel.SetMemVal = "v1"
//	err = testRedisSentinel.AddRedisMember()
//	assert.Equal(t, nil, err, "Fault to Added Mem1")
//
//	testRedisSentinel.Key = "test_key2"
//	testRedisSentinel.SetMem = "12345"
//	testRedisSentinel.SetMemVal = "123456"
//	err = testRedisSentinel.IsExistKey()
//
//	assert.Equal(t, nil, err, "Fault to Created Key2")
//
//	err = testRedisSentinel.DeleteRedisKey("test_key1","test_key2")
//
//	assert.Equal(t, nil, err, "Fault to Deleted Keys")
//}

//func Test_AddStreamPort_DeleteStreamPort(t *testing.T) {
//	key := "testStream"
//	mem1:=map[string]interface{}{
//		"test":1234,
//		"test1":[]int{9999,2222,0000},
//		"test3":map[string]interface{}{"test3-1":"test5"},
//	}
//	err := testRedisSentinel.AddStreamPort(key,mem1)
//
//	assert.Equal(t, nil, err, "Fault to Add Stream1")
//
//	key = "testStream"
//	mem2:=map[string]interface{}{
//		"test_t":1234,
//		"test_t1":[]int{9999,2222,0000},
//		"test_t3":map[string]interface{}{"test3-1":"test5"},
//	}
//
//	err = testRedisSentinel.AddStreamPort(key,mem2)
//
//	assert.Equal(t, nil, err, "Fault to Add Stream2")
//
//	err = testRedisSentinel.RedisClient.LRange(key,0,-1).Err()
//
//	assert.Equal(t, nil, err, "Fault to Got Streams")
//
//	err = testRedisSentinel.DeleteStreamPort(key,0, mem1)
//
//	assert.Equal(t, nil, err, "Fault to Deleted Streams1")
//
//	err = testRedisSentinel.DeleteStreamPort(key,0, mem2)
//
//	assert.Equal(t, nil, err, "Fault to Deleted Streams2")
//}
