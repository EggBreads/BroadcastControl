package utils

import (
	"errors"
	"github.com/catenoid-company/wrController/logger"
	"github.com/catenoid-company/wrController/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MysqlGorm *gorm.DB

func ConnMysql() *gorm.DB {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "kollus:kollus@tcp(182.252.140.149:3305)/kollus_base?charset=utf8&parseTime=false"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		logger.Error("mysql Error :% ", err.Error())
	}

	return db
}

func GetMysqlBroadCastKey(channelKey string, streamKey string) (string, error)  {
	lmc := models.LiveChannelAndStream{}
	//streamKey := values["stream_key"]
	//channelKey := values["channel_key"]
	err := MysqlGorm.
		Joins(
			"left join live_creators on live_media_channels.id = live_creators.live_media_channel_id",
		).
		Select("live_media_channels.*, live_creators.id as creator_id, live_creators.live_media_channel_id, live_creators.stream_key").
		Where("live_media_channels.key = ? and live_creators.stream_key = ?",
			channelKey,
			streamKey,
		).Table("live_media_channels").Find(&lmc).Debug().Error

	if err != nil {
		return "", err
	}

	if lmc.LiveMediaChannels.ID == 0 || lmc.LiveCreators.ID == 0{
		return "",errors.New("Empty Channel\n")
	}

	lmb := models.LiveMediaBroadcasts{}
	err = MysqlGorm.
		Where("live_media_broadcasts.live_media_channel_id = ? and live_media_broadcasts.live_creator_id = ?",
			lmc.LiveMediaChannels.ID, lmc.LiveCreators.ID).
		Last(&lmb).
		Debug().Error

	if err != nil {
		return "", err
	}

	if lmb.Key == ""{
		return "",errors.New("Empty BroadCast\n")
	}

	return lmb.Key, nil
}
