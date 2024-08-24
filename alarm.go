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

func FeishuBotWebHookAlarm(webHookURL string) gin.HandlerFunc {
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
			paramMap["user"] = u.(Claims).ID
		}

		paramMap["status"] = strconv.FormatInt(int64(status), 10)

		jsonData := make(map[string]any)
		_ = json.Unmarshal(w.body.Bytes(), &jsonData)
		respBodyString, _ := json.Marshal(jsonData)

		paramMap["response"] = string(respBodyString)

		paramMap["timing"] = strconv.FormatInt(nowAt.Sub(beginAt).Milliseconds(), 10)

		query := url.Values{}
		for k, v := range paramMap {
			query.Add(k, v)
		}
		query.Add("curl", getCURL(c))
		queryString := query.Encode()

		fmt.Println("Encoded query string:", queryString)

		resp, _ := http.Get(fmt.Sprintf("%s?%s", webHookURL, queryString))
		defer resp.Body.Close()
	}
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
