package models


type NginxParameters struct {
	Record       bool 	`json:"record" default:"false"`
	ChannelKey   string `json:"channel_key" validate:"required"`
	Rtmp         string `json:"rtmp" validate:"required"`
	Client       string `json:"client" validate:"required"`
	Host         string `json:"host" validate:"required"`
	Server       string `json:"server" validate:"required"`
	BroadcastKey string	`json:"broadcast_key"`
	StreamKey	 string	`json:"stream_key" validate:"required"`
}

type NginxReadyRes struct {
	ErrCode int	`json:"err_code"`
}

//example="{\"audio_codec\":\"22tcp\",\"audio_sample_rate\":22,\"audio_bitrate\":11111,\"audio_port\":123123}"
type NginxBroadCastOpenRes struct {
	ChannelKey		string						`json:"channel_key"`
	ListenPort		int							`json:"listen_port"`
	BroadcastKey	string						`json:"broadcast_key"`
	Client			string						`json:"client"`
	VideoProfiles	[]map[string]interface{}	`json:"video_profiles" swaggertype:"array,object"`
	AudioProfile	map[string]interface{}		`json:"audio_profile" swaggertype:"object"`
	StreamPluginId	int							`json:"stream_plugin_id"`
	Servers			[]map[string]string			`json:"servers" swaggertype:"array,object"`
	ErrorCode	 	string						`json:"error_code"`
}

type NginxBroadCastCloseReq struct {
	ChannelKey		string	`json:"channel_key" validate:"required"`
	BroadcastKey	string	`json:"broadcast_key" validate:"required"`
}

type MonitoringNginxParams struct {
	Server		string					`json:"server"`
	Common		monitoringNginxCommon	`json:"common"`
	Nginx		monitoringNginx			`json:"nginx"`
	Rtmp2rtp	monitoringRtmp			`json:"rtmp2rtp"`
}

type CommonRes struct {
	ErrorCode		int		`json:"error_code"`
	Message			string	`json:"message"`
}

type monitoringNginxCommon struct {
	CpuNum			int			`json:"cpu_num"`
	CpuUsed			int         `json:"cpu_used"`       	// percent
	MemoryTotal		int64   	`json:"memory_total"`		// bytes
	MemoryFree		int64       `json:"memory_free"`		// bytes
}

type monitoringNginx struct {
	ConfCount 	int         `json:"conf_count"`
	HealthCheck bool		`json:"health_check"`
}

type monitoringRtmp struct {
	ModuleCount int							`json:"module_count"`
	Modules		[]monitoringRtmpModules		`json:"modules"`
}

type monitoringRtmpModules struct {
	Channel string							`json:"channel"`
	Targets	[]monitoringRtmpModulesTargets	`json:"targets"`
}

type monitoringRtmpModulesTargets struct {
	Ip 			string	`json:"ip"`
	AudioPort 	int		`json:"audio_port"`
	VideoPort1 	int		`json:"video_port1"`
	VideoPort2 	int		`json:"video_port2"`
	VideoPort3 	int		`json:"video_port3"`
}

