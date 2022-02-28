package mq

import "log"

func (s *amqpService) initUpstream() (err error) {
	log.Println("init upstream")
	ch, err := s.newCh()
	if err != nil {
		log.Println("upstream failed")
		log.Println(err)
		return
	}

	s.upstream.ch = ch
	return
}
