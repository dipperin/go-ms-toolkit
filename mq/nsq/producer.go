package nsq

import (
	"errors"
	"github.com/dipperin/go-ms-toolkit/json"
	"github.com/nsqio/go-nsq"
)

type Producer struct {
	producers []*nsq.Producer
}

func (p *Producer) newProduce(addr string) *nsq.Producer {
	producer, err := nsq.NewProducer(addr, nsq.NewConfig())

	if err != nil {
		// todo 这里应该不用panic
		panic(err)
	}

	return producer
}

func (p *Producer) PubMsg(topic string, msg interface{}) error {
	producer, err := p.getProducer()

	if err != nil {
		return err
	}

	if msg == nil {
		return errors.New("msg no nil")
	}

	bMsg, err := json.StringifyJsonToBytesWithErr(msg)
	if err != nil {
		return err
	}

	if err := producer.Publish(topic, bMsg); err != nil {
		return err
	}

	return nil
}

func (p *Producer) getProducer() (*nsq.Producer, error) {
	if len(p.producers) == 0 {
		return nil, errors.New("no producer")
	}

	// todo 以后要随机取
	producer := p.producers[0]

	if producer == nil {
		return nil, errors.New("get nil producer")
	}

	return producer, nil
}
