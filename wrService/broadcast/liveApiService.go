package broadcast

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/catenoid-company/wrController/config"
	"github.com/catenoid-company/wrController/logger"
	"github.com/catenoid-company/wrController/models"
	"github.com/catenoid-company/wrController/utils"
	"net/http"
	"strconv"
)

type LiveApi struct {
	ChannelKey    string
	RId           string
	Authorization string
	//Exception		*utils.ExceptionProcess
}

/**
방송시작
LiveApi 의 profiles 의 정보를 가져옴
*/
func (la *LiveApi) GetBroadCastProfiles(nginxParams models.NginxParameters, key string) (*models.BroadCastProfiles, *ManagePort, map[string]interface{}, error) {
	//defer la.Exception.ExceptionRecovery("GetProfileException")

	// Convert to map from parameters
	params := make(map[string]interface{})

	byt, err := utils.BytesFromObject(nginxParams)
	if err != nil {
		return nil, nil, nil, err
	}
	r := bytes.NewReader(byt)

	_, err = utils.JsonToReader(&params, r)
	if err != nil {
		return nil, nil, nil, err
	}

	// Get to Profiles by Live Api
	liveProfiles, err := la.getProfilesFromLiveApi(params)

	if err != nil {
		return nil, nil, nil, err
	}

	//webrtcKey := "webrtc_"+ nginxParams.ChannelKey+ "_"+ liveProfiles.BroadcastKey

	// stream_plugin_id 생성
	availablePort, err := getStreamAvailablePort(len(liveProfiles.VideoProfiles), key)

	if err != nil {
		return liveProfiles, nil, nil, err
	}

	logger.WithField("availablePort:", availablePort).Info(la.RId, "Find to use ports")

	// webrtc_{channel_key}_{broadcast_key} 의 멤버에 추가
	params["stream_plugin_id"] = availablePort.StreamPluginId

	// Set to janus parameters
	janusParams := map[string]interface{}{}

	janusParams["channel_key"] = nginxParams.ChannelKey
	janusParams["broadcast_key"] = nginxParams.BroadcastKey
	janusParams["stream_plugin_id"] = availablePort.StreamPluginId
	janusParams["audio_port"] = availablePort.AudioPort

	for i, port := range availablePort.VideoPort {
		key := "video_port" + strconv.Itoa(i+1)
		janusParams[key] = port
	}

	// Key 의 정보 저장
	err = la.saveNginxInfo(key, params)

	if err != nil {
		return liveProfiles, availablePort, janusParams, err
		//errors.New("Don't save to nginx info\n")
	}

	return liveProfiles, availablePort, janusParams, nil
}

func (la *LiveApi) saveNginxInfo(key string, values map[string]interface{}) error {
	logger.WithField("nginxParams:", values).Info(la.RId, "BroadCast open request params from nginx agent")

	// Nginx 의 정보 저장
	err := utils.AddRedisMembers(key, values)

	if err != nil {
		return err
	}
	logger.Info(la.RId, "Save to NginxParameters by redis")

	return nil
}

func (la *LiveApi) getProfilesFromLiveApi(values map[string]interface{}) (*models.BroadCastProfiles, error) {
	// Get to Profiles by Live Api
	api := utils.WebrtcApi{
		Host: config.WrConfig.LiveApiHost,
		Data: values,
		Headers: map[string]string{
			"X-Request-Id":  la.RId,
			"Authorization": la.Authorization,
			//"Authorization": utils.AuthorizationHeader(config.WrConfig.AuthUser, config.WrConfig.AuthPass),
		},
	}

	liveProfiles := &models.BroadCastProfiles{}

	res, err := api.CallApi(la.RId, "api/v1/live/channels/{channel_key}/profiles", map[string]string{"channel_key": values["channel_key"].(string)})
	if err != nil {
		return nil, err
	}

	b, err := utils.JsonToReader(liveProfiles, res.Body)

	if err != nil {
		return nil, err
	}

	logger.Info(la.RId, "Get profiles response body : %s (%d)", string(b), res.StatusCode)

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(string(b))
	}

	return liveProfiles, nil
}

