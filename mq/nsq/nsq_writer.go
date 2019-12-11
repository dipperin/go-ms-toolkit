package nsq

import (
	"errors"

	"github.com/dipperin/go-ms-toolkit/json"
	qylog "github.com/dipperin/go-ms-toolkit/log"
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
)

type NsqWriter interface {
	Refresh()
	Publish(topic string, jsonObj interface{}) error
	PublishString(topic string, msg string) error
	Stop()
}

type BaseNsqWriter struct {
	addrs     []string
	producers []*nsq.Producer
}

func NewNsqWriter(addrs []string) NsqWriter {
	if len(addrs) <= 0 {
		panic("nsq addrs length is 0")
	}
	return (&BaseNsqWriter{addrs: addrs}).newProducers()
}

func (writer *BaseNsqWriter) newProducers() NsqWriter {
	for _, addr := range writer.addrs {
		writer.newProducer(addr)
	}
	if len(writer.producers) <= 0 {
		// panic? or do refresh? or error handler?
		panic("NsqWriter.producers length is 0")
	}
	return writer
}

func (writer *BaseNsqWriter) Stop() {
	for _, producer := range writer.producers {
		producer.Stop()
	}
}

func (writer *BaseNsqWriter) Refresh() {
	refreshed := writer.refreshProducer()
	if len(refreshed) > 0 {
		writer.producers = nil
		copy(writer.producers, refreshed)
	}
}

func (writer *BaseNsqWriter) Publish(topic string, jsonObj interface{}) error {
	if len(writer.producers) <= 0 {
		return errors.New("no producer on topic: '" + topic + "'")
	}
	return writer.pubMsg(topic, json.StringifyJsonToBytes(jsonObj))
}

func (writer *BaseNsqWriter) PublishString(topic string, msg string) error {
	if len(writer.producers) <= 0 {
		return errors.New("no producer on topic: '" + topic + "'")
	}
	return writer.pubMsg(topic, []byte(msg))
}

func (writer *BaseNsqWriter) newProducer(addr string) {
	producer, err := nsq.NewProducer(addr, nsq.NewConfig())

	if err != nil {
		qylog.QyLogger.Error("NsqWriter new nsq producer failed", zap.String("addr", addr), zap.Error(err))
		return
	}

	if err := producer.Ping(); err != nil {
		qylog.QyLogger.Error("NsqWriter nsq producer ping check failed",
			zap.String("addr", addr), zap.Error(err))
		return
	}

	writer.producers = append(writer.producers, producer)
}

func (writer *BaseNsqWriter) refreshProducer() (refreshedProducers []*nsq.Producer) {
	for _, addr := range writer.addrs {
		producer, err := nsq.NewProducer(addr, nsq.NewConfig())

		if err != nil {
			qylog.QyLogger.Error("NsqOrderWriter new nsq producer failed", zap.String("addr", addr), zap.Error(err))
			return
		}

		if err := producer.Ping(); err != nil {
			qylog.QyLogger.Error("NsqOrderWriter nsq producer ping check failed",
				zap.String("addr", addr), zap.Error(err))
			return
		}

		refreshedProducers = append(writer.producers, producer)
	}
	return
}

func (writer *BaseNsqWriter) pubMsg(topic string, msg []byte) error {

	// todo, 随机取
	for i := range writer.producers {
		if err := writer.producers[i].Publish(topic, msg); err != nil {
			qylog.QyLogger.Error("NsqWriter nsq producer publish msg failed", zap.String("topic", topic),
				zap.String("addrs", writer.producers[i].String()), zap.String("err", err.Error()))
			continue
		}

		return nil
	}

	return errors.New("NsqWriter all nsq producers publish failed on topic: '" + topic + "'")
}
