/*
 * MIT License
 *
 * Copyright (c) 2018 SmartestEE Co., Ltd..
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2018/05/14        Tong Yuehong
 */

package nsq

import (
	"fmt"
	"log"
	"strconv"

	"github.com/nsqio/go-nsq"
	"github.com/pquerna/ffjson/ffjson"
)

type Mq struct {
	Conf *nsq.Config
}

type Consumer struct {
	*nsq.Consumer
}

type Producer struct {
	*nsq.Producer
}

func (mq *Mq) NewConf(setConf ...func(*nsq.Config)) *nsq.Config {
	conf := nsq.NewConfig()
	for _, set := range setConf {
		set(conf)
	}

	mq.Conf = conf
	return conf
}

func (m *Mq) CreateProducer(addr string) (*Producer, error) {
	pro, err := nsq.NewProducer(addr, m.Conf)
	if err != nil {
		return nil, err
	}

	return &Producer{
		pro,
	}, nil
}

func (p *Producer) Pub(topic string, message interface{}) error {
	if len(topic) == 0 {
		return ErrTopicRequired
	}

	err := p.Ping()
	if err != nil {
		log.Fatal("should not be able to ping after Stop()")
		return err
	}

	msg, err := marshalMsg(message)
	if err != nil {
		return err
	}

	err = p.Publish(topic, msg)
	if err != nil {
		return err
	}

	return nil
}

func (p *Producer) StopIt() {
	p.Stop()
}

func (m *Mq) CreateConsumer(topic, channel, addr string, handler nsq.Handler) (*Consumer, error) {
	if len(topic) == 0 {
		return nil, ErrTopicRequired
	}

	if len(channel) == 0 {
		return nil, ErrChannelRequired
	}

	if handler == nil {
		return nil, ErrHandlerRequired
	}

	c, err := nsq.NewConsumer(topic, channel, m.Conf)
	if err != nil {
		return nil, err
	}

	c.AddHandler(handler)
	err = c.ConnectToNSQD(addr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &Consumer{
		c,
	}, nil
}

func (c *Consumer) Stop() {
	c.Stop()
}

func marshalMsg(msg interface{}) (m []byte, err error) {
	switch t := msg.(type) {
	case []byte:
		m = t
	case float64:
		m = []byte(strconv.FormatFloat(t, 'f', -1, 64))
	case int64:
		m = []byte(strconv.FormatInt(t, 10))
	case string:
		m = []byte(t)
	default:
		m, err = ffjson.Marshal(msg)
	}

	return
}
