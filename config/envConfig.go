package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
)

var WrConfig *Configuration

type Configuration struct {
	Protocol				string		`envconfig:"WRC_PROTOCOL" default:"http"`
	ProductMode				string		`envconfig:"WRC_PRODUCT_MODE" default:"debug"`

	ProcessUID				string		`envconfig:"WRC_PROCESS_UID" default:"0"`
	ProcessGID				string		`envconfig:"WRC_PROCESS_GID" default:"0"`

	WrcBaseHost				string		`envconfig:"WRC_BASE_HOST" default:"localhost:8888"`

	LogFilePath				string		`envconfig:"WRC_LOG_FILE_PATH" default:""`
	LogLevel				string		`envconfig:"WRC_LOG_LEVEL" default:"debug"`

	MonitoringLogFilePath	string		`envconfig:"WRC_MONITORING_LOG_FILE_PATH" default:""`
	MonitoringLogLevel		string		`envconfig:"WRC_MONITORING_LOG_LEVEL" default:"debug"`

	AuthUser 				string		`envconfig:"WRC_AUTH_USER" default:"kollus"`
	AuthPass 				string		`envconfig:"WRC_AUTH_PASS" default:"0catenoid"`

	SentinelMasterName		string		`envconfig:"WRC_SENTINEL_MASTER_NAME" default:"webrtcMaster"`
	SentinelHost			string		`envconfig:"WRC_SENTINEL_HOST" default:"127.0.0.1"`
	SentinelPort			string		`envconfig:"WRC_SENTINEL_PORT" default:"26379"`
	//SentinelMasterName		string		`envconfig:"WRC_SENTINEL_MASTER_NAME" default:"kollus"`
	//SentinelHost			string		`envconfig:"WRC_SENTINEL_HOST" default:"182.252.140.143"`
	//SentinelPort			string		`envconfig:"WRC_SENTINEL_PORT" default:"16379"`


	LiveApiHost				string		`envconfig:"WRC_LIVE_API_HOST" default:"live-dev-kr.kollus.com"`

	StreamPluginIdKey		string		`envconfig:"WRC_STREAM_PLUGIN_ID_KEY" default:"webrtc_stream_port_manage"`
	StreamPluginIdOnAirKeys	string		`envconfig:"WRC_STREAM_PLUGIN_ID_ON_AIR_KEY" default:"webrtc_streamplugin_id_on_air_manage"`

	JanusServerInfoKey		string		`envconfig:"WRC_JANUS_SERVER_INFO_KEY" default:"webrtc_janus_server_info"`
	NginxServerInfoKey		string		`envconfig:"WRC_NGINX_SERVER_INFO_KEY" default:"webrtc_nginx_server_info"`


	JanusHealthCheckIsUse	string		`envconfig:"WRC_JANUS_HEALTH_CHECK_IS_USE" default:"false"`
	JanusHealthCheckTime	string		`envconfig:"WRC_JANUS_HEALTH_CHECK_TIME" default:"1"`

	ThreadTimeout			string		`envconfig:"WRC_THREAD_TIMEOUT" default:"180"`
	ClientTimeout			string		`envconfig:"WRC_CLIENT_TIMEOUT" default:"30"`

	AllowLimitTimeTerm		string		`envconfig:"WRC_ALLOW_LIMIT_TIME_TERM" default:"60"`

	BroadCastMode			string		`envconfig:"WRC_BROADCAST_MODE" default:"test"`
}

func InitConfig() *Configuration {
	config := &Configuration{}

	err := envconfig.Process("wrc", config)

	if err != nil {
		log.Println("[ERROR]", err.Error())
		os.Exit(1)
		return nil
	}

	return config
}

//func (c *Configuration) PrintConfiguration()  {
//	el := reflect.ValueOf(c).Elem()
//	for i:=0; i < el.NumField(); i++{
//		log.Println("[INFO] Configuration Information\t",
//			"[Name : ",el.Type().Field(i).Name,"]\t",
//			"[Value : ",el.Field(i),"]\t")
//	}
//}

