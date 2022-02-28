package mq

import (
	"github.com/lishimeng/go-log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type DownstreamHandler struct {
	ds               *Downstream
	onLostConnection chan byte
}

func (s *amqpService) initDownstream() {
	if len(s.downStreams) > 0 {
		for _, ds := range s.downStreams {
			go s.initSubscribe(ds)
		}
	}
}

func (s *amqpService) initSubscribe(ds Downstream) {
	go s.subscriberProcess(ds)
}

func (s *amqpService) subscriberProcess(ds Downstream) {
	ch, err := s.newCh()
	if err != nil {
		log.Info(err)
		return
	}
	defer func() {
		_ = ch.Close()
	}()

	// TODO transaction(ack)
	log.Info("declare queue:" + ds.Topic())
	q, err := ch.QueueDeclare(
		ds.Topic(), // name
		true,       // durable
		true,       // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Info("queue failed")
		log.Info(err)
		return
	}

	consume, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Info("consume failed")
		log.Info(err)
		return
	}

	err = ch.QueueBind(q.Name, q.Name, s.upstream.exchange, false, nil)
	if err != nil {
		log.Info(err)
		return
	}

	for {
		select {
		case <-s.ctx.Done():
			return
		case m, active := <-consume:
			if !active { // 异常关闭
				go func() {
					s.onLostConnection <- 0x00
				}()
				return
			}
			s.handleMsg(ds, m)
		}
	}
}

func (s *amqpService) handleMsg(ds Downstream, ms amqp.Delivery) {
	log.Info("received downLink stream [exchange:%s router:%s]", ms.Exchange, ms.RoutingKey)
	log.Info(string(ms.Body))
	ds.Subscribe(ds.Topic(), ms.Body, s.upstream)
}
