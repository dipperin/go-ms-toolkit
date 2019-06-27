package nsq

import (
	"github.com/nsqio/go-nsq"
)

type NsqHandler interface {
	GenTask() (topic, channel string, handler nsq.HandlerFunc)
}

type NsqHandlerFunc func() (topic, channel string, handler nsq.HandlerFunc)

func (f NsqHandlerFunc) GenTask() (topic, channel string, handler nsq.HandlerFunc) {
	return f()
}

type ReceiverManager struct {
	receiver MqReceiver
}

// todo, 多个receiver?
func NewReceiverManager(receiver MqReceiver, h ...NsqHandler) *ReceiverManager {
	rm := &ReceiverManager{ receiver: receiver }
	rm.Add(h...)
	return rm
}

func (a *ReceiverManager) Add(h ...NsqHandler) {
	for i := range h {
		a.receiver.AddTask(NewNsqTask(h[i].GenTask()))
	}
}

func (a *ReceiverManager) Start() {
	a.receiver.Start()
}