func (la *LiveApi) SendToOpenFromLiveApi(liveProfiles *models.BroadCastProfiles, availablePort *ManagePort, liveParams map[string]interface{}) (*models.NginxBroadCastOpenRes, error) {
	// Nginx 응답값 처리
	nginxRes := la.generateNginxRes(liveProfiles, availablePort, liveParams)

	logger.WithField("nginxRes:", nginxRes).Debug(la.RId, "Nginx Response Json is filled")

	// liveApi janus 서버정보 응답처리
	liveOpenRes, err := la.sendJanusServerInfoFromLiveApi(liveParams)

	if err != nil {
		return nil, err
	}

	//logger.Info("[%s]LiveApi request to liveOpenParams : %s", la.RId, liveOpenRes)

	logger.WithField("liveOpenRes:", liveOpenRes).Info(la.RId, "LiveChannel pass on nginxRes to LiveService")

	return nginxRes, nil
}

/*
	nginx Response 생성
*/

func (la *LiveApi) generateNginxRes(liveProfiles *models.BroadCastProfiles, availablePort *ManagePort, liveParams map[string]interface{}) *models.NginxBroadCastOpenRes {
	// Format to VideoProfiles
	videoProfiles := make([]map[string]interface{}, 0)

	// VideoProfile 을 nginx agent Res 값으로 변환
	for j := 0; j < len(liveProfiles.VideoProfiles); j++ {
		profile := liveProfiles.VideoProfiles[j]

		var videoProfile map[string]interface{}

		videoProfile, _ = utils.ConvertMapFromStruct(profile, "json")
		videoProfile["video_port"] = availablePort.VideoPort[j]
		videoProfiles = append(videoProfiles, videoProfile)
	}

	// AudioProfile 을 nginx agent Res 값으로 변환
	var audioProfile map[string]interface{}

	audioProfile, _ = utils.ConvertMapFromStruct(liveProfiles.AudioProfile, "json")
	audioProfile["audio_port"] = availablePort.AudioPort

	// NginxAgent response 값
	nginxRes := &models.NginxBroadCastOpenRes{
		ChannelKey:     liveParams["channel_key"].(string),
		VideoProfiles:  videoProfiles,
		AudioProfile:   audioProfile,
		StreamPluginId: availablePort.StreamPluginId,
		Servers:        liveParams["servers"].([]map[string]string),
		ListenPort:     availablePort.StreamPluginId,
	}

	return nginxRes
}

/*
	LiveApi 에 Janus 정보 전송
*/
func (la *LiveApi) sendJanusServerInfoFromLiveApi(liveParams map[string]interface{}) (*models.BroadCastOpenRes, error) {
	// liveApi Server 정보 Request Parameter
	liveOpenReqParams := make(map[string]interface{}, 0)
	janusServers := make([]map[string]interface{}, 0)
	janusOpenServers := liveParams["servers"].([]map[string]string)

	channelKey := liveParams["channel_key"].(string)
	broadcastKey := liveParams["broadcast_key"].(string)

	serverInfo, _ := utils.ConvertMapFromStruct(liveParams, "")
	delete(serverInfo, "servers")
	delete(serverInfo, "channel_key")
	delete(serverInfo, "broadcast_key")

	for _, v := range janusOpenServers {
		serverInfo["ip"] = v["server_ip"]
		//serverInfo["ip"] = "1.2.3.4"
		janusServers = append(janusServers, serverInfo)
	}

	//liveOpenReqParams["channel_key"] = liveParams["channel_key"]
	liveOpenReqParams["servers"] = janusServers

	// Get to Profiles by Live Api
	api := utils.WebrtcApi{
		Host:   config.WrConfig.LiveApiHost,
		Method: "POST",
		Data:   liveOpenReqParams,
		Headers: map[string]string{
			"X-Request-Id":  la.RId,
			"Authorization": la.Authorization,
			//"Authorization": utils.AuthorizationHeader(config.WrConfig.AuthUser, config.WrConfig.AuthPass),
		},
	}

	liveOpenRes := &models.BroadCastOpenRes{}

	res, err := api.CallApi(la.RId,
		"api/v1/live/channels/{channel_key}/broadcasts/{broadcast_key}/streaming-plugin",
		map[string]string{"channel_key": channelKey, "broadcast_key": broadcastKey})
	if err != nil {
		return nil, err
	}

	//logger.Info(la.RId, "Status : %s, ResponseBody : %s", res.Status, string(b))
	logger.Info(la.RId, "Status : %s", res.Status)

	if res.StatusCode != http.StatusNoContent {
		b, err := utils.JsonToReader(nil, res.Body)
		if err != nil {
			return nil, err
		}

		return nil, errors.New("No send request from live kollus : " + string(b))
	}

	return liveOpenRes, nil
}

