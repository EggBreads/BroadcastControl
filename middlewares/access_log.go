package middlewares

import (
	"github.com/catenoid-company/wrController/config"
	"github.com/catenoid-company/wrController/logger"
	"github.com/catenoid-company/wrController/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func AccessLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		// Request Id Setup

		rId := c.Request.Header.Get("X-Request-Id")

		if rId == "" {
			rId = c.Request.Header.Get("X-Request-ID")
			if rId == "" {
				rId = utils.GetRIdUUID()
			}
		}
		
		c.Request.Header.Set("X-Request-Id", rId)

		c.Writer.Header().Add("X-Request-Id", rId)

		//c.Writer.Header().Add("Authorization", utils.AuthorizationHeader(config.AuthUser, config.AuthPass))

		var entry *logger.Entry

		reqUri := c.Request.RequestURI
		// Monitoring Logger Request 출력
		if strings.Contains(reqUri, config.JANUSINFO) || strings.Contains(reqUri, config.NGINXINFO) || strings.Contains(reqUri, config.HEALTH) {
			entry = (*logger.Entry)(logger.MonitoringLogger.WithFields(logrus.Fields{
				"r_status_code": c.Writer.Status(),
				"r_latency":     utils.GetDurationInMilliseconds(startTime),
				"r_client_ip":   utils.GetClientIP(c),
				"r_method":      c.Request.Method,
				"r_uri":         c.Request.RequestURI,
				"rid":           rId,
			}))
		} else {
			// BroadCast Logger Request 출력
			entry = logger.WithFields(logger.Fields{
				"r_status_code": c.Writer.Status(),
				"r_latency":     utils.GetDurationInMilliseconds(startTime),
				"r_client_ip":   utils.GetClientIP(c),
				"r_method":      c.Request.Method,
				"r_uri":         c.Request.RequestURI,
				"rid":           rId,
			})
		}

		if c.Writer.Status() >= 500 {
			entry.Warn("main", c.Errors.String())
		} else {
			entry.Info(rId, "Accessing login success")
		}

		c.Next()
	}
}
