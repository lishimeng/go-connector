package amqp

// Connector Exchange: 'amq.direct', 'amq.fanout', 'amq.topic'
type Connector struct {
	Conn     string
	Exchange string
}

type Downstream interface {
	Subscribe(topic string, v interface{}, upstream Upstream)
	Topic() string
}

type Upstream interface {
	Submit(topic string, v interface{}) error // 提交
}

type Service interface {
	Run() error
	UpstreamHandler() Upstream
	Close()
}
