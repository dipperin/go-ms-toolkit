package msginrouter

import (
	"errors"
	"github.com/dipperin/go-ms-toolkit/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

var (
	ParamsParseErr = errors.New("params parse error")
)

type R struct {
	rg     *gin.RouterGroup
	authRg *gin.RouterGroup
}

func (s *R) AuthRg() *gin.RouterGroup {
	return s.authRg
}

func (s *R) Rg() *gin.RouterGroup {
	return s.rg
}

func (s *R) SetAuthRg(authRg *gin.RouterGroup) {
	s.authRg = authRg
}

func (s *R) SetRg(rg *gin.RouterGroup) {
	s.rg = rg
}

func (s *R) InitRouter() {
	panic("please overwrite this func")
}

func (s *R) Handle(newCtx func() IContext, run func(ctx IContext) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := newCtx()
		if err := c.BindJSON(ctx); err != nil {
			log.QyLogger.Warn(c.Request.URL.Path+"#bind json failed", zap.Error(err))
			httpRender400(c, ParamsParseErr)
			return
		}

		if id, ok := c.Get("userid"); ok {
			uid, err := strconv.Atoi(id.(string))
			if err != nil {
				httpRender403(c)
				return
			}

			if uid != 0 {
				ctx.SetUserID(uint(uid))
			}
		}

		if err := run(ctx); err != nil {
			log.QyLogger.Warn(c.Request.URL.Path+"#run handler get error", zap.Error(err))
			httpRender400(c, err)
			return
		}

		httpRenderSuccess(c, ctx.GetResult())
		return
	}
}

type BaseResp struct {
	Data    interface{} `json:"data,omitempty"`
	Success bool        `json:"success"`
	ErrMsg  string      `json:"err_msg,omitempty"`
}

func httpRender400(c *gin.Context, err error) {
	c.JSON(400, BaseResp{Success: false, ErrMsg: err.Error()})
}

func httpRenderSuccess(c *gin.Context, data interface{}) {
	c.JSON(200, BaseResp{Success: true, Data: data})
}

func httpRender403(c *gin.Context) {
	c.JSON(403, BaseResp{Success: false, ErrMsg: "Permission denied"})
}
