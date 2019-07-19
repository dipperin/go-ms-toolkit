package nsq

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"testing"
)

type testMsg struct {
	Msg   string `json:"msg"`
	Index int    `json:"index"`
}

func TestNewNsqProducer(t *testing.T) {
	t.Skip("no run this test")

	p := NewNsqProducer([]string{"172.16.0.30:4150"})

	for {
		err := p.PubMsg("test1", &testMsg{Msg: "4324324", Index: 2})

		if err != nil {
			t.Fatal("err", err.Error())
		}
	}

}

func TestInitConsumer(t *testing.T) {
	t.Skip("no run this test")
	c := NewConsumer("test1", "aab")
	c.Set("nsqds", []string{"172.16.0.30:4150", "172.16.0.30:4152"})
	//c.Set("nsqlookupds", []string{"172.16.0.30:4161"})
	c.Set("concurrency", 15)
	c.Set("max_attempts", 10)
	c.Set("max_in_flight", 150)
	err := c.Start(nsq.HandlerFunc(func(msg *nsq.Message) error {
		fmt.Println("customer2:", string(msg.Body))
		return nil
	}))
	if err != nil {
		fmt.Errorf(err.Error())
	}

	select {}
}

//// 消费者
//func startConsumer() {
//	cfg := Nsq.NewConfig()
//	consumer, err := Nsq.NewConsumer("test1", "aab", cfg)
//	if err != nil {
//		log.Fatal(err)
//	}
//	// 设置消息处理函数
//	consumer.AddHandler(Nsq.HandlerFunc(func(message *Nsq.Message) error {
//		log.Println(string(message.Body))
//		return nil
//	}))
//	// 连接到单例nsqd
//	if err := consumer.ConnectToNSQD("172.16.0.30:4150"); err != nil {
//		log.Fatal(err)
//	}
//	<-consumer.StopChan
//}