/*
	LiveApi 에 방송 시작 전송
*/
func (la *LiveApi) SendOpenSignalToLivaApi(params models.NginxParameters) (map[string]interface{}, error) {
	//streamRes := &models.BroadCastStreamRes{}

	api := utils.WebrtcApi{
		Host:   config.WrConfig.LiveApiHost,
		Method: "PUT",
		Data: map[string]interface{}{
			"streaming_server_type": "webrtc",
			"streaming_server_host": params.Host,
			"streaming_server_ip":   params.Server,
			"shooter_ip":            params.Client,
		},
		Headers: map[string]string{
			"X-Request-Id":  la.RId,
			"Authorization": la.Authorization,
			"Content-Type":  "application/x-www-form-urlencoded",
		},
	}

	// BroadCast Stream 으로 방송 시작을 보냄
	res, err := api.CallApi(la.RId, "api/v1/live/broadcasts/start/{stream_key}", map[string]string{"stream_key": params.StreamKey})
	if err != nil {
		return nil, err
	}

	body := make(map[string]interface{}, 0)

	b, err := utils.JsonToReader(&body, res.Body)
	if err != nil {
		return nil, err
	}

	logger.Info(la.RId, "Live Stream Start Api Response Body : %s (%d)", string(b), res.StatusCode)

	// 방송 Live 시작 요청 파라미터가 정상적으로 처리가 안된경우
	if res.StatusCode == http.StatusBadRequest {
		return nil, errors.New(fmt.Sprintf("%d : %s", config.INVALID_PARAMS, config.INVALID_PARAMETERS))
	}

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("%d : %s", config.API_OR_NON_DATA, config.FORMAT_ALREADY_ONAIR))
	}

	body, ok := body["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New(fmt.Sprintf("%d : %s", config.CONTROLLER_EXCEPTION, config.FORMAT_FAIL_PARSE_BYTES_CONTROLLER_EXCEPTION))
	}

	return body, nil
}

/*
	Send to Live Stream Api
	방송이 종료되었다는 시그널을 Stream Api 에 전달
*/
func (la *LiveApi) CloseToLiveBroadCast(broadCastKey string) error {
	api := utils.WebrtcApi{
		Host:   config.WrConfig.LiveApiHost,
		Method: "PUT",
		Headers: map[string]string{
			"X-Request-Id":  la.RId,
			"Authorization": la.Authorization,
			"Content-Type":  "application/x-www-form-urlencoded",
			//"Authorization": utils.AuthorizationHeader(config.WrConfig.AuthUser, config.WrConfig.AuthPass),
		},
	}

	res, err := api.CallApi(la.RId, "api/v1/live/broadcasts/{broadcast_key}/pause", map[string]string{"broadcast_key": broadCastKey})
	if err != nil {
		return err
	}

	logger.WithField("StatusCode:", res.StatusCode).Info(la.RId, "success to closed LiveApi stream broadcastKey: %s (%d)", broadCastKey, res.StatusCode)

	// 응답 실패
	if res.StatusCode != http.StatusNoContent {

		return errors.New("Fail to broadcast stop from liveApiStream\n")
	}

	// 응답 완료
	logger.Info(la.RId, "Success to broadcast stop from liveApiStream")

	return nil

}
