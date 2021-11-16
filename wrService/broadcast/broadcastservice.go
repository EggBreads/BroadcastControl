package broadcast

import (
	"errors"
	"github.com/catenoid-company/wrController/config"
	"github.com/catenoid-company/wrController/logger"
	"github.com/catenoid-company/wrController/utils"
	"strconv"
	"strings"
	"time"
)

type ManagePort struct {
	StreamPluginId	int			`json:"stream_plugin_id"`
	AudioPort		int			`json:"audio_port"`
	VideoPort		[]int		`json:"video_port"`
}

/*
	포트관리 redis raw 형태
	Key webrtc_stream_port_manage
	val [
		{"stream_plugin_id":9010,"audio_port":9010,"video_port":[9012,9014,9016]}
	]
*/

/*
	*채널키는 반드시 존재
	@Port 관리 Version 1
	신규 port 생성 방법
	* Port 가 하나도 없을경우 초기값으로 생성
	* 생성된 Port 가 있을경우 가장 마지막 port 번호 이후 번호로 생성
 */

func getStreamAvailablePort(videoPortCnt int, webrtcKey string) (*ManagePort, error) {
	// 초기값 고정
	managePort := &ManagePort{
		StreamPluginId: 9000,
		AudioPort: 9000,
	}

	// StreamPluginId가 존재하는지 확인
	streamPluginInfo := utils.RedisClient.Get(config.WrConfig.StreamPluginIdKey).Val()

	// 마지막 StreamPluginId 조회
	streamPluginId, err := getUseLastPort(streamPluginInfo)
	if err !=nil {
		return nil, err
	}

	// StreamPluginId 정보가 있는경우만 조회된 StreamPluginId 사용
	if streamPluginId > 0 {
		// 마지막 StreamPluginId가 9000 같거나 크면 새로 갱신
		managePort.StreamPluginId = streamPluginId + 10
		managePort.AudioPort = streamPluginId+ 10
		managePort.VideoPort = make([]int, 0)
	}

	videoPort := managePort.StreamPluginId

	for i := 0; i < videoPortCnt; i++ {
		videoPort = videoPort +2
		managePort.VideoPort = append(managePort.VideoPort, videoPort)
	}

	//update last stream_plugin_id
	lastPort := strconv.Itoa(managePort.StreamPluginId)
	val := map[string]interface{}{
		webrtcKey: lastPort,
	}

	byt, err := utils.BytesFromObject(val)
	if err != nil {
		return nil, err
	}

	err = utils.RedisClient.Set(config.WrConfig.StreamPluginIdKey, string(byt), time.Duration(30 * 24)*time.Hour).Err()

	if err != nil {
		return nil, err
	}

	return managePort, nil
}

func getUseLastPort(streamPluginInfo string) (int, error) {
	streamPluginId := 0

	// 마지막 StreamPluginId 조회시 StreamPluginId 할당
	if streamPluginInfo != "" {
		r := strings.NewReader(streamPluginInfo)
		val := map[string]interface{}{}

		_, err := utils.JsonToReader(&val, r)
		if err != nil {
			return  streamPluginId, err
		}

		for _, v := range val{
			streamPluginId, err = strconv.Atoi(v.(string))
			if err != nil {
				return  streamPluginId, err
			}
		}

		if streamPluginId == 0 {
			return streamPluginId,errors.New("Fail to created stream_plugin_id\n")
		}

		return streamPluginId, nil
	}

	// 마지막 StreamPluginId 조회시 StreamPluginId 없을경우 현재 사용중인 StreamPluginId로 할당
	janusInfo, err := utils.GetAllMembers(config.WrConfig.JanusServerInfoKey)
	//logger.Info("onAir StreamPluginId Info : %s", streamPluginInfo)
	if err != nil {
		return streamPluginId, err
	}

	for _, v := range janusInfo {
		info, err := utils.ConvertMapFromStruct(v, "")
		if err != nil {
			return streamPluginId ,err
		}

		// janus field 가 없음
		onJanus, ok := info["janus"].(map[string]interface{})
		if !ok {
			if onJanus == nil{
				//logger.Warn("Empty used janusInfo : %s", info)
				continue
			}

			return streamPluginId, errors.New("Janus type error\n")
		}

		if len(onJanus) > 0 && onJanus != nil{
			streamings, ok := onJanus["streamings"].([]interface{})
			if !ok {
				// janusinfo 에 streamings filed 가 없음
				if streamings == nil{
					//logger.Warn(rId, "Empty used streampluginId : %s", info["janus"])
					continue
				}

				//logger.Error(rId, "Janus type error : %s", info["janus"])
				return streamPluginId, errors.New("Streamings type error\n")
			}

			if len(streamings) > 0 && streamings != nil {
				for _, m := range streamings{
					streaming := m.(map[string]interface{})

					id := int(streaming["id"].(float64))
					if id > streamPluginId{
						streamPluginId = id
					}
				}
			}
		}
	}

	// StreamPluginId가 사용중인 부분이 없고, 마지막 StreamPluginId도 제거된상황
	// 최초 초기화된 처음상황
	if streamPluginId == 0 {
		return streamPluginId, nil
	}

	return streamPluginId, nil
}

