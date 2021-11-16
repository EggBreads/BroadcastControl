package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/catenoid-company/wrController/config"
	"github.com/catenoid-company/wrController/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"
)

/*
	Byte 을 Json 으로 변경
*/
func JsonToBytes(obj interface{}, byt []byte) error {
	err := json.Unmarshal(byt, obj)

	if err != nil {
		return err
	}

	return nil
}

/*
	Bytes 로 변경
*/
func BytesFromObject(obj interface{}) ([]byte, error) {
	b, e := json.Marshal(obj)

	if e != nil {
		return nil, e
	}

	return b, nil
}

/*
	Struct 을 Map 형태로 변경
	해당 변경을 원하는 Struct 는 반드시 Tag 가 존재
*/
func ConvertMapFromStruct(obj interface{}, tagName string) (map[string]interface{}, error) {
	values := make(map[string]interface{}, 0)

	rv := reflect.ValueOf(obj)

	kind := rv.Type().Kind().String()

	switch kind {
	case "ptr":
		el := reflect.ValueOf(obj).Elem()

		for i := 0; i < el.NumField(); i++ {
			elFiled := el.Field(i)
			tag := el.Type().Field(i).Name
			if tagName != "" {
				tag = el.Type().Field(i).Tag.Get(tagName)
			}
			values[tag] = elFiled.Interface()
		}

		break
	case "struct":
		for i := 0; i < rv.Type().NumField(); i++ {
			filedName := rv.Type().Field(i).Name

			tag := filedName
			if tagName != "" {
				tag = rv.Type().Field(i).Tag.Get(tagName)
			}
			val := rv.FieldByName(filedName).Interface()

			values[tag] = val
		}
		break
	case "slice":
		for i := 0; i < rv.Len(); i++ {
			//sRv := reflect.ValueOf(rv.Index(i))
			values[strconv.Itoa(i)] = rv.Index(i).Interface()
		}

		break
	case "map":
		for _, k := range rv.MapKeys() {
			values[k.String()] = rv.MapIndex(k).Interface()
		}
		break
	case "string":
		err := JsonToBytes(&values, []byte(obj.(string)))
		if err != nil {
			return nil, err
		}

		break
	default:
		break
	}

	return values, nil
}

/*
	Response 의 정보를 Parsing 하기 위해 사용
*/
//func responseTypeParse(obj interface{}, res *http.Response) error {
//	byt , err := ioutil.ReadAll(res.Body)
//
//	if err != nil {
//		return err
//	}
//
//	err = JsonToBytes(obj, byt)
//
//	if err != nil {
//		return err
//	}
//
//	return err
//}
/*
	Method 이름으로 호출하여 사용하기 위해 사용
*/
func RunFromCallMethodName(obj interface{}, methodName string, args ...string) error {
	rt := reflect.TypeOf(obj)
	fn, isMethod := rt.MethodByName(methodName)

	if !isMethod {
		return errors.New("[ERROR] Not match method name")
	}

	arg := make([]reflect.Value, 0)
	arg = append(arg, reflect.ValueOf(obj))

	//if len(args) > 0 && args != nil{
	if len(args) > 0 {
		for _, v := range args {
			arg = append(arg, reflect.ValueOf(v))
		}
	}

	val := fn.Func.Call(arg)

	if len(val) == 0 {
		return nil
	}

	if val[0].Interface() != nil {
		return val[0].Interface().(error)
	}

	return nil
}

/*
	Parameter 정보를 한번에 처리하기위해 사용
	obj => pointer 로 사용해야 처리가됨
*/
func GetParameters(obj interface{}, c *gin.Context, isValid bool, rId string) error {
	//params := make(map[string]interface{}, 0)

	b, err := JsonToReader(obj, c.Request.Body)
	if err != nil {
		return err
	}

	reqUri := c.Request.RequestURI
	// Monitoring Logger Request Parameter 출력
	if strings.Contains(reqUri, config.JANUSINFO) || strings.Contains(reqUri, config.NGINXINFO) || strings.Contains(reqUri, config.HEALTH) {
		logger.MonitoringLogger.WithField("rId", rId).Infof("Requests Params : %s", string(b))
	} else {
		// BroadCast Logger Request Parameter 출력
		logger.Info(rId, "Requests Params : %s", string(b))
	}

	if !isValid {
		return nil
	}

	validate := validator.New()
	err = validate.Struct(obj)
	if err != nil {
		return err
	}

	return nil
}

