package models

type LiveParameters struct {
	ChannelKey			string		`json:"channel_key" validate:"required"`
	StreamKeys			[]string	`json:"stream_keys" validate:"required"`
	ContentProviderKey	string 		`json:"content_provider_key" validate:"required"`
	Record				bool		`json:"record" swaggerignore:"true" default:"false" `
}

type LiveCancelParameters struct {
	ChannelKey			string		`json:"channel_key" validate:"required"`
	StreamKeys			[]string	`json:"stream_keys"`
	ContentProviderKey	string 		`json:"content_provider_key" validate:"required"`
}

type LivePrepareSuccessRes struct {
	ErrorCode		string	`json:"error_code"`
	ChannelKey		string	`json:"channel_key"`
	StreamKeys		[]string`json:"stream_keys"`
	Message			string	`json:"message"`
}

type BroadCastProfiles struct {
	ChannelKey		string			`json:"channel_key"`
	VideoProfiles	[]videoProfiles	`json:"video_profiles"`
	AudioProfile	audioProfile	`json:"audio_profile"`
}

type videoProfiles struct {
	VideoCodec		string	`json:"video_codec"`
	VideoWidth 		int		`json:"video_width"`
	VideoHeight 	int		`json:"video_height"`
	VideoFramerate 	int		`json:"video_framerate"`
	VideoBitrate	int		`json:"video_bitrate"`
}

type audioProfile struct {
	AudioCodec 		string	`json:"audio_codec"`
	AudioSampleRate int		`json:"audio_sample_rate"`
	AudioBitrate 	int		`json:"audio_bitrate"`
}

type BroadCastStreamRes struct {
	BroadcastKey	string			`json:"broadcast_key"`
}

type BroadCastOpenRes struct {
	ErrCode	interface{}		`json:"err_code"`
}