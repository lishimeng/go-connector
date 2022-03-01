package amqp

import (
	"context"
)

const (
	defaultExchange = "amq.direct"
	ExchangeDirect  = "amq.direct"
	ExchangeFanout  = "amq.fanout"
	ExchangeTopic   = "amq.topic"
)

func Start(ctx context.Context, c Connector, s ...Downstream) (svs Service, err error) {
	exchange := defaultExchange
	if c.Exchange == ExchangeDirect || c.Exchange == ExchangeFanout || c.Exchange == ExchangeTopic {
		exchange = c.Exchange
	}

	upstream := AmqpUpstream{
		exchange: exchange,
	}

	svs = New(ctx, c, &upstream, s...)

	err = svs.Run()
	return
}
