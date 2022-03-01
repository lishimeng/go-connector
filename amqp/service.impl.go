package amqp

import (
	"context"
	"github.com/lishimeng/go-log"
	amqp "github.com/rabbitmq/amqp091-go"

	"time"
)

var (
	defaultAutoReconnectTimes = 65535
)

type amqpService struct {
	ctxParent        context.Context // 外部系统传进来的关闭信号
	ctx              context.Context // 调用Close关闭后的信号
	cancel           context.CancelFunc
	onLostConnection chan byte // 连接丢失后的信号

	upstream    *AmqpUpstream
	downStreams []Downstream

	conn *amqp.Connection
	chs  []*amqp.Channel

	broker string

	autoReconnect     bool
	maxReconnectTimes int
}

func New(ctx context.Context, c Connector, upstream *AmqpUpstream, s ...Downstream) Service {
	var svs Service

	as := amqpService{
		ctxParent:         ctx,
		upstream:          upstream,
		downStreams:       s,
		conn:              nil,
		chs:               nil,
		broker:            c.Conn,
		maxReconnectTimes: defaultAutoReconnectTimes,
	}
	as.autoReconnect = true // 自动重连开启

	svs = &as
	return svs
}

func (s *amqpService) Close() {
	defer func() {
		recover()
	}()
	s.cancel()
}

func (s *amqpService) Run() (err error) {
	err = s.init()
	if err != nil {
		return
	}
	return
}

func (s *amqpService) UpstreamHandler() Upstream {
	return s.upstream
}

func (s *amqpService) init() (err error) {
	log.Info("init service")
	log.Info("connect mq broker %s", s.broker)

	ctxMe, cancel := context.WithCancel(s.ctxParent)
	s.ctx = ctxMe
	s.cancel = cancel

	s.onLostConnection = make(chan byte)
	var amqpChs []*amqp.Channel
	s.chs = amqpChs // 使用全新的列表保存所有连接

	err = s.initConn(s.broker, func(connection *amqp.Connection) {
		// init upstream handler
		err = s.initUpstream()
		if err != nil {
			return
		}

		// init downstream handler
		s.initDownstream()

		go s.destroy() // watch
		return
	})

	return
}

func (s *amqpService) destroy() {
	defer func() {
		log.Info("amqp watch process done")
	}()

	for {
		select {
		case <-s.ctxParent.Done(): // parent context close
			s.Close()
			return
		case <-s.onLostConnection:
			s.Close()
			if s.autoReconnect {
				go s.reconnect()
			}
			return
		}
	}
}

func (s *amqpService) reconnect() {
	log.Info("reconnect")
	var err error
	var times = 0
	for {
		err = s.init()
		if err == nil {
			return
		}
		times++
		if times > s.maxReconnectTimes {
			return
		}
		log.Info("wait %d try again\n", times)
		time.Sleep(time.Second * time.Duration(times))
	}
}
