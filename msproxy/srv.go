package msproxy

import (
	"errors"
	"github.com/dipperin/go-ms-toolkit/log"
	"go.uber.org/zap"
)

type Context interface {
	SetResult(data interface{})
}

type IProxyResp interface {
	GetErrMsg() string
	GetData() interface{}
	GetSuccess() bool
}

type Srv struct {
	requester Requester
}

func NewSrv(requester Requester) *Srv {
	return &Srv{requester: requester}
}

func (s *Srv) PostProxy(api string, ctx Context, resp IProxyResp) error {
	if err := s.requester.Post(api, ctx, resp); err != nil {
		log.QyLogger.Warn("post proxy failed", zap.String("api", api), zap.Any("req", ctx), zap.Error(err))
		return err
	}

	if !resp.GetSuccess() {
		log.QyLogger.Warn("post proxy resp success false", zap.Any("resp", resp), zap.String("api", api), zap.Any("req", ctx))
		return errors.New(resp.GetErrMsg())
	}

	return nil
}
