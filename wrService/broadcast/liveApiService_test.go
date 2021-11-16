package broadcast

import (
	"github.com/catenoid-company/wrController/models"
	"github.com/catenoid-company/wrController/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetBroadCastProfiles_SendToOpenFromLiveApi(t *testing.T) {
	api := &LiveApi{
		ChannelKey: "wvu94hscinf5dsa2",
	}

	key := "webrtc_wvu94hscinf5dsa2_20210427-sdfjsnd"

	params := models.NginxParameters{
		Record: false,
		ChannelKey: "jbgi96f5xf7ikhpm",
		Rtmp: "rtmp://localhost:1935/jbgi96f5xf7ikhpm",
		Client: "11.22.11.22",
		Host: "kr01wr01",
		Server: "1.2.3.4",
		BroadcastKey:"20210427-sdfjsnd",
		StreamKey:"abcdefs3",
	}

	r1 := &models.BroadCastProfiles{}
	r2 := &ManagePort{}
	r3 := map[string]interface{}{}

	profiles, ports, janus, _ := api.GetBroadCastProfiles(params,"")

	assert.ObjectsAreEqual(profiles, r1)

	assert.ObjectsAreEqual(ports, r2)

	assert.ObjectsAreEqual(janus, r3)

	liveParams := map[string]interface{}{
		"channel_key" : "wvu94hscinf5dsa2",
		"stream_key" : "asdfgb",
		"stream_plugin_id": 9010,
		"audio_port" : 9010,
		"video_port1" : 9012,
		"video_port2" : 9014,
		"servers" : []map[string]interface{}{
			{"ip": "1.2.3.4"},
			{"ip": "1.2.3.4"},
		},
	}

	res := &models.NginxBroadCastOpenRes{}
	res, err := api.SendToOpenFromLiveApi(profiles, ports, liveParams)

	assert.NoError(t, err, "Fault to LiveApi Open err")

	assert.ObjectsAreEqual(&models.NginxBroadCastOpenRes{}, res)

	webrtcStreamPluginIdKey := "webrtc_stream_plugin_id"
	err = utils.RedisClient.Del(webrtcStreamPluginIdKey, key).Err()

	assert.NoError(t, err, "Fault to keys delete")
}