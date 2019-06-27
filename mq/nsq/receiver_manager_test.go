package nsq

import (
	"testing"
	"github.com/nsqio/go-nsq"
	"fmt"
)

func TestReceiverManager(t *testing.T) {
	t.Skip("no run this test")
	receiver := NewMqReceiver(&MqHostConfigs{Lookup:[]string{"127.0.0.1:4161", "127.0.0.1:4162"}, Nsq: []string{"127.0.0.1:4150", "127.0.0.1:4152"}})
	manager := NewReceiverManager(receiver)
	manager.Add(&testNsqHandler{})
	manager.Start()
	select {}
}

type testNsqHandler struct {

}

func (*testNsqHandler) GenTask() (topic, channel string, handler nsq.HandlerFunc, host *MqHostConfigs) {
	return "topic", "1", func(message *nsq.Message) error {
		fmt.Println(string(message.Body))
		return nil
	}, nil
}
