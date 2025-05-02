package pubsub

import "sync"

type PubSub struct {
	Channs []chan struct{}
	Lock   *sync.RWMutex
}

func NewPubSub() *PubSub {
	return &PubSub{
		Channs: make([]chan struct{}, 0),
		Lock:   new(sync.RWMutex),
	}
}

func (p *PubSub) SubScribe() (<-chan struct{}, func()) {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	c := make(chan struct{}, 1)
	p.Channs = append(p.Channs, c)

	return c, func() {
		p.Lock.Lock()
		defer p.Lock.Unlock()

		for i, channel := range p.Channs {
			if channel == c {
				p.Channs = append(p.Channs[:i], p.Channs[i+1:]...)
				close(c)
				return
			}
		}
	}
}

func (p *PubSub) Publish() {
	p.Lock.RLock()
	defer p.Lock.RUnlock()

	for _, channel := range p.Channs {
		channel <- struct{}{}
	}
}
