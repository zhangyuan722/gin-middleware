package gm

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

func QueryPipe[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var instance T

		reflectInstance := reflect.ValueOf(&instance).Elem()

		if err := c.ShouldBindQuery(reflectInstance.Addr().Interface()); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"code": 1000, "msg": err.Error()})
			c.Abort()
			return
		}

		c.Set(CtxQuery, instance)
		c.Next()
	}
}

func BodyPipe[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var instance T

		reflectInstance := reflect.ValueOf(&instance).Elem()

		if err := c.ShouldBindBodyWithJSON(reflectInstance.Addr().Interface()); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"code": 1000, "msg": err.Error()})
			c.Abort()
			return
		}

		c.Set(CtxBody, instance)
		c.Next()
	}
}
