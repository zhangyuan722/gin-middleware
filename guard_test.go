package m

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthGuard(t *testing.T) {
	r := gin.New()

	authGuard := AuthGuard(&AuthGuardPayload{
		RequestHeaderAuthKey: requestHeaderAuthKey,
		WhiteList:            whiteList,
		SecretKey:            secretKey,
		CallBack: func(c *gin.Context, claims *Claims, accessToken string) {
			c.JSON(http.StatusOK, gin.H{"msg": "Authenticated"})
		},
		FailBack: func(c *gin.Context, code int, msg string, data any) {
			c.JSON(http.StatusOK, gin.H{"msg": msg})
		},
	})

	r.GET("/api/protected", authGuard, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "This is a protected endpoint"})
	})

	whiteListPath := "/api/public"
	whiteListReq := httptest.NewRequest(http.MethodGet, whiteListPath, nil)
	whiteListRecorder := httptest.NewRecorder()
	r.ServeHTTP(whiteListRecorder, whiteListReq)
	if whiteListRecorder.Code != http.StatusOK {
		t.Errorf("Expected 200 for whitelisted path, got %d", whiteListRecorder.Code)
	}

	noTokenReq := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	noTokenRecorder := httptest.NewRecorder()
	r.ServeHTTP(noTokenRecorder, noTokenReq)
	if noTokenRecorder.Code != http.StatusOK {
		t.Errorf("Expected 200 for no token case, got %d", noTokenRecorder.Code)
	}

}
