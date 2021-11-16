package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	g *gin.Engine
	res *httptest.ResponseRecorder
)

func Test_Init(t *testing.T) {
	//utils.RedisClient = SetConfigAndRedis()
	//g = SetWrControllerServer()

}

func Test_PrepareBroadCast(t *testing.T) {
	res = httptest.NewRecorder()
	d1,buf1 := map[string]interface{}{
		"channel_key":"qwertyuiop2",
		"stream_key":"abcdef2",
	},new(bytes.Buffer)

	e := json.NewEncoder(buf1).Encode(d1)

	assert.NoError(t, e, "Fault to parse json reader")

	r1 := httptest.NewRequest("POST", "/channel", buf1)
	r1.Header.Add("Content-Type","application/json")
	g.ServeHTTP(res,r1)

	assert.Equal(t, http.StatusOK, res.Code, "Fault to Prepared BroadCast")
}

func Test_OpenBroadCast(t *testing.T) {
	res = httptest.NewRecorder()
	d2,buf2 := map[string]interface{}{
		"r_id": "qwertysdfgh2",
		"channel_key": "qwertyuiop2",
		"rtmp": "rtmp://localhost:1935/qwertyuiop2",
		"client": "11.22.11.22",
		"host": "kr01wr01",
		"server": "1.2.3.4",
	},new(bytes.Buffer)

	e := json.NewEncoder(buf2).Encode(d2)

	assert.NoError(t, e, "Fault to parse Open BroadCast")

	r2 := httptest.NewRequest("POST", "/publish", buf2)
	r2.Header.Add("Content-Type","application/json")
	g.ServeHTTP(res,r2)

	assert.Equal(t, http.StatusOK, res.Code, "Fault to Open Broadcast")
}

func Test_BreakBroadCast(t *testing.T) {
	res = httptest.NewRecorder()
	d3,buf3 := map[string]interface{}{
		"channel_key": "qwertyuiop2",
		"broadcast_key": "abcdef2",
	},new(bytes.Buffer)

	e := json.NewEncoder(buf3).Encode(d3)

	assert.NoError(t, e, "Fault to parse Open BroadCast")

	r3 := httptest.NewRequest("POST", "/unPublish", buf3)
	r3.Header.Add("Content-Type","application/json")
	g.ServeHTTP(res,r3)

	assert.Equal(t, http.StatusOK, res.Code, "Fault to parse Break Broadcast")
}


