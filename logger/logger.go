package logger

import (
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"
)

var logger = logrus.New()
type Fields logrus.Fields
type Entry logrus.Entry

/*
	Logger Initial
 */
func Init() {
	path := os.Getenv("WRC_LOG_FILE_PATH")
	level := os.Getenv("WRC_LOG_LEVEL")

	lv := getLevel(level)

	// Set Logger File Save physical
	logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	if len(path) > 0 {
		// wrc.log set rolling
		wrcLogger, err := SetRollingLogFile(path)

		if err != nil {
			log.Printf(path + " : %s", err.Error())
			return
		}

		logger.SetOutput(wrcLogger)
	} else {
		logger.SetOutput(os.Stdout)
	}
	logger.SetLevel(lv)
}

func SetRollingLogFile(path string) (*rotatelogs.RotateLogs, error) {
	wrcLogger, err := rotatelogs.New(
		path + ".%Y%m%d",
		rotatelogs.WithMaxAge(-1),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithLinkName(path),
	)

	if err != nil {
		return nil, err
	}

	return wrcLogger, nil
}

func getLevel(level string) (lv logrus.Level) {
	lv = logrus.InfoLevel
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
		logrus.Info("Unknown level string.")
	}
	return
}

func Info(rid string, format string, args ...interface{}) {
	Infos(3, rid, format, args...)
}
func Infos(skipCnt int, rid string, format string, args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(skipCnt)
		entry.Data["rid"] = rid
		entry.Infof(format, args...)
	}
}
func Trace(rid string, format string, args ...interface{}) {
	Traces(3, rid, format, args...)
}
func Traces(skipCnt int, rid string, format string, args ...interface{}) {
	if logger.Level >= logrus.TraceLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(skipCnt)
		entry.Data["rid"] = rid
		entry.Tracef(format, args...)
	}
}
func Debug(rid string, format string, args ...interface{}) {
	Debugs(3, rid, format, args...)
}
func Debugs(skipCnt int, rid string, format string, args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(skipCnt)
		entry.Data["rid"] = rid
		entry.Debugf(format, args...)
	}
}
func Warn(rid string, format string, args ...interface{}) {
	Warns(3, rid, format, args...)
}
func Warns(skipCnt int, rid string, format string, args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(skipCnt)
		entry.Data["rid"] = rid
		entry.Warnf(format, args...)
	}
}
func Error(rid string, format string, args ...interface{}) {
	//sArgs := make([]interface{} ,0)
	//for _, arg := range args {
	//	b ,e := json.Marshal(arg)
	//	if e != nil {
	//		sArgs = append(sArgs, arg)
	//	}else{
	//		sArgs = append(sArgs, string(b))
	//	}
	//}
	//Errors(3, rid, format, sArgs...)
	Errors(3, rid, format, args...)

}
func Errors(skipCnt int, rid string, format string, args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(skipCnt)
		entry.Data["rid"] = rid
		entry.Errorf(format, args...)
	}
}
// slice 타입은 지원하지 않으니 이걸 사용 하지 말고 WithField 를 사용 해라.
func WithStruct(obj interface{}) *Entry {
	params := make(map[string]interface{}, 0)
	if obj != nil{
		rv := reflect.ValueOf(obj)
		switch rv.Kind().String() {
		//case "slice":
		//	params["slice"] = obj
		case "struct":
			for i := 0; i < rv.Type().NumField(); i++ {
				filedName := rv.Type().Field(i).Name
				tag := filedName
				val := rv.FieldByName(filedName).Interface()
				params[tag] = val
			}
		case "map":
			for _,k := range rv.MapKeys(){
				params[k.String()] = rv.MapIndex(k).Interface()
			}
		default:
			entry := logger.WithFields(logrus.Fields{})
			entry.Data["file"] = fileInfo(2)
			entry.Warnf("Unknown Type: %s", rv.Kind().String())
		}
	}
	return (*Entry)(logger.WithFields(params))
}
func WithField(key string, obj interface{}) *Entry {
	return (*Entry)(logger.WithField(key, obj))
}
func WithFields(fields Fields) *Entry {
	return (*Entry)(logger.WithFields(logrus.Fields(fields)))
}
// for Entry
func (e *Entry) Info(rId, format string, args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		e.Data["file"] = fileInfo(2)
		e.Data["rid"] = rId
		(*logrus.Entry)(e).Infof(format, args...)
	}
}
func (e *Entry) Debug(rId, format string, args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		e.Data["file"] = fileInfo(2)
		e.Data["rid"] = rId
		(*logrus.Entry)(e).Debugf(format, args...)
	}
}
func (e *Entry) Warn(rId, format string, args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		e.Data["file"] = fileInfo(2)
		e.Data["rid"] = rId
		(*logrus.Entry)(e).Warnf(format, args...)
	}
}
func (e *Entry) Error(rId, format string, args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		e.Data["file"] = fileInfo(2)
		e.Data["rid"] = rId
		(*logrus.Entry)(e).Errorf(format, args...)
	}
}
func (e *Entry) WithField(key string, obj interface{}) *Entry {
	return (*Entry)(logger.WithField(key, obj))
}
func (e *Entry) WithFields(fields Fields) *Entry {
	return (*Entry)(logger.WithFields(logrus.Fields(fields)))
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
