package msginengine

import (
	"bytes"
	"github.com/dipperin/go-ms-toolkit/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
)

func printReq() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 不打印有文件的请求
		if c.ContentType() == "multipart/form-data" {
			return
		}
		params, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body.Close()
		// 这里有意思，呵呵
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(params))
		paramsStr := string(params)
		log.QyLogger.Info("[Get Request]", zap.String("req_url", c.Request.RequestURI), zap.String("params", paramsStr))
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func printResp() gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		respStr := blw.body.String()
		if len(respStr) < 500 {
			log.QyLogger.Info("[Response Result]", zap.String("req_url", c.Request.RequestURI), zap.String("resp", respStr))
		} else {
			log.QyLogger.Info("[Response Partial Result]", zap.String("req_url", c.Request.RequestURI), zap.String("resp", respStr[0:500]))
		}
	}
}
