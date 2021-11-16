package broadcast

import (
	"github.com/stretchr/testify/assert"
	"testing"
)
//
////var testConfig *config.Configuration
////var testRedisSentinel *lib.RedisSentinel
////
////func Test_Init(t *testing.T) {
////	testConfig = config.InitConfig()
////	testRedisSentinel = &lib.RedisSentinel{}
////	testRedisSentinel.RedisClient = testRedisSentinel.ConnectSentinel(testConfig)
////}
//
func Test_OpenJanusBroadCastPorts(t *testing.T) {
	api := &JanusApi{}

	params := map[string]interface{}{
		"channel_key": "wvu94hscinf5dsa8",
		"stream_plugin_id": 9000,
		"audio_port": 9000,
		"video_port1": 9002,
		"video_port2": 9004,
	}

	result,_ := api.OpenJanusBroadCastPorts(params)

	expect := map[string]interface{}{
		"channel_key": "wvu94hscinf5dsa8",
		"stream_plugin_id": 9000,
		"audio_port": 9000,
		"video_port1": 9002,
		"video_port2": 9004,
		"servers":[]map[string]string{
			{"ip":"1.2.3.4"},
			{"ip":"1.2.3.4"},
			{"ip":"1.2.3.4"},
		},
	}
	assert.Equal(t, expect, result, "Fault to Checked Janus Server")
}

//func Test_CloseToJanusServer(t *testing.T) {
//	wg := &sync.WaitGroup{}
//
//	api := &JanusApi{
//		RedisSentinel: testRedisSentinel,
//	}
//
//	params := map[string]interface{}{
//		"channel_key": "wvu94hscinf5dsa8",
//		"stream_plugin_id": 9000,
//		"audio_port": 9000,
//		"video_port1": 9002,
//		"video_port2": 9004,
//	}
//
//	go api.CloseToJanusServer(wg,params)
//
//}
