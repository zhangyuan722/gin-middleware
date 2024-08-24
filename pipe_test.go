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
