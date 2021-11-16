package models

type LiveMediaBroadcasts struct {
	Id					int
	Key					string
	LiveMediaChannelId	int
	LiveCreatorId		int
}

type LiveChannelAndStream struct {
	LiveMediaChannels
	LiveCreators
}

type LiveMediaChannels struct {
	ID				int		`gorm:"primary_key"`
	Key				string	`json:"key"`
}

type LiveCreators struct {
	ID					int		`gorm:"column:creator_id"`
	LiveMediaChannelId	int		`json:"live_media_channel_id"`
	StreamKey			string	`json:"stream_key"`
}


