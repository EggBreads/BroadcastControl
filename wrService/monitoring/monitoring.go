package monitoring

import (
	"github.com/catenoid-company/wrController/config"
	"github.com/catenoid-company/wrController/logger"
	"github.com/catenoid-company/wrController/utils"
	"strconv"
	"time"
)

func saveComputeInfo(obj map[string]interface{}, key, rId string) error {
	val := make(map[string]interface{}, 0)

	server := obj["server"].(string)
	delete(obj, "server")

	b, err := utils.BytesFromObject(obj)
	if err != nil {
		return err
	}

	val[server] = string(b)

	logger.MonitoringLogger.WithField("rId", rId).WithField("MonitoringParams:", val).Infof("key : %s", key)

	err = utils.AddRedisMembers(key, val)

	if err != nil {
		return err
	}

	logger.MonitoringLogger.WithField("rId", rId).Infof("Finish to Save %s Monitoring Info", key)
	return nil
}

func getMonitoring(rId, key string) (map[string]interface{}, error) {
	monitoringInfo := make(map[string]interface{}, 0)

	result, err := utils.GetAllMembers(key)

	if err != nil {
		return nil, err
	}

	logger.MonitoringLogger.WithField("rId", rId).WithField("result", result).Infof("%s Monitoring Result", key)

	monitoringInfo, _ = utils.ConvertMapFromStruct(result, "json")

	logger.MonitoringLogger.WithField("rId", rId).Infof("Finish to Get %s Monitoring Info", key)

	return monitoringInfo, nil
}

func getServerMonitoring(rid, key string, server string) (map[string]interface{}, error) {
	m, err := utils.GetMembers(key, server)
	if err != nil {
		return nil, err
	}

	info, _ := utils.ConvertMapFromStruct(m[server], "")

	logger.MonitoringLogger.WithField("rid", rid).WithField("info", info).Infof("Get Server Monitoring Information %s", key)

	return info, nil
}

func AgentManageServers() {
	rId := utils.GetRIdUUID()
	logger.MonitoringLogger.WithField("rId", rId).Info("======== Start to manage hosts ========")

	healthCheckTime, _ := strconv.Atoi(config.WrConfig.JanusHealthCheckTime)

	for {
		logger.MonitoringLogger.WithField("rId", rId).Info("======== Run to manage hosts ========")
		time.Sleep(time.Duration(healthCheckTime) * time.Minute)

		// Janus Info 조회
		janusInfo, err := utils.GetAllMembers(config.WrConfig.JanusServerInfoKey)

		nginxInfo, err := utils.GetAllMembers(config.WrConfig.NginxServerInfoKey)

		if err != nil {
			logger.MonitoringLogger.WithField("rId", rId).Errorf("Stopped to JanusHealthCheck :%s", err.Error())
			continue
		}

		if len(janusInfo) > 0 {
			removeServers := checkLimitTime(janusInfo, rId)

			b, err := utils.BytesFromObject(removeServers)
			if err != nil {
				logger.MonitoringLogger.WithField("rId", rId).Errorf("Fail to converted to bytes : %s", err.Error())
				continue
			}

			logger.MonitoringLogger.WithField("rId", rId).Infof("Remove to janus servers list : %s", string(b))

			if len(removeServers) > 0 {
				delAgentMem(rId, config.WrConfig.JanusServerInfoKey, removeServers...)
			}
		}

		if len(nginxInfo) > 0 {
			removeServers := checkLimitTime(nginxInfo, rId)

			b, err := utils.BytesFromObject(removeServers)
			if err != nil {
				logger.MonitoringLogger.WithField("rId", rId).Errorf("Fail to converted to bytes : %s", err.Error())
				continue
			}

			logger.MonitoringLogger.WithField("rId", rId).Infof("Remove to nginx servers list : %s", string(b))

			if len(removeServers) > 0 {
				delAgentMem(rId, config.WrConfig.NginxServerInfoKey, removeServers...)
			}
		}
		logger.MonitoringLogger.WithField("rId", rId).Info("======== Finish to manage hosts ========")
	}
}

func checkLimitTime(m map[string]string, rId string) []string {
	removeServers := make([]string, 0)

	for k, v := range m {
		denyTime := float64(time.Now().Unix())

		info := map[string]interface{}{}
		info, err := utils.ConvertMapFromStruct(v, "")
		if err != nil {
			logger.MonitoringLogger.WithField("rId",rId).Errorf("Failed to convert map : %s", err.Error())
			continue
		}

		allowTime := info["allow_limit_time"].(float64)

		// Limit Time 없음
		if allowTime == 0 {
			logger.MonitoringLogger.WithField("rId",rId).Errorf("Empty Limit Time. Server %s", k)
			continue
		}

		// 시간내에 사용중인 Server
		if allowTime > denyTime {
			continue
		}

		// 미사용중인 Server 제거하기위해 추가
		removeServers = append(removeServers, k)
	}

	return removeServers
}

func delAgentMem(rId, key string, h ...string) {
	err := utils.RedisClient.HDel(key, h...).Err()

	if err != nil {
		logger.MonitoringLogger.WithField("rId", rId).Errorf("Janus Agent Host is failed to unUsed Host : %s", err.Error())
	}

	logger.MonitoringLogger.WithField("rId", rId).Infof("Janus Host Agent is deleted. Host is broken out: %s", h)
}

func updateAllowTime(key, mapKey string, info map[string]interface{}) error {
	val := make(map[string]interface{}, 0)

	term, _ := strconv.Atoi(config.WrConfig.AllowLimitTimeTerm)
	info["allow_limit_time"] = time.Now().Add(time.Duration(term) * time.Minute).Unix()

	b, err := utils.BytesFromObject(info)
	if err != nil {
		return err
	}

	val[mapKey] = string(b)

	return utils.AddRedisMembers(key, val)
}
