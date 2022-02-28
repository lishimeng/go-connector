package mq

import (
	"encoding/json"
	"github.com/lishimeng/go-log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AmqpUpstream struct {
	ch       *amqp.Channel
	exchange string
}

// Submit 如返回中有异常应视为发送失败了
func (au *AmqpUpstream) Submit(topic string, v interface{}) (err error) {
	log.Info("submit to %s\n", topic)
	bs, err := json.Marshal(v)
	if err != nil {
		log.Info(err)
		return
	}
	log.Info("submit:" + string(bs))
	err = au.ch.Publish(au.exchange, topic, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        bs,
	})
	if err != nil {
		log.Info("submit failed")
		log.Info(err)
		return
	}
	log.Info("submit completed")
	return
}
