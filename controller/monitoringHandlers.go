package controller

import (
	"github.com/catenoid-company/wrController/config"
	"github.com/catenoid-company/wrController/logger"
	"github.com/catenoid-company/wrController/models"
	"github.com/catenoid-company/wrController/utils"
	"github.com/catenoid-company/wrController/wrService/monitoring"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type MonitoringHandler struct{}

// @Summary Register Janus Monitoring information
// @Description Janus Agent and Janus Turn Status Monitoring information register
// @Accept  json
// @Produce  json
// @tags Monitoring
// @Param BroadCastCloseInfo body models.CommonRes true "BroadCast is close. All Agent and live Stream Api send to broadcast close signal"
// @Router /janusinfo [POST]
// @Success 200 {object} models.CommonRes "Janus Monitoring Info registered"
// @Failure 400 {object} models.CommonRes "Invalid Parameters"
// @Failure 401 {object} models.CommonRes "Validation Authorization"
// @Failure 417 {object} models.CommonRes "Janus Monitoring Info registered Exception"
func (mh MonitoringHandler) SaveJanusMonitoringHandler(ctx *gin.Context) {
	rId := ctx.Request.Header.Get("X-Request-Id")
	authorization := ctx.Request.Header.Get("Authorization")

	res := models.CommonRes{
		ErrorCode: 0,
		Message:   "Janus information is saved",
	}

	// Get Parameters
	//params := models.JanusMonitoringParams{}
	params := make(map[string]interface{}, 0)

	err := utils.GetParameters(&params, ctx, false, rId)
	if err != nil {
		logger.MonitoringLogger.WithField("rid", rId).Errorf("Failed to parsed parameters : %s", err.Error())

		res = config.FormatExceptionRes(config.INVALID_PARAMS, config.INVALID_PARAMETERS)

		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	logger.MonitoringLogger.WithField("rid", rId).WithField("params", params).Info("Monitoring Params")

	// 최종 등록된 시간 저장
	term, _ := strconv.Atoi(config.WrConfig.AllowLimitTimeTerm)
	params["allow_limit_time"] = time.Now().Add(time.Duration(term) * time.Minute).Unix()

	logger.MonitoringLogger.WithField("rid", rId).Infof(rId, "Monitoring Parsed Params : %s", params)

	service := &monitoring.JanusMonitoringService{
		RId: rId,
		Authorization: authorization,
	}

	janus, ok := params["janus"].(map[string]interface{})
	if !ok {
		logger.MonitoringLogger.WithField("rid", rId).Errorf("Empty janus field : %s", janus)
	}

	if len(janus) > 0 {
		arrStreamings, ok := janus["streamings"].([]interface{})
		if !ok {
			logger.MonitoringLogger.WithField("rid", rId).Errorf("Empty streamings field : %s", janus)
		}

		arrM := make([]map[string]interface{}, 0)

		if len(arrStreamings) > 0 {
			for _, v := range arrStreamings {
				m, ok := v.(map[string]interface{})
				if !ok {
					logger.MonitoringLogger.WithField("rid", rId).Warnf("Non parsed streamings: %s", janus)
					continue
				}

				if len(m) > 0 {
					arrM = append(arrM, m)
				}
			}

			if len(arrM) > 0 {
				// Remain to closed janus server
				janus["streamings"] = service.CloseToRemainJanusBroadCast(arrM, params["server"].(string))
				params["janus"] = janus
			}
		}
	}

	// Save to Janus Information by RedisClient
	err = service.SaveJanusMonitoring(params)

	if err != nil {
		logger.MonitoringLogger.WithField("rid", rId).Errorf("RedisClient don't save monitoring information : %s", err.Error())

		res = config.FormatExceptionRes(config.CONTROLLER_EXCEPTION, config.FORMAT_INSERT_REDIS_MONITORING_CONTROLLER_EXCEPTION)

		//res.ErrorCode = "-1"
		//res.Message = "RedisClient don't save monitoring information"

		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	logger.MonitoringLogger.WithField("rid", rId).Info("Janus Information save to success")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Register Janus Monitoring information
// @Description Janus Agent and Janus Trun Status Monitoring information register
// @Accept  json
// @Produce  json
// @tags Monitoring
// @Param BroadCastCloseInfo body models.MonitoringNginxParams true "BroadCast is close. All Agent and live Stream Api send to broadcast close signal"
// @Router /nginxinfo [POST]
// @Success 200 {object} models.CommonRes "Janus Monitoring Info registered"
// @Failure 400 {object} models.CommonRes "Invalid Parameters"
// @Failure 401 {object} models.CommonRes "Validation Authorization"
// @Failure 417 {object} models.CommonRes "Janus Monitoring Info registered Exception"
func (mh MonitoringHandler) SaveNginxMonitoringHandler(ctx *gin.Context) {
	rId := ctx.Request.Header.Get("X-Request-Id")
	authorization := ctx.Request.Header.Get("Authorization")

	res := models.CommonRes{
		ErrorCode: 0,
		Message:   "Nginx information is saved",
	}

	// Get Parameters
	//params := models.MonitoringNginxParams{}
	params := make(map[string]interface{}, 0)
	err := utils.GetParameters(&params, ctx, false, rId)

	if err != nil {
		logger.MonitoringLogger.WithField("rid", rId).Errorf("Failed to parsed parameters : %s", err.Error())

		res = config.FormatExceptionRes(config.INVALID_PARAMS, config.INVALID_PARAMETERS)

		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	p, err := utils.BytesFromObject(params)
	if err != nil {
		logger.MonitoringLogger.WithField("rid", rId).Errorf("Failed reader json : %s", err.Error())

		res = config.FormatExceptionRes(config.CONTROLLER_EXCEPTION, config.FORMAT_FAIL_JSON_READER_CONTROLLER_EXCEPTION)

		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	logger.MonitoringLogger.WithField("rid", rId).Infof("Get Nginx Info Params : %s", string(p))

	// 최종 등록된 시간 저장
	term, _ := strconv.Atoi(config.WrConfig.AllowLimitTimeTerm)
	params["allow_limit_time"] = time.Now().Add(time.Duration(term) * time.Minute).Unix()

	// Set NginxMonitoring Struct Field by Type
	service := &monitoring.NginxMonitoringService{
		RId: rId,
		Authorization: authorization,
	}

	// Save to Nginx Information by RedisClient
	err = service.SaveNginxMonitoring(params)

	if err != nil {
		logger.MonitoringLogger.WithField("rid", rId).Errorf("RedisClient don't save monitoring information : %s", err.Error())

		res = config.FormatExceptionRes(config.CONTROLLER_EXCEPTION, config.FORMAT_INSERT_REDIS_MONITORING_CONTROLLER_EXCEPTION)

		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// @Summary Get Janus Monitoring information
// @Description Get to Janus Agent and Janus Turn Status Monitoring information
// @Accept  json
// @Produce  json
// @tags Monitoring
// @Router /janusinfo [GET]
// @Success 200 {object} []object{}
// @Failure 400 {object} models.CommonRes "Invalid Parameters"
// @Failure 401 {object} models.CommonRes "Validation Authorization"
// @Failure 417 {object} models.CommonRes "Janus Monitoring Info got Exception"
func (mh MonitoringHandler) GetJanusServersInfo(ctx *gin.Context) {
	rId := ctx.Request.Header.Get("X-Request-Id")
	authorization := ctx.Request.Header.Get("Authorization")

	errRes := models.CommonRes{
		ErrorCode: 0,
		Message:   "",
	}

	service := &monitoring.JanusMonitoringService{
		RId: rId,
		Authorization: authorization,
	}

	monitoringResult, err := service.GetMonitoring()

	if err != nil {
		logger.MonitoringLogger.WithField("rid", rId).Errorf("Not found monitoring information : %s", err.Error())

		errRes = config.FormatExceptionRes(config.CONTROLLER_EXCEPTION, config.FORMAT_GET_REDIS_MONITORING_CONTROLLER_EXCEPTION)

		ctx.JSON(http.StatusInternalServerError, errRes)
		return
	}

	//res := make([]models.JanusMonitoringParams, 0)
	res := make([]map[string]interface{}, 0)

	for k, v := range monitoringResult {
		m := make(map[string]interface{})

		err := utils.JsonToBytes(&m, []byte(v.(string)))
		if err != nil {
			logger.MonitoringLogger.WithField("rid", rId).Errorf("Failed bytes json : %s", err.Error())
		} else {
			m["server"] = k
			res = append(res, m)
		}
	}

	if len(res) == 0 {
		logger.MonitoringLogger.WithField("rid", rId).Error("Empty JanusMonitoring Info")

		errRes = config.FormatExceptionRes(config.API_OR_NON_DATA, config.FORMAT_EMPTY_REDIS_MONITORING_API_OR_NON_DATA)

		ctx.JSON(http.StatusPreconditionFailed, errRes)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// @Summary Get Nginx Agent Monitoring information
// @Description Get to Nginx Agent and Nginx Turn Status Monitoring information
// @Accept  json
// @Produce  json
// @tags Monitoring
// @Router /nginxinfo [GET]
// @Success 200 {object} []object{}
// @Failure 400 {object} models.CommonRes "Invalid Parameters"
// @Failure 401 {object} models.CommonRes "Validation Authorization"
// @Failure 417 {object} models.CommonRes "Nginx Monitoring Info got Exception"
func (mh MonitoringHandler) GetNginxServersInfo(ctx *gin.Context) {
	rId := ctx.Request.Header.Get("X-Request-Id")
	authorization := ctx.Request.Header.Get("Authorization")

	errRes := models.CommonRes{
		ErrorCode: 0,
		Message:   "",
	}

	service := &monitoring.NginxMonitoringService{
		RId: rId,
		Authorization: authorization,
	}

	monitoringResult, err := service.GetMonitoring()

	if err != nil {
		logger.MonitoringLogger.WithField("rid", rId).Errorf("Not found monitoring information : %s", err.Error())

		errRes = config.FormatExceptionRes(config.CONTROLLER_EXCEPTION, config.FORMAT_GET_REDIS_MONITORING_CONTROLLER_EXCEPTION)

		ctx.JSON(http.StatusInternalServerError, errRes)
		return
	}

	res := make([]map[string]interface{}, 0)
	for k, v := range monitoringResult {
		m := make(map[string]interface{}, 0)

		err := utils.JsonToBytes(&m, []byte(v.(string)))
		if err != nil {
			logger.MonitoringLogger.Errorf("q : %s", err.Error())
			continue
		}

		m["server"] = k
		res = append(res, m)
	}

	//res := &models.MonitoringNginxParams{}
	//
	//for k,v := range monitoringResult {
	//	strJson := "{\""+ k + "\":" + v.(string) + "}"
	//	ok := false
	//	res, ok = utils.JsonToBytes(res,[]byte(strJson)).(*models.MonitoringNginxParams)
	//
	//	if !ok {
	//		ctx.Status(http.StatusExpectationFailed)
	//		return
	//	}
	//}

	if len(res) == 0 {
		logger.MonitoringLogger.WithField("rid", rId).Error("Empty NginxMonitoring Info")

		errRes = config.FormatExceptionRes(config.API_OR_NON_DATA, config.FORMAT_EMPTY_REDIS_MONITORING_API_OR_NON_DATA, config.NGINX)

		ctx.JSON(http.StatusPreconditionFailed, errRes)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// @Summary Nginx Agent And Janus Agent Health Check
// @Description Nginx Agent And Janus Agent Health Check
// @Accept  application/json
// @Produce  json
// @tags Monitoring
// @param server body object{server=string} true "Server Ip"
// @Router /health [POST]
// @Success 200 {object} models.CommonRes
// @Failure 400 {object} models.CommonRes "Invalid Parameters"
// @Failure 401 {object} models.CommonRes "Validation Authorization"
// @Failure 417 {object} models.CommonRes "Health Check Process Exception"
func (mh MonitoringHandler) AgentHealthCheck(ctx *gin.Context) {
	rId := ctx.Request.Header.Get("X-Request-Id")

	res := models.CommonRes{
		ErrorCode: 0,
		Message:   "Success Health Check",
	}

	params := map[string]interface{}{}

	b, err := utils.JsonToReader(&params, ctx.Request.Body)
	if err != nil {
		logger.MonitoringLogger.WithField("rid", rId).Errorf("Failed reader json : %s", err.Error())

		res = config.FormatExceptionRes(config.CONTROLLER_EXCEPTION, config.FORMAT_FAIL_JSON_READER_CONTROLLER_EXCEPTION)

		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	logger.MonitoringLogger.Info(rId, "HealthCheck Parameters : %s", string(b))

	server := params["server"].(string)

	if server == "" {
		logger.MonitoringLogger.WithField("rid", rId).Errorf("Invalid Parameter : %s", params)

		res = config.FormatExceptionRes(config.INVALID_PARAMS, config.INVALID_PARAMETERS)

		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	nginxService := &monitoring.NginxMonitoringService{
		RId: rId,
	}
	janusService := &monitoring.JanusMonitoringService{
		RId: rId,
	}

	logger.MonitoringLogger.WithField("rid", rId).Infof("Find to service : %s", server)

	nErr := nginxService.NginxRefreshAllowTime(server)
	jErr := janusService.JanusRefreshAllowTime(server)

	if nErr != nil && jErr != nil {
		if nErr != nil {
			logger.MonitoringLogger.WithField("rid", rId).Errorf("Fail to updated healthCheck : %s", nErr.Error())
		}

		if jErr != nil {
			logger.MonitoringLogger.WithField("rid", rId).Errorf("Fail to updated healthCheck : %s", jErr.Error())
		}

		res = config.FormatExceptionRes(config.CONTROLLER_EXCEPTION, config.FORMAT_FAIL_HEALTH_CHECK_CONTROLLER_EXCEPTION)

		ctx.JSON(http.StatusInternalServerError, res)
		return
	}


	ctx.JSON(http.StatusOK, res)
}
