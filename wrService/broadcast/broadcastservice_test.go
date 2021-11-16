package broadcast

import (
	"github.com/catenoid-company/wrController/config"
	"github.com/catenoid-company/wrController/logger"
	"github.com/catenoid-company/wrController/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Init(t *testing.T) {
	config.WrConfig = config.InitConfig()

	utils.RedisClient, _ = utils.ConnectSentinel(config.WrConfig)
	// Logger Set
	logger.Init()

}

func Test_GetStreamAvailablePort(t *testing.T) {
	key := "webrtc_wvu94hscinf5dsa2_20210303-abcdefg"

	webrtcStreamPluginIdKey := "webrtc_stream_plugin_id"

	m1 := &ManagePort{
		StreamPluginId: 9000,
		AudioPort: 9000,
		VideoPort: []int{9002, 9004, 9006},
	}

	p1 , err := getStreamAvailablePort(3, key)

	assert.NoError(t, err, "Fault to got ports1")

	assert.Equal(t, m1, p1, "Fault to difference managePorts1")

	m2 := &ManagePort{
		StreamPluginId: 9010,
		AudioPort: 9010,
		VideoPort: []int{9012, 9014, 9016},
	}

	p2 , err :=getStreamAvailablePort(3, key)

	assert.NoError(t, err, "Fault to got ports2")

	assert.Equal(t, m2, p2, "Fault to difference managePorts2")

	id, err := utils.RedisClient.Get(webrtcStreamPluginIdKey).Result()

	assert.NoError(t, err, "Fault to got id")

	assert.Equal(t, "9010", id, "Fault to difference id" + id)

	err = utils.DeleteRedisKey(webrtcStreamPluginIdKey)

	assert.NoError(t, err, "Fault to deleted")
}
