package nsq

import (
	"github.com/nsqio/go-nsq"
	"time"
)

type MsgHandlerInterface interface {
	HandleMessage(msg *nsq.Message) error
}

func InitConsumer(topic string, channel string, address string, handler MsgHandlerInterface) {
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = time.Second          //设置重连时间
	c, err := nsq.NewConsumer(topic, channel, cfg) // 新建一个消费者
	if err != nil {
		panic(err)
	}
	c.SetLogger(nil, 0)

	c.AddHandler(handler)

	if err := c.ConnectToNSQLookupd(address); err != nil {
		panic(err)
	}
}