//func GetParameters(obj interface{}, c *gin.Context, isValid bool, rId string)  error {
//	el := reflect.ValueOf(obj).Elem()
//	params := make(map[string]interface{},0)
//
//	if strings.LastIndex(c.Request.Header.Get("Content-Type"), "json") > -1 {
//		_, err := JsonToReader(&params, c.Request.Body)
//		if err != nil {
//			return err
//		}
//	}else{
//		for i:=0; i < el.NumField(); i++{
//			if c.Request.Method == http.MethodGet{
//				tag := el.Type().Field(i).Tag.Get("json")
//				val := c.Param(tag)
//				if val == ""{
//					val = c.Query(tag)
//				}
//				if val != "" {
//					params[tag] = val
//				}
//			}else{
//				tag := el.Type().Field(i).Tag.Get("json")
//				val := c.Param(tag)
//
//				if val == ""{
//					val = c.PostForm(tag)
//				}
//				if val == ""{
//					val = c.Query(tag)
//				}
//				if val != "" {
//					params[tag] = val
//				}
//			}
//		}
//	}
//
//	byt, err := json.Marshal(params)
//	if err != nil {
//		return err
//	}
//
//	reqUri := c.Request.RequestURI
//	// Monitoring Logger Request Parameter 출력
//	if strings.Contains(reqUri, config.JANUSINFO) || strings.Contains(reqUri, config.NGINXINFO) || strings.Contains(reqUri, config.HEALTH) {
//		logger.MonitoringLogger.WithField("rId", rId).Infof( "Requests Params : %s", string(byt))
//	}else{
//		// BroadCast Logger Request Parameter 출력
//		logger.Info(rId, "Requests Params : %s", string(byt))
//	}
//
//	err = JsonToBytes(obj, byt)
//	if err != nil {
//		return err
//	}
//
//	if !isValid{
//		return nil
//	}
//
//	validate := validator.New()
//	err = validate.Struct(obj)
//	if err != nil{
//		return err
//	}
//
//	return nil
//}

