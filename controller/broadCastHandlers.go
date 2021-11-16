package controller

import (
	"github.com/catenoid-company/wrController/config"
	"github.com/catenoid-company/wrController/logger"
	"github.com/catenoid-company/wrController/models"
	"github.com/catenoid-company/wrController/utils"
	"github.com/catenoid-company/wrController/wrService/broadcast"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
)

type LiveApiHandlers struct{}

// @Summary Prepared BroadCast Channel
// @Description Prepare to Broadcast and Send to ready signal to nginx agent
// @Accept  application/json
// @Produce  json
// @tags Broadcast
// @Param BroadCastPrepareInfo body models.LiveParameters true "Prepare to Broadcast and Send to ready signal to nginx agent"
// @Router /channel [POST]
// @Success 200 {object} models.LivePrepareSuccessRes "The ChannelKey Prepared Broadcast"
// @Failure 400 {object} models.CommonRes "Invalid Parameters"
// @Failure 401 {object} models.CommonRes "Validation Authorization"
// @Failure 408 {object} models.CommonRes "Nginx Agent Request is timeout"
// @Failure 412 {object} models.CommonRes "All Nginx Agent Request isn't success"
// @Failure 417 {object} models.CommonRes "Webrtc Controller Server Exception"
func (lah *LiveApiHandlers) PrepareBroadCast(ctx *gin.Context) {
	rId := ctx.Request.Header.Get("X-Request-Id")
	authorization := ctx.Request.Header.Get("Authorization")

	res := models.CommonRes{}

	params := models.LiveParameters{}

	err := utils.GetParameters(&params, ctx, true, rId)
	if err != nil {
		logger.Error(rId, "Invalid Parameters : %s", params)
		res = config.FormatExceptionRes(config.INVALID_PARAMS, config.INVALID_PARAMETERS)

		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	nginxData, nginxHosts, err := broadcast.GetNginxSendInfo(rId, params)
	if err != nil {
		logger.Error(rId, "%s : ", err.Error())
		res = config.FormatExceptionRes(config.CONTROLLER_EXCEPTION, err.Error())
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(nginxHosts))

	successHosts := make([]string, 0)

	for _, host := range nginxHosts {
		nginxApi := &broadcast.NginxApi{
			Data:          nginxData,
			RId:           rId,
			Authorization: authorization,
			Host:          host,
		}

		go nginxApi.ChannelBroadService(wg, ctx.Request.Method, &successHosts)
	}
	// Goroutine Time out
	isTimeout := utils.WaitTimeout(wg)

	if isTimeout {
		// Routine Time over 되여도 실제 Nginx 에 처리에 성공한 Host 가 있을수 있기때문에 확인
		if len(successHosts) == 0 {
			logger.Error(rId, "Nginx Agent Request is timeout")

			res = config.FormatExceptionRes(config.TIMEOUT, config.FORMAT_TIMEOUT, config.NGINX)

			ctx.JSON(http.StatusRequestTimeout, res)
			return
		}
	}

	// Nginx Host 들이 Response 을 받았으나, 모두 실패응답을 받음
	if len(successHosts) == 0 {
		logger.Error(rId, "Not success to responsive all Nginx hosts")

		res = config.FormatExceptionRes(config.API_OR_NON_DATA, config.FORMAT_UNRESPONSIVE_ALL_HOST_API_OR_NON_DATA, config.NGINX)

		ctx.JSON(http.StatusPreconditionFailed, res)
		return
	}

	logger.WithField("successHosts", successHosts).Info(rId, "Finish Nginx Agent Prepared")

	ctx.JSON(http.StatusOK, models.LivePrepareSuccessRes{
		ErrorCode:  "0",
		ChannelKey: params.ChannelKey,
		StreamKeys: params.StreamKeys,
		Message:    "Success",
	})
}

// @Summary Cancel Prepared BroadCast Channel
// @Description Send to cancel signal to nginx agent
// @Accept  application/json
// @Produce  json
// @tags Broadcast
// @Param BroadCastPrepareInfo body models.LiveCancelParameters true "Send to cancel signal to nginx agent"
// @Router /channel [DELETE]
// @Success 204 {null} nil "The ChannelKey cancel on prepared broadcast"
// @Failure 400 {object} models.CommonRes "Invalid Parameters"
// @Failure 401 {object} models.CommonRes "Validation Authorization"
// @Failure 408 {object} models.CommonRes "Nginx Agent Request is timeout"
// @Failure 412 {object} models.CommonRes "All Nginx Agent Request isn't success"
// @Failure 417 {object} models.CommonRes "Webrtc Controller Server Exception"
func (lah *LiveApiHandlers) CancelPrepareBroadCast(ctx *gin.Context) {
	rId := ctx.Request.Header.Get("X-Request-Id")
	authorization := ctx.Request.Header.Get("Authorization")

	res := models.CommonRes{}

	params := models.LiveCancelParameters{}

	err := utils.GetParameters(&params, ctx, true, rId)
	if err != nil {
		logger.Error(rId, "Invalid Parameters : %s", params)
		res = config.FormatExceptionRes(config.INVALID_PARAMS, config.INVALID_PARAMETERS)

		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	nginxData, nginxHosts, err := broadcast.GetNginxSendInfo(rId, params)

	if err != nil {
		logger.Error(rId, "%s : ", err.Error())
		res = config.FormatExceptionRes(config.CONTROLLER_EXCEPTION, err.Error())
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(nginxData))

	successHosts := make([]string, 0)

	for _, host := range nginxHosts {
		nginxApi := &broadcast.NginxApi{
			Data:          nginxData,
			RId:           rId,
			Authorization: authorization,
			Host:          host,
		}

		go nginxApi.ChannelBroadService(wg, ctx.Request.Method, &successHosts)
	}

	// Goroutine Time out
	isTimeout := utils.WaitTimeout(wg)

	if isTimeout {
		// Routine Time over 되여도 실제 Nginx 에 처리에 성공한 Host 가 있을수 있기때문에 확인
		if len(successHosts) == 0 {
			logger.Error(rId, "Nginx Agent Request is timeout")

			res = config.FormatExceptionRes(config.TIMEOUT, config.FORMAT_TIMEOUT, config.NGINX)

			ctx.JSON(http.StatusRequestTimeout, res)
			return
		}
	}

	// Nginx Host 들이 Response 을 받았으나, 모두 실패응답을 받음
	if len(successHosts) == 0 {
		logger.Error(rId, "Not success to responsive all Nginx hosts")

		res = config.FormatExceptionRes(config.API_OR_NON_DATA, config.FORMAT_UNRESPONSIVE_ALL_HOST_API_OR_NON_DATA, config.NGINX)

		ctx.JSON(http.StatusPreconditionFailed, res)
		return
	}

	logger.WithField("successHosts", successHosts).Info(rId, "Finish Nginx Agent Cancel Prepared")

	ctx.Status(http.StatusNoContent)
}

// @Summary Register BroadCast
// @Description BroadCast is start. All Agent and live Stream Api send to broadcast open signal
// @Accept  json
// @Produce  json
// @tags Broadcast
// @Param BroadCastStartInfo body models.NginxParameters true "BroadCast is start. All Agent and live Stream Api send to broadcast open signal"
// @Router /publish [POST]
// @Success 200 {object} models.NginxBroadCastOpenRes{servers=[]object{ip=string},audio_profile=object{audio_codec=string,audio_sample_rate=int32,audio_bitrate=int32,audio_port=int32},video_profiles=[]object{video_codec=string,video_width=int32,video_height=int32,video_framerate=int32,video_bitrate=int64,video_port=int64}} "The ChannelKey Start Broadcast"
// @Failure 400 {object} models.CommonRes "Invalid Parameters"
// @Failure 401 {object} models.CommonRes "Validation Authorization"
// @Failure 412 {object} models.CommonRes "The response requested to another agent has failed"
// @Failure 417 {object} models.CommonRes "Webrtc Controller Server Exception"
func (lah *LiveApiHandlers) StartBroadCasting(ctx *gin.Context) {
	rId := ctx.Request.Header.Get("X-Request-Id")
	authorization := ctx.Request.Header.Get("Authorization")

	params := models.NginxParameters{}
	errorRes := models.CommonRes{}

	err := utils.GetParameters(&params, ctx, true, rId)
	if err != nil {
		logger.Error(rId, "Invalid Parameters : %s", params)
		errorRes = config.FormatExceptionRes(config.INVALID_PARAMS, config.INVALID_PARAMETERS)

		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	logger.Info(rId, "Start BroadCast Parameters from Nginx Agent : %s", params)

	// ChannelKey 가 Redis 에 있는지 확인
	storedChannelKey := utils.RedisClient.Keys("webrtc_" + params.ChannelKey + "_*").Val()
	if len(storedChannelKey) > 0 {
		logger.WithField("storedChannelKey", storedChannelKey).Error(rId, "%s", "Already used channel")

		errorRes = config.FormatExceptionRes(config.API_OR_NON_DATA, "Already used channel")

		ctx.JSON(http.StatusPreconditionFailed, errorRes)
		return
	}

	//Janus Host 가 존재하는지 확인
	janusInfo, err := utils.GetAllMembers(config.WrConfig.JanusServerInfoKey)
	if err != nil {
		logger.Error(rId, "Failed to got JanusInfo %s", err.Error())

		errorRes = config.FormatExceptionRes(config.CONTROLLER_EXCEPTION, config.FORMAT_GET_REDIS_HOST_CONTROLLER_EXCEPTION, config.JANUS)

		ctx.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	if len(janusInfo) == 0 {
		logger.Error(rId, "Not success to responsive all Janus Hosts")

		errorRes = config.FormatExceptionRes(config.API_OR_NON_DATA, config.FORMAT_UNRESPONSIVE_ALL_HOST_API_OR_NON_DATA, config.JANUS)

		ctx.JSON(http.StatusPreconditionFailed, errorRes)
		return
	}

	// 방송시작
	// liveApi service struct
	liveApi := &broadcast.LiveApi{
		RId:           rId,
		Authorization: authorization,
	}

	//liveApi 방송시작 시그널 전달 처리
	streamRes, err := liveApi.SendOpenSignalToLivaApi(params)

	if err != nil {
		logger.Error(rId, "Fail to process api to stream Api : %s", err.Error())
		errorRes = config.FormatExceptionRes(config.CONTROLLER_EXCEPTION, config.FORMAT_FAIL_CALL_API_ERR)
		ctx.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	params.BroadcastKey = streamRes["key"].(string)

	// 진행할 방송의 Unique Key 생성
	key := "webrtc_" + params.ChannelKey + "_" + params.BroadcastKey

	janusHosts := broadcast.GetServerHosts(janusInfo)

	// janusApi service struct
	janusApi := &broadcast.JanusApi{
		RId:           rId,
		Authorization: authorization,
		Hosts:         janusHosts,
	}

	/** nginx service 는 현재 Request 가 nginx 쪽이기때문에 최종적으로 완료시 Response
	  	LiveApi 에서 profile 및 port 정보 조회
	  	stream_plugin_id가 업데이트가 되었든 안되었든 무시, 어차피 새로시작된 방송의 번호로 변경되기때문
		예외가 발생된경우 해당 Key 을 조회후 제거
	*/
	liveProfiles, availablePort, janusParams, err := liveApi.GetBroadCastProfiles(params, key)

	if err != nil {
		logger.Error(rId, "Not proceed to start BroadCast from live profile : %s", err.Error())

		sErr := janusApi.StopJanusBroadCastOnStart(params.BroadcastKey, key, nil)
		if sErr != nil {
			logger.Error(rId, "%s", sErr.Error())
		}

		errorRes = config.FormatExceptionRes(config.CONTROLLER_EXCEPTION,"%s", config.FORMAT_NOT_PROCEED_START_BROADCAST_API_OR_NON_DATA)

		ctx.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	logger.WithField("janusParams:", janusParams).Info(rId, "JanusChannel pass on JanusParams to LiveService")

	// TODO janus Agent 로 StreampluginId 및 방송 정보 전달
	// TODO Janus 호출시 예외의 기준을 변경해야됨 아직 안됨
	liveParams, err := janusApi.OpenJanusBroadCastPorts(janusParams)

	if err != nil {
		//실패처리시 방송종료
		logger.Error(rId, "Not proceed to start BroadCast from : %s", err.Error())

		sErr := janusApi.StopJanusBroadCastOnStart(params.BroadcastKey, key, nil)
		if sErr != nil {
			logger.Error(rId, "%s", sErr.Error())
		}

		errorRes = config.FormatExceptionRes(config.CONTROLLER_EXCEPTION, "%s", config.FORMAT_NOT_PROCEED_START_BROADCAST_API_OR_NON_DATA)

		ctx.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	liveParams["stream_key"] = params.StreamKey

	logger.WithField("liveApiParams:", liveParams).Info(rId, "JanusService pass on janusServers to LiveChannel")

	// 방송 시작 정보를 live API 에 전달
	nginxRes, err := liveApi.SendToOpenFromLiveApi(liveProfiles, availablePort, liveParams)

	if err != nil {
		logger.Error(rId, "Not proceed to start BroadCast from live stream : %s", err.Error())

		closeJanusParams := make(map[string]interface{})

		closeJanusParams["hosts"] = liveParams["servers"]
		closeJanusParams["channel_key"] = params.ChannelKey
		closeJanusParams["stream_plugin_id"] = availablePort.StreamPluginId

		sErr := janusApi.StopJanusBroadCastOnStart(params.BroadcastKey, key, closeJanusParams)
		if sErr != nil {
			logger.Error(rId, "%s", sErr.Error())
		}

		errorRes = config.FormatExceptionRes(config.CONTROLLER_EXCEPTION, "%s", config.FORMAT_NOT_PROCEED_START_BROADCAST_API_OR_NON_DATA)

		ctx.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	// 방송시작 Nginx 응답
	logger.Info(rId, "Start to BroadCasting by NginxAgent")

	/**
		TODO 현재 방송중인 정보의 StreamPluginId을 따로 저장 하는부분 취소
	 */
	//servers := liveParams["servers"].([]map[string]string)

	//err = janusApi.SaveSuccessResJanusHosts(key, params.BroadcastKey, janusParams, servers)
	//
	//if err != nil {
	//	logger.Error(rId, config.FORMAT_SAVE_ON_AIR_BROADCAST + "%s", err.Error())
	//	errorRes = config.FormatExceptionRes(config.CONTROLLER_EXCEPTION, config.FORMAT_SAVE_ON_AIR_BROADCAST)
	//
	//	ctx.JSON(http.StatusInternalServerError, errorRes)
	//	return
	//}

	nginxRes.ErrorCode = "0"
	nginxRes.BroadcastKey = params.BroadcastKey
	nginxRes.Client = params.Client

	logger.WithField("nginxRes", nginxRes).Info(rId, "Broadcast is stating....")

	ctx.JSON(http.StatusOK, nginxRes)
}

// @Summary Close BroadCast
// @Description BroadCast is close. All Agent and live Stream Api send to broadcast close signal
// @Accept  json
// @Produce  json
// @tags Broadcast
// @Param BroadCastCloseInfo body models.NginxBroadCastCloseReq true "BroadCast is close. All Agent and live Stream Api send to broadcast close signal"
// @Router /unPublish [POST]
// @Success 200 {object} models.CommonRes "The Broadcast success to closed"
// @Success 206 {object} models.CommonRes "The Broadcast is success to closed partially"
// @Failure 400 {object} models.CommonRes "Invalid Parameters"
// @Failure 401 {object} models.CommonRes "Validation Authorization"
// @Failure 408 {object} models.CommonRes "Janus Agent Request is timeout"
// @Failure 417 {object} models.CommonRes "Webrtc BroadcastKey wasn't deleted"
func (lah *LiveApiHandlers) CloseBroadCasting(ctx *gin.Context) {
	rId := ctx.Request.Header.Get("X-Request-Id")
	authorization := ctx.Request.Header.Get("Authorization")

	logger.Info(rId, "Close to Broadcast start")

	res := models.CommonRes{
		ErrorCode: 0,
		Message:   "BroadCast is closed",
	}

	params := models.NginxBroadCastCloseReq{}

	err := utils.GetParameters(&params, ctx, true, rId)
	if err != nil {
		logger.Error(rId, "Invalid Parameters : %s", params)
		res = config.FormatExceptionRes(config.INVALID_PARAMS, config.INVALID_PARAMETERS)

		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	//Send to close Janus Agents
	janusInfo, err := utils.GetAllMembers(config.WrConfig.JanusServerInfoKey)

	if err != nil {
		logger.Error(rId, "Failed to got JanusInfo %s", err.Error())

		res = config.FormatExceptionRes(config.CONTROLLER_EXCEPTION, config.FORMAT_GET_REDIS_HOST_CONTROLLER_EXCEPTION, config.JANUS)

		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	if len(janusInfo) == 0 {
		logger.Error(rId, "Empty send to Janus hosts")

		res = config.FormatExceptionRes(config.API_OR_NON_DATA, config.FORMAT_EMPTY_REDIS_HOST_INFO_OR_NON_DATA, config.JANUS)

		ctx.JSON(http.StatusPreconditionFailed, res)
		return
	}

	key := "webrtc_" + params.ChannelKey + "_" + params.BroadcastKey

	logger.Info(rId, "Close Broadcast Key : " + key)

	// Get BroadCast Info
	broadcastInfo, err := utils.GetAllMembers(key)

	if err != nil {
		logger.Error(rId, "Fail to got broadcast information : %s", err.Error())

		res = config.FormatExceptionRes(config.CONTROLLER_EXCEPTION, config.FORMAT_GET_REDIS_WEBRTC_INFO_CONTROLLER_EXCEPTION, config.JANUS)

		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	if len(broadcastInfo) == 0 {
		logger.Error(rId, "Empty BroadCast information")

		res = config.FormatExceptionRes(config.API_OR_NON_DATA, config.FORMAT_EMPTY_REDIS_BROADCAST_API_OR_NON_DATA, config.JANUS)

		ctx.JSON(http.StatusPreconditionFailed, res)
		return
	}

	logger.WithField("ChannelKeyInfo:", broadcastInfo).Info(rId, "Get ChannelKey Info for removed")

	janusHosts := broadcast.GetServerHosts(janusInfo)

	if len(janusHosts) == 0 {
		logger.Error(rId, "Not success to responsive all janus hosts %s : %s", params.BroadcastKey)

		res = config.FormatExceptionRes(config.API_OR_NON_DATA, config.FORMAT_UNRESPONSIVE_ALL_HOST_API_OR_NON_DATA, config.JANUS)

		ctx.JSON(http.StatusPreconditionFailed, res)
		return
	}
	// TODO 일단 제외
	// TODO Janus Server Port Close
	// TODO Janus 방송종료 전송 실패와는 별개로 성공되여도
	// TODO 별도의 키도 저정된 StreamPluginId는 방송종료 StreamPluginId 이므로 제거해야됨
	//streamPluginId := broadcastInfo["stream_plugin_id"]
	//err = utils.RedisClient.HDel(config.WrConfig.StreamPluginIdOnAirKeys, streamPluginId).Err()
	//if err != nil {
	//	logger.Error(rId, "Failed to removed on air StreamPluginId %s : ", err.Error())
	//	res = config.FormatExceptionRes(config.CONTROLLER_EXCEPTION, config.FORMAT_DELETE_REDIS_KEY_MEMBER_CONTROLLER_EXCEPTION, config.JANUS)
	//	ctx.JSON(http.StatusInternalServerError, res)
	//
	//	return
	//}

	// Webrtc Broadcast Key is remove
	streamPluginId := broadcastInfo["stream_plugin_id"]
	err = utils.DeleteRedisKey(key)
	if err != nil {
		logger.Error(rId, "Don't deleted webrtc Key %s : %s", key, err.Error())

		res = config.FormatExceptionRes(config.CONTROLLER_EXCEPTION, config.FORMAT_DELETE_REDIS_KEY_CONTROLLER_EXCEPTION)

		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// LiveApi 방송종료
	liveApi := &broadcast.LiveApi{
		ChannelKey:    params.ChannelKey,
		RId:           rId,
		Authorization: authorization,
	}

	err = liveApi.CloseToLiveBroadCast(params.BroadcastKey)

	if err != nil {
		logger.Error(rId, "Not proceed to close broadcast from Kollus Stream Api %s : %s", params.BroadcastKey, err.Error())

		res = config.FormatExceptionRes(config.API_OR_NON_DATA, config.FORMAT_NOT_PROCEED_CLOSE_BROADCAST_API_OR_NON_DATA, config.KOLLUS_STREAM)

		ctx.JSON(http.StatusPreconditionFailed, res)
		return
	}

	// 마지막 StreamPluginId 관리
	err = broadcast.ManageLastStreamPluginId(broadcastInfo, rId)
	if err != nil {
		logger.Error(rId, "Failed to remove StreamPlugin and port : %s", err.Error())

		res = config.FormatExceptionRes(config.CONTROLLER_EXCEPTION, config.FORMAT_FAIL_CLOSE_BROADCAST_CONTROLLER_EXCEPTION)

		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	janusApi := &broadcast.JanusApi{
		RId:           rId,
		Authorization: authorization,
		Hosts:         janusHosts,
	}

	// Send to Janus Server Break
	wg := &sync.WaitGroup{}
	wg.Add(len(janusApi.Hosts))

	successHosts := make([]string, 0)
	UnresponsiveServers := make([]string, 0)

	for _, h := range janusApi.Hosts {
		janusParam := make(map[string]interface{}, 0)

		janusParam["host"] = h
		janusParam["stream_plugin_id"], _ = strconv.Atoi(streamPluginId)
		janusParam["channel_key"] = params.ChannelKey

		go janusApi.CloseToJanusServer(wg, janusParam, &successHosts, &UnresponsiveServers)
	}

	isTimeout := utils.WaitTimeout(wg)

	if isTimeout {
		logger.Error(rId, "Janus Agent Request is timeout")

		res = config.FormatExceptionRes(config.TIMEOUT, config.FORMAT_TIMEOUT, config.JANUS)

		ctx.JSON(http.StatusRequestTimeout, res)
		return
	}

	logger.Info(rId, "Success Janus Servers : %s", successHosts)

	if len(UnresponsiveServers) > 0 {
		logger.Warn(rId, "Unresponsive Janus Servers : %s", UnresponsiveServers)

		res = config.FormatExceptionRes(config.PARTIAL_SUCCESS, "FORMAT_UNRESPONSIVE_HOSTS_API_OR_NON_DATA")

		//ctx.JSON(http.StatusPartialContent, res)
		ctx.JSON(http.StatusOK, res)
		return
	}

	logger.Info(rId, "Remove broadcast: %s", key)

	ctx.JSON(http.StatusOK, res)
}
