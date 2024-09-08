package gm

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogAlarm(t *testing.T) {
	gin.SetMode(gin.TestMode)

	//webHookURL := "http://example.com/webhook"
	// use your webhook url
	webHookURL := "https://www.feishu.cn/flow/api/trigger-webhook/7d643d2362a93141cfb69e17523c28d8"

	payload := &LogAlarmPayload{
		Trigger: AlarmTrigger{
			WebHook: []string{webHookURL},
		},
	}

	router := gin.New()
	router.Use(LogAlarm(payload))
	router.GET("/test/log", func(c *gin.Context) {
		c.Status(http.StatusInternalServerError)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test/log", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
