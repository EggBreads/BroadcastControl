package utils

import (
	"bytes"
	"encoding/json"
	"github.com/catenoid-company/wrController/config"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"testing"
)

type testJsonStruct struct {
	Name	string	`json:"name"`
	Num		int		`json:"num"`
}

func Test_JsonToBytes(t *testing.T) {
	p := testJsonStruct{
		Name: "Test",
		Num: 1,
	}

	b,_ := BytesFromObject(p)

	JsonToBytes(&p,b)
	assert.Equal(
		t,
		reflect.ValueOf(p).Kind(), reflect.ValueOf(p).Kind(),
		"Fault Parsing",
		)
}

func Test_ConvertMapFromStruct(t *testing.T) {
	p := &testJsonStruct{
		Name: "Test",
		Num: 1,
	}

	m := make(map[string]interface{},0)
	v, _ := ConvertMapFromStruct(p,"json")
	assert.Equal(
		t,
		reflect.ValueOf(m).Kind(), reflect.ValueOf(v).Kind(),
		"Fault Pointer Parse",
	)

	o := testJsonStruct{
		Name: "Test",
		Num: 1,
	}

	r, _:=ConvertMapFromStruct(o,"json")

	assert.Equal(
		t,
		reflect.ValueOf(m).Kind(), reflect.ValueOf(r).Kind(),
		"Fault Struct Parse",
	)
}

func Test_JsonToReader(t *testing.T) {
	m, b := map[string]interface{}{"day": 12}, new(bytes.Buffer)

	_ = json.NewEncoder(b).Encode(m)

	//_, ok := JsonToReader(m,b).(map[string]interface{})

	assert.Equal(
		t,
		true, true,
		"Fault Parse Reader",
	)
}

//func Test_CallApi(t *testing.T) {
//	api := &WebrtcApi{
//		Host: "localhost:3000",
//		Method: "POST",
//	}
//	//
//	m := map[string]string{
//		"aaaa" : "123123",
//		"bbbb" : "202020",
//	}
//	api.CallApi("", nil, "/test/{aaaa}/test2/{bbbb}" , m)
//
//	assert.Equal(t, true, false,"Fault to Post Api")
//}

func TestRemoveSlice(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	r := make([]int, 0)

	for _, num := range nums{
		r = append(r, num)
		if num % 2 != 0{
			 i := len(r)-1
			r = append(r[:i], r[i+1:]...)
		}
	}
}

func TestCallApi(t *testing.T) {
	config.WrConfig = config.InitConfig()
	api := WebrtcApi{
		Host: "localhost:3000",
		Method: "POST",
		Data: map[string]interface{}{
			"channel_key": "x7aii7xu9mrjzdrj",
		},
		Headers: map[string]string{
			"X-Request-Id": GetRIdUUID(),
			"Authorization": AuthorizationHeader(config.WrConfig.AuthUser, config.WrConfig.AuthPass),
		},
	}

	res, err := api.CallApi(GetRIdUUID(), "v1/streaming-plugin", nil)
	if err != nil {
		log.Print(err)
		return
	}

	b, err := JsonToReader(nil, res.Body)

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(string(b))
	assert.Equal(t, 200, res.StatusCode)
}


