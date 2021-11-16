package monitoring

import (
	"errors"
	"github.com/catenoid-company/wrController/config"
	"github.com/catenoid-company/wrController/logger"
	"github.com/catenoid-company/wrController/utils"
	"net/http"
	"strconv"
	"strings"
)

type JanusMonitoringService struct {
	RId           string
	Authorization string
}

/*
	JanusServer 모니터링 정보 저장
	RedisKey : "webrtc_janus_servers_info"
	Member:	각각의 Janus IP
	Value:	janus server info
	결과
	ex) "1.2.3.4" : {Janus 정보들}
*/
func (jms *JanusMonitoringService) SaveJanusMonitoring(janusInfo map[string]interface{}) error {
	logger.MonitoringLogger.WithField("rid", jms.RId).Info("Start JanusMonitoringService")

	return saveComputeInfo(janusInfo, config.WrConfig.JanusServerInfoKey, jms.RId)
}

func (jms *JanusMonitoringService) GetMonitoring() (map[string]interface{}, error) {
	logger.MonitoringLogger.WithField("rid", jms.RId).Info("Start to Get Janus Monitoring Info")
	return getMonitoring(jms.RId, config.WrConfig.JanusServerInfoKey)
}

/*
	Janus Server Allow Time check
	허용된 시간을 넘을경우 해당 Server 정보는 제거
*/
func (jms *JanusMonitoringService) JanusRefreshAllowTime(server string) error {
	isJanusServer := utils.RedisClient.HExists(config.WrConfig.JanusServerInfoKey, server).Val()

	if !isJanusServer {
		// 정보없음
		return errors.New("Not Match Server On JanusRedis\n")
	}

	info, err := getServerMonitoring(jms.RId, config.WrConfig.JanusServerInfoKey, server)

	if err != nil {
		return err
	}

	if len(info) == 0 {
		// 정보없음
		return errors.New("Empty Server Information On JanusRedis\n")
	}

	err = updateAllowTime(config.WrConfig.JanusServerInfoKey, server, info)

	if err != nil {
		return err
	}

	return nil
}

/*
	실제 방송을 종료되였지만, Janus Agent 에 방송 정보가 남아있을경우 제거
	모니터링 정보들은 저장되어야 함으로 Return 은 의미가 없으며, Log 로 출력해야됨
*/
func (jms *JanusMonitoringService) CloseToRemainJanusBroadCast(streamings []map[string]interface{}, server string) []interface{} {
	result := make([]interface{}, 0)

	// 현재 Streaming 중인 StreamPluginId가 있는 Janus
	for _, streaming := range streamings {
		result = append(result, streaming)
		/*
			TODO 변경점
			TODO 1.채널키를 기준으로 Redis Key 을 찾은후 O
			TODO 2.해당 Redis Key 을 조회
			TODO 2-1 조회시 Key 가 존재할경우 정상처리
			TODO 2-2 조회시 Key 가 존재하지 않을경우 Janus 종료처리
			TODO 3. 방송중인 StreamPluginId을 조회
			TODO 3-1 동일할경우 그대로 진행
			TODO 3-2 서로 다를경우 onAir Key 에서 해당 StreamPluginId 제거

			TODO 미해결
			TODO 방송중인지 확인할수있는 방법이 키를 추가하지 않으면 방법이 없음
		*/

		des, ok := streaming["description"].(string)
		if !ok {
			logger.MonitoringLogger.WithField("rid", jms.RId).Warnf("Non parsed description (server %s) ", server)
			continue
		}

		channelKey := strings.ReplaceAll(des, "Channel-", "")

		logger.MonitoringLogger.WithField("rid", jms.RId).Infof("ChannelKey : %s", channelKey)

		webrtcBroadCastKey := utils.RedisClient.Keys("webrtc_" + channelKey + "_*").Val()

		logger.MonitoringLogger.WithField("rid", jms.RId).WithField("webrtcBroadCastKey", webrtcBroadCastKey).Info("Find to key")

		id, ok := streaming["id"].(float64)
		if !ok {
			logger.MonitoringLogger.WithField("rid", jms.RId).Warnf("Non parsed streampluginId (server %s) ", server)
			continue
		}

		streamPluginId := strconv.Itoa(int(id))

		// Key 가 존재하지 않는경우는 방송이 종료된것으로 Janus 에 종료처리를 해줌
		// StreamPluginIdOnAirKeys 의 StreamPlugInId 정보는 새로 생성된
		// 방송에 의해서 할당되어 사용되었을수도 있으므로 제거안함
		if len(webrtcBroadCastKey) == 0 {
			err := jms.sendCloseToJanus(server, channelKey, int(id))
			if err != nil {
				logger.MonitoringLogger.WithField("rid", jms.RId).Error(err.Error())
			}
			continue
		}

		webrtcInfo := utils.RedisClient.HGetAll(webrtcBroadCastKey[0]).Val()

		logger.MonitoringLogger.WithField("webrtcInfo", webrtcInfo).Info("OnAir Info")

		findStreamPluginId := webrtcInfo["stream_plugin_id"]

		// StreamPluginId가 서로 다르면 onAir의 PluginId가 변경되었다는것이므로
		// Janus 쪽에 해당 StreamPluginId을 종료 요청을 해주어야됨
		if streamPluginId != findStreamPluginId {
			err := jms.sendCloseToJanus(server, channelKey, int(id))
			if err != nil {
				logger.MonitoringLogger.WithField("rid", jms.RId).Error(err.Error())
			}
			continue
		}
	}

	logger.MonitoringLogger.WithField("rid", jms.RId).WithField("result", result).Infof("Janus Streaming Info")

	return result
}

// Janus 방송종료
func (jms *JanusMonitoringService) sendCloseToJanus(server, channelKey string, streamPluginId int) error {
	api := utils.WebrtcApi{
		Host:   server,
		Method: "DELETE",
		Data: map[string]interface{}{
			"channel_key":      channelKey,
			"stream_plugin_id": streamPluginId,
		},
		Headers: map[string]string{
			"X-Request-Id":  jms.RId,
			"Authorization": jms.Authorization,
			//"Authorization": utils.AuthorizationHeader(config.WrConfig.AuthUser, config.WrConfig.AuthPass),
		},
	}

	res, err := api.CallApi(jms.RId, "v1/streaming-plugin", nil)

	if err != nil {
		logger.MonitoringLogger.WithField("rid", jms.RId).Errorf("Failed to got stream score : %s", err.Error())
		return err
	}

	b, err := utils.JsonToReader(nil, res.Body)

	if err != nil {
		logger.Error(jms.RId, "%s", err.Error())
		return err
	}

	logger.Info(jms.RId, "Status : %s, ResponseBody : %s", res.Status, string(b))

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		logger.MonitoringLogger.WithField("rid", jms.RId).Errorf("Failed to janus broadcast closed  %s(server %s) ", streamPluginId, server)
		return err
	}

	return nil
}
