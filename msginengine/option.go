package msginengine

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type EngineOption interface {
	apply(*Engine)
}

type optionFunc func(*Engine)

func (f optionFunc) apply(s *Engine) {
	f(s)
}

func defaultPingCheck(c *gin.Context) {
	c.JSON(http.StatusOK, "OK")
}

func EnablePingCheck() EngineOption {
	return optionFunc(func(s *Engine) {
		s.ping(defaultPingCheck)
	})
}

// 开启请求body解密
func EnablePostBodyDecrypt(d reqDecrypt) EngineOption {
	return optionFunc(func(s *Engine) {
		if d == nil {
			return
		}
		s.decryptUtil = d
		s.ginEngine.Use(s.decryptMiddleware())
	})
}

// 开启jwt 验证
func EnableJwtAuth() EngineOption {
	return optionFunc(func(s *Engine) {
		//jwtAuth := jwt.NewSrv()
		//s.auth = jwtAuth
	})
}

// 打印请求
func PrintReq() EngineOption {
	return optionFunc(func(s *Engine) {
		s.ginEngine.Use(printReq())
	})
}

// 打印响应
func PrintResp() EngineOption {
	return optionFunc(func(s *Engine) {
		s.ginEngine.Use(printResp())
	})
}

