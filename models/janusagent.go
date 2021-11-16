package models

type JanusStreamRes struct {
	ErrCode	int						`json:"err_code"`
	ErrMsg	string					`json:"err_msg"`
	Data	JanusOpenStreamServer	`json:"data"`
}

type JanusOpenStreamServer struct {
	ServerIp	string	`json:"server_ip"`
}

type JanusMonitoringParams struct {
	Server		string					`json:"server"`
	Common		monitoringJanusCommon	`json:"common"`
	Janus		monitoringJanus			`json:"janus"`
	Turn		monitoringTurn			`json:"turn"`
}

type monitoringJanusCommon struct {
	CpuNum			int			`json:"cpu_num"`
	CpuUsed			int         `json:"cpu_used"`       	// percent
	MemoryTotal		int64   	`json:"memory_total"`		// bytes
	MemoryFree		int64       `json:"memory_free"`		// bytes
}

type monitoringJanus struct {
	StreamPluginCount 		int							`json:"stream_plugin_count"`
	HealthCheck 			bool						`json:"health_check"`
	SessionCount			int							`json:"session_count"`
	Streams					[]monitoringJanusStream		`json:"streams"`
	Streamings				[]streamings				`json:"streamings"`
}

type monitoringJanusStream struct {
	Id 				int		`json:"id"`
	AudioPort 		int		`json:"audio_port"`
	VideoPort1 		int		`json:"video_port1"`
	VideoPort2 		int		`json:"video_port2"`
	VideoPort3		int		`json:"video_port3"`
	SessionCount 	int		`json:"session_count"`
}

type monitoringTurn struct {
	HealthCheck bool	`json:"health_check"`
}

type streamings struct {
	Id				int			`json:"id"`
	Description		string		`json:"description"`
	AudioAgeMs		int64		`json:"audio_age_ms"`
	VideoAgeMs		int64		`json:"video_age_ms"`
}
