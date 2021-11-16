package monitoring

import (
	"errors"
	"github.com/catenoid-company/wrController/config"
	"github.com/catenoid-company/wrController/logger"
	"github.com/catenoid-company/wrController/utils"
)

type NginxMonitoringService struct {
	RId string
	Authorization	string
}

/*
	JanusServer 모니터링 정보 저장
	RedisKey : "webrtc_janus_servers_info"
	Member:	각각의 Janus IP
	Value:	janus server info
	결과
	ex) "1.2.3.4" : {Janus 정보들}
*/
func (nms *NginxMonitoringService) SaveNginxMonitoring(nginxInfo map[string]interface{}) error {
	logger.MonitoringLogger.WithField("rid" , nms.RId).Infof("Start NginxMonitoringService")
	return saveComputeInfo(nginxInfo, config.WrConfig.NginxServerInfoKey, nms.RId)
}

func (nms *NginxMonitoringService) GetMonitoring() (map[string]interface{}, error){
	logger.MonitoringLogger.WithField("rid" , nms.RId).Infof("Start NginxMonitoringService")
	return getMonitoring(nms.RId, config.WrConfig.NginxServerInfoKey)
}

/*
	Nginx Server Allow Time check
	허용된 시간을 넘을경우 해당 Server 정보는 제거
 */
func (nms *NginxMonitoringService) NginxRefreshAllowTime(server string) error {
	isNginxServer := utils.RedisClient.HExists(config.WrConfig.NginxServerInfoKey, server).Val()

	if !isNginxServer {
		// 정보없음
		return  errors.New("Not Match Server On NginxRedis\n")
	}

	info, err := getServerMonitoring(nms.RId, config.WrConfig.NginxServerInfoKey, server)

	if err != nil {
		return err
	}

	if len(info) == 0 {
		// 정보없음
		return errors.New("Empty Server Information On NginxRedis\n")
	}

	err = updateAllowTime(config.WrConfig.NginxServerInfoKey, server, info)

	if err != nil {
		return err
	}

	return nil
}
