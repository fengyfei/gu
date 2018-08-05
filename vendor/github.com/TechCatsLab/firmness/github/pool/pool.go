/*
 * Revision History:
 *     Initial: 2018/07/04        Li Zebang
 */

package pool

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// Pool -
type Pool interface {
	Get(string) (*Client, error)
	Put(*Client) error
	AddClient(tokens ...*Token) error
	ResetInterval(time.Duration)
}

type pool struct {
	lock            *sync.Mutex
	clients         map[string][]*Client
	resetInterval   time.Duration
	latestResetTime time.Time
	index           int
	timers          []*timer
	done            chan struct{}
}

type timer struct {
	clients []*Client
}

var (
	// MinResetInterval -
	MinResetInterval = time.Minute
	// ResetPeriod -
	ResetPeriod = time.Hour
	// ResetIntervalOffset -
	ResetIntervalOffset = 1.2
	// MinReuseNum -
	MinReuseNum = 10
)

// NewPool -
func NewPool(interval time.Duration, tokens ...*Token) (Pool, error) {
	for index := 0; index < len(tokens); index++ {
		if tokens[index] == nil {
			log.Printf("warning: the token with index '%d' is nil", index)
			tokens = append(tokens[:index], tokens[index+1:]...)
			index--
		}
	}

	if len(tokens) == 0 {
		return nil, errors.New("the number of tokens cannot be 0")
	}

	if interval < MinResetInterval {
		interval = MinResetInterval
	}

	ts := make(map[string][]*Token)
	for index := 0; index < len(tokens); index++ {
		if _, exist := ts[tokens[index].Tag]; !exist {
			ts[tokens[index].Tag] = make([]*Token, 0)
		}
		ts[tokens[index].Tag] = append(ts[tokens[index].Tag], tokens[index])
	}

	cs := make(map[string][]*Client)
	for tag, tokens := range ts {
		cs[tag] = make([]*Client, len(tokens))
		for index := 0; index < len(tokens); index++ {
			cs[tag][index], _ = client(tokens[index])
		}
	}

	var timers []*timer
	if ResetPeriod%interval != 0 {
		timers = make([]*timer, ResetPeriod/interval+1)
	} else {
		timers = make([]*timer, ResetPeriod/interval+1)
	}
	for index := range timers {
		timers[index] = &timer{clients: make([]*Client, 0)}
	}

	p := &pool{
		clients:         cs,
		resetInterval:   interval,
		latestResetTime: time.Now(),
		lock:            &sync.Mutex{},
		timers:          timers,
		done:            make(chan struct{}),
	}

	go p.resetClient()

	return p, nil
}

func (p *pool) resetClient() {
	timer := time.NewTimer(p.resetInterval)
	for {
		select {
		case t := <-timer.C:
			func() {
				p.lock.Lock()
				defer p.lock.Unlock()
				p.latestResetTime = t.Local()
				if p.index == len(p.timers) {
					p.index = 0
				}
				for _, client := range p.timers[p.index].clients {
					client.Reset()
					p.clients[client.tag] = append(p.clients[client.tag], client)
				}
				p.index++
				timer.Reset(p.resetInterval)
			}()
		case <-p.done:
			return
		}
	}
}

func (p *pool) Get(tag string) (*Client, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if _, exist := p.clients[tag]; !exist {
		return nil, fmt.Errorf("the client with tag '%s' doesn't exist", tag)
	}

	if len(p.clients[tag]) == 0 {
		return nil, errors.New("no client available")
	}

	client := p.clients[tag][len(p.clients[tag])-1]
	if client.remaining <= MinReuseNum && client.remaining != defaultRemaining {
		p.clients[tag] = p.clients[tag][:len(p.clients[tag])-1]
	}

	return client, nil
}

func (p *pool) Put(client *Client) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if client == nil {
		return errors.New("client is nil")
	}

	if _, exist := p.clients[client.tag]; !exist {
		return errors.New("the client is invalid ")
	}

	if !client.isUsed && client.remaining != defaultRemaining {
		client.isUsed = true
		sub := client.reset.Sub(time.Now())
		offset := int(sub / p.resetInterval)
		if p.latestResetTime.Add(p.resetInterval).Sub(time.Now()).Seconds() < (sub%p.resetInterval).Seconds()*ResetIntervalOffset {
			offset++
		}
		index := (p.index + offset - 1) % len(p.timers)
		p.timers[index].clients = append(p.timers[index].clients, client)
	}

	return nil
}

func (p *pool) AddClient(tokens ...*Token) error {
	return nil
}

func (p *pool) ResetInterval(interval time.Duration) {

}
