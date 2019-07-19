package nsq

import (
	"testing"
	"github.com/nsqio/go-nsq"
	"fmt"
)

func TestReceiverManager(t *testing.T) {
	t.Skip("no run this test")
	SetLogLv(nsq.LogLevelError)
	receiver := NewNsqReceiver(&MqHostConfigs{Lookup: []string{"127.0.0.1:4161", "127.0.0.1:4162"}, Nsq: []string{"127.0.0.1:4150", "127.0.0.1:4152"}})
	manager := NewReceiverManager(receiver)
	manager.Add(GenTask("Topic", "Channel"))
	manager.Start()
	select {}
}

func GenTask(topic, channel string) NsqHandlerFunc {
	return func() (config *MqTaskConfigs) {
		return &MqTaskConfigs{
			Topic:   topic,
			Channel: channel,
			Handler: func(message *nsq.Message) error {
				fmt.Println(string(message.Body))
				return nil
			},
		}
	}
}

func TestReceiverManagerWithMultipleTask(t *testing.T) {
	t.Skip("no run this test")
	receiver := NewNsqReceiver(&MqHostConfigs{Lookup: []string{"127.0.0.1:4161", "127.0.0.1:4162"}, Nsq: []string{"127.0.0.1:4150", "127.0.0.1:4152"}})
	task := NewCreateOrder()
	manager := NewReceiverManager(receiver)
	manager.Add(task.tmp1Handler("1", "2"), task.tmp2Handler("3", "4"))
	manager.Start()
	select {}
}

type NsqCreateOrder interface {
	tmp1Handler(topic, channel string) NsqHandlerFunc
	tmp2Handler(topic, channel string) NsqHandlerFunc
}

func NewCreateOrder() NsqCreateOrder {
	return &TestCreateOrder{db: nil, logic: nil}
}

type TestCreateOrder struct {
	db    interface{}
	logic interface{}
}

func (*TestCreateOrder) tmp1Handler(topic, channel string) NsqHandlerFunc {
	return func() (config *MqTaskConfigs) {
		return &MqTaskConfigs{
			Topic:   topic,
			Channel: channel,
			Handler: func(message *nsq.Message) error {
				fmt.Println(string(message.Body))
				return nil
			},
		}
	}
}

func (*TestCreateOrder) tmp2Handler(topic, channel string) NsqHandlerFunc {
	return func() (config *MqTaskConfigs) {
		return &MqTaskConfigs{
			Topic:   topic,
			Channel: channel,
			Handler: func(message *nsq.Message) error {
				fmt.Println(string(message.Body))
				return nil
			},
		}
	}
}