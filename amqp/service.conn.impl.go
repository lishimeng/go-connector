package mq

import (
	"github.com/lishimeng/go-log"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (s *amqpService) newCh() (ch *amqp.Channel, err error) {

	ch, err = s.conn.Channel()
	if err == nil {
		s.chs = append(s.chs, ch)
	} else {
		log.Info("create ch failed")
		log.Info(err)
	}
	return
}

func (s *amqpService) initConn(addr string, onConnected func(connection *amqp.Connection)) (err error) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		log.Info("connect failed")
		log.Info(err)
		return
	}

	log.Info("connect success")
	s.conn = conn
	if onConnected != nil {
		go onConnected(conn)
	}
	go s.connectionWatch()
	return
}

func (s *amqpService) connectionWatch() {
	defer func() {
		log.Info("connection watch process done")
	}()

	for {
		select {
		case <-s.ctx.Done():
			if s.conn != nil {
				log.Info("close connection")
				_ = s.conn.Close()
			}
			return
		}
	}
}
