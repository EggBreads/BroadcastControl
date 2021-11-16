package broadcast

import (
	"errors"
	"fmt"
	"github.com/catenoid-company/wrController/config"
	"github.com/catenoid-company/wrController/logger"
	"github.com/catenoid-company/wrController/models"
	"github.com/catenoid-company/wrController/utils"
	"sync"
)

type NginxApi struct {
	Data				map[string]interface{}
	RId					string
	Authorization		string
	Host				string
}
/*
	Nginx 방송 준비
	Nginx 에 방송 준비 요청
	완료후 방송 상태값 준비로 추가
 */
func (na *NginxApi) ChannelBroadService(wg *sync.WaitGroup, method string, successHosts *[]string) {
	defer wg.Done()

	webrtcApi := &utils.WebrtcApi{
		Host:   na.Host,
		Method: method,
		Data: na.Data,
		Headers: map[string]string{
			"X-Request-Id":na.RId,
			"Authorization": na.Authorization,
		},
	}

	resReadyResult := models.NginxReadyRes{}

	res, err := webrtcApi.CallApi(na.RId,"v1/channel", nil)
	if err != nil {
		logger.Error(na.RId, err.Error())
		return
	}

	b, err := utils.JsonToReader(&resReadyResult, res.Body)

	if err != nil {
		logger.Error(na.RId,"%s", err.Error())
		return
	}

	logger.Info(na.RId, "Status : %s, ResponseBody : %s", res.Status, string(b))

	if resReadyResult.ErrCode == 0 {
		*successHosts = append(*successHosts, na.Host)
		logger.Info(na.RId, "Nginx Agent request is sent")
		return
	}

	logger.WithField("resReadyResult", resReadyResult).WithField("nginxRequest", na).Error(na.RId,"NginxAgent error_code is status failed")
}

/**
	방송 준비시 공통사용 Nginx Function
	Parameters Format and Get Nginx Agent Info
 */
func GetNginxSendInfo(rId string, params interface{}) (map[string]interface{},[]string, error) {
	logger.WithField("liveApiParams:", params).Info(rId, "Get LiveApi Params")

	nginxData, err := utils.ConvertMapFromStruct(params, "json")

	if err != nil {
		logger.Error(rId, "Failed to parsed parameters : %s", err.Error())
		return nil, nil, errors.New(fmt.Sprintf("%s", config.FORMAT_FAIL_PARSE_PARAMETERS_CONTROLLER_EXCEPTION))
	}

	// Nginx Host 가 존재하는지 확인
	nginxInfo, err := utils.GetAllMembers(config.WrConfig.NginxServerInfoKey)
	if err != nil {
		logger.Error(rId, "Fail to get Nginx Host Info : %s", err.Error())
		return nil, nil, errors.New(fmt.Sprintf(config.FORMAT_GET_REDIS_HOST_CONTROLLER_EXCEPTION, config.NGINX))
	}

	nginxHosts := GetServerHosts(nginxInfo)

	// 방송준비 비지니스 로직 시작
	// LiveApi 준비 Api
	if len(nginxHosts) == 0 {
		logger.Error(rId, "Not success to responsive all Nginx hosts")
		return nil, nil, errors.New(fmt.Sprintf(config.FORMAT_UNRESPONSIVE_ALL_HOST_API_OR_NON_DATA, config.NGINX))
	}

	streamKeys, ok := nginxData["stream_keys"].([]string)
	if ok {
		if len(streamKeys) == 0{
			delete(nginxData, "stream_keys")
		}
	}else{
		delete(nginxData, "stream_keys")
	}

	return nginxData, nginxHosts, nil
}