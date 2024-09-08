package gm

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type TestQuery struct {
	Id   int    `form:"id" binding:"required"`
	Name string `form:"name" binding:"required"`
}

type TestBody struct {
	Message string `json:"message" binding:"required"`
}

func TestQueryPipe(t *testing.T) {
	r := gin.Default()

	r.GET("/test/query", QueryPipe[TestQuery](), func(c *gin.Context) {
		query := c.MustGet(CtxQuery).(TestQuery)
		if query.Id != 1 || query.Name != "test" {
			t.Errorf("QueryPipe did not bind query parameters correctly")
		}
		c.JSON(http.StatusOK, gin.H{"msg": "Query parameters are correct"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test/query?id=1&name=test", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %v", w.Code)
	}
}

func TestBodyPipe(t *testing.T) {
	r := gin.Default()

	r.POST("/test/body", BodyPipe[TestBody](), func(c *gin.Context) {
		body := c.MustGet(CtxBody).(TestBody)
		if body.Message != "test" {
			t.Errorf("BodyPipe did not bind body correctly")
		}
		c.JSON(http.StatusOK, gin.H{"msg": "Body is correct"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/test/body", strings.NewReader(`{"message": "test"}`))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %v", w.Code)
	}
}

func TestUserAgentPipe(t *testing.T) {
	r := gin.Default()

	type It struct {
		BrowserName string
		OSName      string
		OSVersion   string
		DeviceType  string
	}

	tests := map[string]It{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8": {
			BrowserName: "Safari",
			OSName:      "macOS",
			OSVersion:   "10.12.6",
			DeviceType:  DeviceTypeDesktop,
		},
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36": {
			BrowserName: "Chrome",
			OSName:      "Windows",
			OSVersion:   "6.1",
			DeviceType:  DeviceTypeDesktop,
		},
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1": {
			BrowserName: "Safari",
			OSName:      "iOS",
			OSVersion:   "10.3.2",
			DeviceType:  DeviceTypeMobile,
		},
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) FxiOS/8.1.1b4948 Mobile/14F89 Safari/603.2.4": {
			BrowserName: "Firefox",
			OSName:      "iOS",
			OSVersion:   "10.3.2",
			DeviceType:  DeviceTypeMobile,
		},
		"Mozilla/5.0 (iPad; CPU OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1": {
			BrowserName: "Safari",
			OSName:      "iOS",
			OSVersion:   "10.3.2",
			DeviceType:  DeviceTypeTablet,
		},
		"Mozilla/5.0 (Linux; Android 4.3; GT-I9300 Build/JSS15J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.125 Mobile Safari/537.36": {
			BrowserName: "Chrome",
			OSName:      "Android",
			OSVersion:   "4.3",
			DeviceType:  DeviceTypeMobile,
		},
		"Mozilla/5.0 (Android 4.3; Mobile; rv:54.0) Gecko/54.0 Firefox/54.0": {
			BrowserName: "Firefox",
			OSName:      "Android",
			OSVersion:   "4.3",
			DeviceType:  DeviceTypeMobile,
		},
		"Mozilla/5.0 (Linux; Android 4.3; GT-I9300 Build/JSS15J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.91 Mobile Safari/537.36 OPR/42.9.2246.119956": {
			BrowserName: "Opera",
			OSName:      "Android",
			OSVersion:   "4.3",
			DeviceType:  DeviceTypeMobile,
		},
		"Opera/9.80 (Android; Opera Mini/28.0.2254/66.318; U; en) Presto/2.12.423 Version/12.16": {
			BrowserName: "Opera Mini",
			OSName:      "Android",
			OSVersion:   "",
			DeviceType:  DeviceTypeMobile,
		},
	}

	r.GET("/test/useragent", UserAgentPipe(), func(c *gin.Context) {
		browserName := c.MustGet(CtxBrowserName)
		deviceType := c.MustGet(CtxDeviceType)
		osName := c.MustGet(CtxOSName)
		osVersion := c.MustGet(CtxOSVersion)

		it := tests[c.Request.UserAgent()]

		if browserName != it.BrowserName {
			t.Errorf("browserName expected %s, but got %s", it.BrowserName, browserName)
		}

		if osName != it.OSName {
			t.Errorf("osName expected %s, but got %s", it.OSName, osName)
		}

		if osVersion != it.OSVersion {
			t.Errorf("osVersion expected %s, but got %s", it.OSVersion, osVersion)
		}

		if deviceType != it.DeviceType {
			t.Errorf("deviceType expected %s, but got %s", it.DeviceType, deviceType)
		}

		c.String(http.StatusOK, "ok")
	})

	for k := range tests {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test/useragent", nil)
		req.Header.Set("User-Agent", k)
		r.ServeHTTP(w, req)
	}
}
