package gm

import (
	"github.com/gin-gonic/gin"
	"github.com/mileusna/useragent"
	"net/http"
	"reflect"
	"strings"
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

func UserAgentPipe() gin.HandlerFunc {
	return func(c *gin.Context) {
		agent := useragent.Parse(c.Request.UserAgent())

		c.Set(CtxBrowserName, agent.Name)
		c.Set(CtxOSName, agent.OS)
		c.Set(CtxOSVersion, agent.OSVersion)

		var (
			deviceType string
		)
		if agent.Desktop {
			deviceType = DeviceTypeDesktop
		} else if agent.Mobile {
			deviceType = DeviceTypeMobile
		} else if agent.Tablet {
			deviceType = DeviceTypeTablet
		} else if agent.Bot {
			deviceType = DeviceTypeBot
		} else {
			if agent.Device != "" {
				deviceType = strings.ToLower(agent.Device)
			} else {
				deviceType = DeviceTypeUnknown
			}
		}
		c.Set(CtxDeviceType, deviceType)

		c.Next()
	}
}
