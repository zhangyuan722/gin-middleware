package gm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type LogAlarmPayload struct {
	Trigger AlarmTrigger
	Params  []string
}

type AlarmTrigger struct {
	WebHook []string
}

func LogAlarm(p *LogAlarmPayload) gin.HandlerFunc {
	return func(c *gin.Context) {
		beginAt := time.Now()
		w := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		status := c.Writer.Status()
		if status < http.StatusInternalServerError {
			return
		}

		nowAt := time.Now()

		paramMap := map[string]string{}

		paramMap["timestamp"] = nowAt.String()

		if u, exists := c.Get(CtxClaims); exists {
			paramMap["user"] = u.(*Claims).ID
		}

		paramMap["status"] = strconv.FormatInt(int64(status), 10)

		jsonData := make(map[string]any)
		_ = json.Unmarshal(w.body.Bytes(), &jsonData)
		respBodyString, _ := json.Marshal(jsonData)

		paramMap["response"] = string(respBodyString)

		paramMap["timing"] = strconv.FormatInt(nowAt.Sub(beginAt).Milliseconds(), 10)

		paramMap["curl"] = getCURL(c)

		if len(p.Trigger.WebHook) > 0 {
			triggerWebHook(p.Trigger.WebHook, paramMap)
		}
	}
}

func triggerWebHook(urls []string, paramMap map[string]string) {
	params := url.Values{}
	for k, v := range paramMap {
		params.Add(k, v)
	}
	queryString := params.Encode()

	var resp *http.Response
	for _, v := range urls {
		resp, _ = http.Get(fmt.Sprintf("%s?%s", v, queryString))
	}
	defer resp.Body.Close()
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func getCURL(c *gin.Context) string {
	request := c.Request

	protocol := request.URL.Scheme
	if strings.Contains(request.URL.Scheme, "http") == false {
		protocol = "http"
	}

	curl := fmt.Sprintf(`curl '%s' \`, fmt.Sprintf("%s://%s%s", protocol, request.Host, request.URL.String()))

	curl += fmt.Sprintf(`-X '%s' \`, request.Method)

	headers := request.Header
	for key, values := range headers {
		for _, value := range values {
			curl += fmt.Sprintf(`-H '%s: %s' \`, key, value)
		}
	}

	if request.ContentLength > 0 {
		body, exists := c.Get(CtxBody)
		if exists {
			m, err := structToMap(body)
			if err == nil {
				var jsonBytes []byte
				jsonBytes, err = json.Marshal(m)
				if err == nil {
					curl += fmt.Sprintf(`--data-raw '%s' \`, string(jsonBytes))
				}
			}
		}

	}

	return curl
}

func structToMap(s any) (map[string]any, error) {
	rt := reflect.TypeOf(s)
	rv := reflect.ValueOf(s)

	m := make(map[string]any)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		value := rv.Field(i)

		lowerCaseFieldName := strings.ToLower(field.Name[:1]) + field.Name[1:]
		m[lowerCaseFieldName] = value.Interface()
	}

	return m, nil
}
