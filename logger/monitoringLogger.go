package logger

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strings"
	"time"
)

var MonitoringLogger = logrus.New()

func MonitoringLoggerInit() {
	path := os.Getenv("WRC_MONITORING_LOG_FILE_PATH")
	level := os.Getenv("WRC_MONITORING_LOG_LEVEL")

	lv := logrus.InfoLevel
	switch strings.ToLower(level) {
	case "debug":
		lv = logrus.DebugLevel
	case "info":
		lv = logrus.InfoLevel
	case "warn":
		lv = logrus.WarnLevel
	case "error":
		lv = logrus.ErrorLevel
	default:
		logrus.Info("%s","Unknown level string.")
	}

	// Set Monitoring Logger File Save physical
	MonitoringLogger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	if len(path) > 0 {
		// wrcMonitoring.log set rolling
		wrcMonitoringLogger, err := SetRollingLogFile(path)

		if err != nil {
			log.Printf(path + " : %s", err.Error())
			return
		}

		MonitoringLogger.SetOutput(wrcMonitoringLogger)
	} else {
		MonitoringLogger.SetOutput(os.Stdout)
	}
	MonitoringLogger.SetLevel(lv)
}


