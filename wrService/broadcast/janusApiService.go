package broadcast

import (
	"errors"
	"fmt"
	"github.com/catenoid-company/wrController/config"
	"github.com/catenoid-company/wrController/logger"
	"github.com/catenoid-company/wrController/models"
	"github.com/catenoid-company/wrController/utils"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type JanusApi struct {
	ChannelKey    string
	StreamKey     string
	RId           string
	Authorization string
	Hosts         []string
}

/*
	Janus Agent Open BroadCast
	janusServer 들에게 Port을 열어달라고 요청하는 서비스
*/
func (ja *JanusApi) OpenJanusBroadCastPorts(params map[string]interface{}) (map[string]interface{}, error) {
	//logger.Info(ja.RId, "JanusChannel receive to janusParams from LiveService : %s", params)
	logger.WithField("janusParams:", params).Info(ja.RId, "JanusChannel receive to janusParams from LiveService")

	wg := &sync.WaitGroup{}
	wg.Add(len(ja.Hosts))

	s := make([]map[string]string, 0)
	f := make([]string, 0)
	//resServers := make([]map[string]interface{}, 0)

	for _, h := range ja.Hosts {
		host := h
		data := params
		go ja.getJanusServersIp(wg, host, data, &s, &f)
	}

	// goroutine timeout check
	isTimeout := utils.WaitTimeout(wg)

	if isTimeout {
		if len(s) == 0 {
			//logger.Error("Janus Servers were timeout")
			return params, errors.New("Janus Servers were timeout\n")
		}
	}

	if len(f) > 0 {
		logger.Warn(ja.RId, "[Start]Fail to Janus Servers : %s", f)
	}

	logger.Info(ja.RId, "[Start]Success Janus Servers : %s", s)
	//logger.Info(ja.RId, "JanusApi complete to request")

	params["servers"] = "0"

	if len(s) == 0 {
		//logger.Error("Empty Janus Servers")
		return nil, errors.New("Empty Janus Servers\n")
	}

	params["servers"] = s

	return params, nil
}

func (ja *JanusApi) getJanusServersIp(wg *sync.WaitGroup, host string, params map[string]interface{}, servers *[]map[string]string, UnresponsiveServers *[]string) {
	defer wg.Done()

	val := make(map[string]string, 0)

	api := utils.WebrtcApi{
		Host:   host,
		Method: "POST",
		Data:   params,
		Headers: map[string]string{
			"X-Request-Id": ja.RId,
			//"Authorization": ja.Authorization,
			"Authorization": utils.AuthorizationHeader(config.WrConfig.AuthUser, config.WrConfig.AuthPass),
		},
	}

	// Janus Api Parameter
	result := models.JanusStreamRes{
		ErrCode: 99,
	}

	res, err := api.CallApi(ja.RId, "v1/streaming-plugin", nil)

	if err != nil {
		logger.Error(ja.RId, "Not response to start janus /streaming-plugin : %s", err.Error())
		*UnresponsiveServers = append(*UnresponsiveServers, host)
		return
	}

	b, err := utils.JsonToReader(&result, res.Body)

	if err != nil {
		logger.Error(ja.RId, "%s", err.Error())
		return
	}

	logger.Info(ja.RId, "Status : %s, ResponseBody : %s", res.Status, string(b))

	if result.ErrCode == 0 {
		// 변경되는 부분이 Janus Host Heath check 가 들어가면 필요없을수 있음
		// 필요없을경우는 Success 된 host 정보로 변경
		val["server_ip"] = result.Data.ServerIp
		*servers = append(*servers, val)
		//logger.WithStruct(res).Info(ja.RId, "Response receive to janusServer")
	} else {
		logger.WithStruct(res).Error(ja.RId, "Failed to Janus Response")
	}
}

// TODO 방송종료시 성공건에 대한 StreamPluginId 정보를 제거하자
func (ja *JanusApi) CloseToJanusServer(wg *sync.WaitGroup, params map[string]interface{}, successHosts *[]string, UnresponsiveServers *[]string) {
	defer wg.Done()

	api := utils.WebrtcApi{
		Host:   params["host"].(string),
		Method: "DELETE",
		Data:   params,
		Headers: map[string]string{
			"X-Request-Id":  ja.RId,
			"Authorization": ja.Authorization,
			//"Authorization": utils.AuthorizationHeader(config.WrConfig.AuthUser, config.WrConfig.AuthPass),
		},
	}

	// Janus Api Parameter
	result := make(map[string]interface{}, 0)

	res, err := api.CallApi(ja.RId, "v1/streaming-plugin", nil)
	if err != nil {
		logger.Error(ja.RId, "Not response to close janus /streaming-plugin : %s", err.Error())
		*UnresponsiveServers = append(*UnresponsiveServers, api.Host)
		return
	}

	b, err := utils.JsonToReader(&result, res.Body)

	logger.Info(ja.RId, "Status : %s, ResponseBody : %s", res.Status, string(b))

	if err != nil {
		logger.Error(ja.RId, "%s", err.Error())
		return
	}

	if res.StatusCode == http.StatusOK {
		*successHosts = append(*successHosts, api.Host)
	} else {
		logger.Error(ja.RId, "Failed to response. %s", res)
	}
}

/*
	성공한 StreamPluginId 기준으로 해당 Hosts 갯수 만큼 저장
*/
func (ja *JanusApi) SaveSuccessResJanusHosts(key string, broadcastKey string, janusParams map[string]interface{}, servers []map[string]string) error {
	successHosts := make([]string, 0)

	for _, v := range servers {
		successHosts = append(successHosts, v["server_ip"])
	}

	streampluginId := janusParams["stream_plugin_id"].(int)
	strStreamPluginId := strconv.Itoa(streampluginId)

	memMap := map[string]interface{}{
		strStreamPluginId: strings.Join(successHosts, ","),
	}

	err := utils.AddRedisMembers(config.WrConfig.StreamPluginIdOnAirKeys, memMap)
	if err != nil {
		logger.Error(ja.RId, "OnAir Broadcast wasn't saved : %s", err.Error())

		// 만약 Redis 에 저장에 실패된다면 방송시작된 Janus 을 종료해줘야된다
		sErr := ja.StopJanusBroadCastOnStart(broadcastKey, key, janusParams)
		if sErr != nil {
			return sErr
		}

		return err
	}

	return nil
}

/*
	방송 시작중 실패시 방송 종료 메세지 처리
	시작된 Janus 방송 종료 처리
*/
func (ja *JanusApi) StopJanusBroadCastOnStart(broadcastKey string, key string, closeJanusParams map[string]interface{}) error {
	liveApi := LiveApi{
		RId:           ja.RId,
		Authorization: ja.Authorization,
	}

	err := liveApi.CloseToLiveBroadCast(broadcastKey)

	if err != nil {
		//logger.Error(ja.RId, "Invalid broadcastKey %s : %s", broadcastKey, err.Error())
		return err
	}

	err = utils.DeleteRedisKey(key)
	if err != nil {
		return err
	}

	if closeJanusParams != nil {
		hosts, ok := closeJanusParams["hosts"].([]map[string]string)
		if !ok {
			return errors.New(fmt.Sprintf(config.FORMAT_UNRESPONSIVE_ALL_HOST_API_OR_NON_DATA, config.JANUS))
		}

		channelKey := closeJanusParams["channel_key"]
		streamPluginId := closeJanusParams["stream_plugin_id"]

		// 방송종료 실패시 정상처리된 JANUS AGENT 에 종료를 전달
		if len(hosts) > 0 {
			s := make([]string, 0)
			f := make([]string, 0)

			wg := &sync.WaitGroup{}
			wg.Add(len(hosts))

			for _, h := range hosts {
				janusParams := make(map[string]interface{}, 0)

				janusParams["channel_key"] = channelKey
				janusParams["stream_plugin_id"] = streamPluginId
				janusParams["host"] = h["server_ip"]
				go ja.CloseToJanusServer(wg, janusParams, &s, &f)
			}

			logger.Warn(ja.RId, "BroadCast isn't start by liveApiStream")

			// Timeout 됬어도 실패된 Host 정보들은 가지고있기때문에 로그만 출력
			isTimeout := utils.WaitTimeout(wg)
			if isTimeout {
				logger.Warn(ja.RId, "Janus Servers were timeout")
			}

			// TODO 성공된 Janus 들을 종료시키는중 발생되여 할당이 되기때문에 실패 StreamPluginId을 추가해야함
			if len(f) > 0 {
				logger.Warn(ja.RId, "Fail to Janus Servers : %s", f)
			}

			logger.Warn(ja.RId, "Success Janus Servers : %s", s)
		}
	}
	return nil
}
