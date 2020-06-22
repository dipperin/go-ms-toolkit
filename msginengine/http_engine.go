package msginengine

import (
	"bytes"
	"github.com/dipperin/go-ms-toolkit/log"
	"github.com/dipperin/go-ms-toolkit/qyenv"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"strconv"
)

// 鉴权
type authorization interface {
	Auth(userID uint, token string) error
}

// 数据加解密
type reqDecrypt interface {
	PostBodyDecrypt(body string) ([]byte, error)
}

type IGRouter interface {
	SetAuthRg(authRg *gin.RouterGroup)
	SetRg(rg *gin.RouterGroup)
	InitRouter()
}

type Engine struct {
	ginEngine   *gin.Engine
	auth        authorization
	decryptUtil reqDecrypt
}

func NewEngine(ginEngine *gin.Engine, opts ...EngineOption) *Engine {
	s := &Engine{ginEngine: ginEngine}
	return s.WithOptions(opts...)
}

func ReleasesEngine() *Engine {
	// release 默认启用的选项: 监控检测， JWT验证
	options := []EngineOption{EnablePingCheck()}

	switch qyenv.GetUseDocker() {
	// 本地 & 开发 & 测试 现在默认打印请求
	case 0, 1:
		options = append(options, PrintReq(), PrintResp())
	case 2:
		if os.Getenv("ENGINE_PRINT") == "1" {
			options = append(options, PrintReq(), PrintResp())
		}
	}

	return NewEngine(gin.New(), options...)
}

func (e *Engine) WithOptions(opts ...EngineOption) *Engine {
	for _, opt := range opts {
		opt.apply(e)
	}
	return e
}

func (e *Engine) Start(port string) {
	log.QyLogger.Info(`server start`, zap.String("port", port))

	if err := e.ginEngine.Run(":" + port); err != nil {
		panic(err)
	}
}

// 如果使用 这个函数的话会自动加载Group，请注意！！！！
func (e *Engine) EasyCombine(group string, grs ...IGRouter) *Engine {
	rg := e.NewRouterGroup(group)
	authRg := e.NewAuthRouterGroup(group)

	for i := range grs {
		grs[i].SetRg(rg)
		grs[i].SetAuthRg(authRg)
		grs[i].InitRouter()
	}

	return e
}

func (e *Engine) NewRouterGroup(group string) *gin.RouterGroup {
	return e.ginEngine.Group(group)
}

func (e *Engine) NewAuthRouterGroup(group string) *gin.RouterGroup {
	rg := e.ginEngine.Group(group)
	rg.Use(e.authMiddleware())
	return rg
}

func (e *Engine) authMiddleware() gin.HandlerFunc {
	if e.auth == nil {
		return func(c *gin.Context) {}
	}

	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		id := c.Request.Header.Get("userid")

		uid, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(403, gin.H{"success": false, "err_msg": "Permission denied"})
			c.Abort()
			return
		}

		if err := e.auth.Auth(uint(uid), token); err == nil {
			c.Set("userid", id)
			return
		}

		c.JSON(403, gin.H{"success": false, "err_msg": "Permission denied"})
		c.Abort()
		return
	}
}

func (e *Engine) decryptMiddleware() gin.HandlerFunc {
	if e.decryptUtil == nil {
		return func(c *gin.Context) {}
	}

	return func(c *gin.Context) {
		// 不处理有文件的请求
		if c.ContentType() == "multipart/form-data" {
			return
		}

		params, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body.Close()

		result, err := e.decryptUtil.PostBodyDecrypt(string(params))
		if err != nil {

			log.QyLogger.Error("decrypt failed", zap.String("raw msg", string(params)), zap.Error(err))

			c.JSON(500, gin.H{"success": false, "err_msg": "server can't handle requester"})
			c.Abort()
			return
		}

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(result))
		return
	}
}

func (e *Engine) ping(handler gin.HandlerFunc) {
	e.ginEngine.GET("/monitor/ping", handler)
}
