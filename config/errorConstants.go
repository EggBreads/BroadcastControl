package config

import (
	"fmt"
	"github.com/catenoid-company/wrController/models"
	"strings"
)
/*
	CONTROLLER_EXCEPTION = http.StatusExpectationFailed
	API_OR_NON_DATA = http.StatusExpectationFailed
	TIMEOUT = http.StatusRequestTimeout
	INVALID_PARAMS = http.StatusBadRequest
 */
const (
	PARTIAL_SUCCESS = - iota
	CONTROLLER_EXCEPTION
	API_OR_NON_DATA
	TIMEOUT
	INVALID_PARAMS
)

const (
	NGINX = "Nginx"
	JANUS = "Janus"
	KOLLUS_PROFILE = "KollusProfile"
	KOLLUS_STREAM = "KollusStream"
)

/*
	Controller 의 처리과정중 발생한 예외메세지
	http.StatusInternalServerError
*/
const (
	FORMAT_FAIL_JSON_READER_CONTROLLER_EXCEPTION = "Failed reader json"
	FORMAT_FAIL_PARSE_PARAMETERS_CONTROLLER_EXCEPTION = "Failed to parsed parameters"
	FORMAT_FAIL_PARSE_BYTES_CONTROLLER_EXCEPTION = "Failed to parsed bytes"

	FORMAT_DELETE_REDIS_KEY_CONTROLLER_EXCEPTION = "Don't deleted webrtc key : %s"
	FORMAT_DELETE_REDIS_KEY_MEMBER_CONTROLLER_EXCEPTION = "Don't deleted streampluginId : %s"

	FORMAT_GET_REDIS_HOST_CONTROLLER_EXCEPTION = "Fail to get %s Host Info"
	FORMAT_GET_REDIS_WEBRTC_INFO_CONTROLLER_EXCEPTION = "Fail to got broadcast information to %s"
	FORMAT_GET_REDIS_MONITORING_CONTROLLER_EXCEPTION = "Not found monitoring information"

	FORMAT_INSERT_REDIS_MONITORING_CONTROLLER_EXCEPTION = "RedisClient don't save monitoring information"

	FORMAT_FAIL_CLOSE_BROADCAST_CONTROLLER_EXCEPTION = "Failed to remove StreamPlugin and port"

	FORMAT_FAIL_HEALTH_CHECK_CONTROLLER_EXCEPTION = "Fail to updated healthCheck"

	FORMAT_FAIL_CALL_API_ERR = "Fail to process api"

	FORMAT_SAVE_ON_AIR_BROADCAST = "Don't saved broadcast server"

)

/*
	API 응답이 정상적으로 오지 않을경우 메세지
	http.StatusPreconditionFailed
*/
const (
	FORMAT_ALREADY_ONAIR = "Already onAir on channel"

	FORMAT_NOT_PROCEED_START_BROADCAST_API_OR_NON_DATA = "Not proceed to start BroadCast"
	FORMAT_NOT_PROCEED_CLOSE_BROADCAST_API_OR_NON_DATA = "Not proceed to close broadcast"

	FORMAT_UNRESPONSIVE_ALL_HOST_API_OR_NON_DATA = "Not success to responsive all %s hosts"

	FORMAT_EMPTY_REDIS_HOST_INFO_OR_NON_DATA = "Empty send to %s hosts"
	FORMAT_EMPTY_REDIS_BROADCAST_API_OR_NON_DATA="Empty BroadCast information to %s"
	FORMAT_EMPTY_REDIS_MONITORING_API_OR_NON_DATA="Empty %s Info"
)

/*
	API 요청 시간이 Timeout
	http.StatusRequestTimeout
*/
const (
	FORMAT_TIMEOUT = "%s Agent Request is timeout"
)

/*
	부적절한 파라미터
	http.StatusBadRequest
 */
const (
	INVALID_PARAMETERS = "Invalid Parameters"
)

/*
	Response format
 */
func FormatExceptionRes(errorType int, strFormat string, format... string) models.CommonRes {
	res := models.CommonRes{
		Message: strFormat,
	}
	if len(format) > 0 {
		res.Message = fmt.Sprintf(strFormat, format)
	}
	res.Message = strings.ReplaceAll(res.Message, "[", "")
	res.Message = strings.ReplaceAll(res.Message, "]", "")
	res.ErrorCode = errorType
	return res
}