func JsonToReader(obj interface{}, reader io.Reader) ([]byte, error) {
	if obj == nil {
		b, err := ioutil.ReadAll(reader)

		if err != nil {
			return nil, err
		}
		return b, nil
	}

	d := json.NewDecoder(reader)

	err := d.Decode(obj)
	if err != nil {
		return nil, err
	}

	b, err := BytesFromObject(obj)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GetKstTime(layout string) time.Time {
	loc, _ := time.LoadLocation("Asia/Seoul")
	t, _ := time.ParseInLocation(layout, time.Now().Format(layout), loc)
	return t
}

// GetClientIP gets the correct IP for the end client instead of the proxy
func GetClientIP(c *gin.Context) string {
	// first check the X-Forwarded-For header
	requester := c.Request.Header.Get("X-Forwarded-For")
	// if empty, check the Real-IP header
	if len(requester) == 0 {
		requester = c.Request.Header.Get("X-Real-IP")
	}
	// if the requester is still empty, use the hard-coded address from the socket
	if len(requester) == 0 {
		requester = c.Request.RemoteAddr
	}

	// if requester is a comma delimited list, take the first one
	// (this happens when proxied via elastic load balancer then again through nginx)
	if strings.Contains(requester, ",") {
		requester = strings.Split(requester, ",")[0]
	}

	return requester
}

func GetDurationInMilliseconds(start time.Time) float64 {
	end := time.Now()
	duration := end.Sub(start)
	milliseconds := float64(duration) / float64(time.Millisecond)
	rounded := float64(int(milliseconds*100+.5)) / 100
	return rounded
}

/*
	UUID 생성
	Bit 변경은 일단 주석
*/
func GetRIdUUID() string {
	u := uuid.NewV4()
	// u.SetVariant(uuid.VariantFuture)
	return u.String()
}

type WebrtcApi struct {
	Host     string
	Method   string `default:"GET"`
	Data     map[string]interface{}
	RawQuery map[string]string
	Headers  map[string]string
	//Err		 error
}

/*
	Api 호출
	*obj => Response 의 응담에 해당하는 model 만약 미사용시 nil 로 표기하며 표기시 Response 가 Return
	WebrtcApi Struct 을 만들어서 함수 CallApi 사용
*/

func (wa *WebrtcApi) CallApi(rId string, path string, pathParams map[string]string) (*http.Response, error) {
	if wa.Host == "" {
		return nil, errors.New("Empty api host info\n")
	}

	formatPath := wa.getPathParams(path, pathParams)

	uri := wa.getUri(formatPath)

	reader, err := wa.getBodyData()
	if err != nil {
		return nil, err
	}

	// New Create Request
	req, err := http.NewRequest(wa.Method, uri, reader)
	if err != nil {
		return nil, err
	}

	// Add to Headers
	if wa.Headers != nil {
		if len(wa.Headers) > 0 {
			for k, v := range wa.Headers {
				req.Header.Add(k, v)
			}
		}
	}

	if wa.Method != "GET" {
		if wa.Headers["Content-Type"] == "" {
			req.Header.Set("Content-Type", "application/json")
		}
	}

	logger.WithField("Host", req.Host).WithField("Method", wa.Data).WithField("Host", req.Host).Info(rId, "Before request submit Host")

	clientTimeout, err := strconv.Atoi(config.WrConfig.ClientTimeout)

	// 형변환 확인
	if err != nil {
		return nil, err
	}

	// Call Request Api Default Values
	client := &http.Client{
		Timeout: time.Duration(clientTimeout) * time.Second,
	}

	res, err := client.Do(req)

	if err != nil {
		return res, err
	}

	logger.Info(rId, "CallApi Result Type is http.response")

	return res, nil
}

func (wa *WebrtcApi) getUri(path string) string {
	queryString := wa.getQueryString()
	u := &url.URL{
		Scheme:   config.WrConfig.Protocol,
		Host:     wa.Host,
		Path:     path,
		RawQuery: queryString,
	}

	return u.String()
}

func (wa *WebrtcApi) getBodyData() (io.Reader, error) {
	if wa.Data == nil {
		return nil, nil
	}
	//ContentType 에따라 DataType 변경
	if wa.Headers["Content-Type"] == "application/x-www-form-urlencoded" {

		data := url.Values{}

		for k, v := range wa.Data {
			data.Set(k, v.(string))
		}

		return strings.NewReader(data.Encode()), nil
	}
	// application/json
	byt, err := BytesFromObject(wa.Data)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(byt), nil
}

func (wa *WebrtcApi) getPathParams(path string, pathParams map[string]string) string {
	if pathParams != nil {
		for k, v := range pathParams {
			parsePath := fmt.Sprintf("{%s}", k)
			path = strings.Replace(path, parsePath, v, 1)
		}
	}
	return path
}

func (wa *WebrtcApi) getQueryString() string {
	if wa.RawQuery == nil {
		return ""
	}
	strSlice := make([]string, 0)
	for k, v := range wa.Data {
		strSlice = append(strSlice, k+"="+v.(string))
	}
	return strings.Join(strSlice, "&")
}

/*
	Base 64 Authorization Parse
*/

func AuthorizationHeader(user, password string) string {
	base := user + ":" + password

	return "Basic " + base64.StdEncoding.EncodeToString(stringToBytes(base))
}

func stringToBytes(s string) (b []byte) {
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return b
}

// Agent 에 Request 요청시 goroutine 의 timeout 처리
// clientTimeout(RequestTimeout) * 전체 hosts 갯수 * goroutine timeout
func WaitTimeout(wg *sync.WaitGroup) bool {
	complete := make(chan bool, 0)
	goroutineTimeOut, _ := strconv.Atoi(config.WrConfig.ThreadTimeout)

	timeout := goroutineTimeOut

	go func() {
		defer close(complete)
		wg.Wait()
	}()

	select {
	case <-complete:
		return false

	case <-time.After(time.Duration(timeout) * time.Second):
		return true
	}
}

//func RemoveIdxSlice(i interface{}, idx int) interface{} {
//	s, ok := i.([]interface{})
//	if !ok {
//		return i
//	}
//	return append(s[:idx], s[idx+1:]...)
//}