//func getStreamAvailablePort(videoPortCnt int, webrtcKey string) (*ManagePort, error) {
//	// 초기값 고정
//	managePort := &ManagePort{
//		StreamPluginId: 9000,
//		AudioPort: 9000,
//	}
//
//	// StreamPluginId가 존재하는지 확인
//	cnt,err := utils.RedisClient.Exists(config.WrConfig.StreamPluginIdKey).Result()
//
//	if err != nil {
//		return nil, err
//	}
//
//	// 존재할경우 마지막번호의 다음번호 생성
//	if cnt > 0{
//		// 마지막 stream_plugin_id 조회
//		result, err := utils.GetAllMembers(config.WrConfig.StreamPluginIdKey)
//
//		if err != nil {
//			return nil, err
//		}
//
//		streamPluginId := 0
//
//		for _, v := range result{
//			streamPluginId, _ = strconv.Atoi(v)
//		}
//
//		if streamPluginId == 0 {
//			return nil,errors.New("Fail to created stream_plugin_id\n")
//		}
//
//		logger.WithField("stream_plugin_id:", streamPluginId).Debug("Get last stream_plugin_id")
//
//		if streamPluginId >= 9000 {
//			managePort.StreamPluginId = streamPluginId + 10
//			managePort.AudioPort = streamPluginId+ 10
//			managePort.VideoPort = make([]int, 0)
//		}
//	}
//
//	videoPort := managePort.StreamPluginId
//
//	for i := 0; i < videoPortCnt; i++ {
//		videoPort = videoPort +2
//		managePort.VideoPort = append(managePort.VideoPort, videoPort)
//	}
//
//	//expiredTime := time.Duration(21600)* time.Second
//
//	//update last stream_plugin_id
//	err = utils.AddRedisMember(config.WrConfig.StreamPluginIdKey, webrtcKey, strconv.Itoa(managePort.StreamPluginId))
//
//	if err != nil {
//		return nil, err
//	}
//
//	return managePort, nil
//}

/*
	Janus Agent 의 Health Check 을 위해사용
	저장된 Janus Agent 에 Ping 을 보내 pong 이 오지않을시 해당 IP 정보는 제거
 */

func GetServerHosts(info map[string]string)[]string  {
	 hosts := make([]string, 0)

	for k := range info{
		hosts = append(hosts, k)
	}

	return hosts
}

/*
	Nginx 종료 요청에 따른 처리
	해당 Channel 에 대한 키 조회 및 제거
	예외 발생은 Webrtc Key 제거가 안될경우에만 error return 이 가야됨
	stream_plugin_id 관련 예외시 Log 만 출력[어차피 신규 방송은 갱신되기 때문에 Log 만 출력해도됨]
*/
func ManageLastStreamPluginId(info map[string]string, rId string) error  {
	// Find last stream_plugin_id
	spi := utils.RedisClient.Get(config.WrConfig.StreamPluginIdKey).Val()

	// 마지막 StreamPluginId가 이미 방송종료되어 제거되었기때문에 Log 만 남김
	if spi == ""{
		logger.Warn("[%s]Empty Last StreamPluginId information :%", rId, info)
		return nil
	}
	// Get to saved streamPluginId
	r := strings.NewReader(spi)
	val := map[string]interface{}{}

	_, err := utils.JsonToReader(&val, r)
	if err != nil {
		return err
	}

	streamPluginId := ""
	for _,v := range val{
		strV, ok := v.(string)
		// 혹시 Type 문제가 발생하더라도 신규로 갱신되면 해결가능하기때문에
		// 로그만 출력
		if !ok {
			logger.Warn("[%s]Removed StreamPluginId type error %s ", rId, v)
			return nil
		}
		streamPluginId = strV
	}

	// 마지막 stream_plugin_id 와 현재 종료시킬 저장된 id가 같은경우 'webrtc_stream_plugin_id'을 제거
	// Janus 에 실패된건이 없을경우 처리
	if streamPluginId == info["stream_plugin_id"] {
		err := utils.DeleteRedisKey(config.WrConfig.StreamPluginIdKey)
		if err != nil {
			return err
			//logger.Error(err.Error())
		}

		// Port 가 없음
		logger.Info(rId, "Remove stream_plugin_id: %s", spi)
	}

	return nil
}
