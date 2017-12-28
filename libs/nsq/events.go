/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd.
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
 *     Initial: 2017/12/28        Feng Yifei
 */

package nsq

import (
	"errors"

	gnsq "github.com/nsqio/go-nsq"
)

// Events is a abstraction of a basic event bus.
type Events struct {
	nsq      string
	topic    string
	config   *gnsq.Config
	consumer []*gnsq.Consumer
	producer *gnsq.Producer
}

// NewEvents creates a new event bus.
func NewEvents(nsq string, topic string) (ev *Events, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	if len(nsq) == 0 {
		nsq = "127.0.0.1:4150"
	}

	if !gnsq.IsValidTopicName(topic) {
		err = errors.New("Invalid topic name")
		return
	}

	ev = &Events{
		nsq:      nsq,
		topic:    topic,
		config:   gnsq.NewConfig(),
		consumer: []*gnsq.Consumer{},
	}

	ev.Set("max_in_flight", 10)

	if err = ev.CreateProducer(); err != nil {
		ev = nil
		return
	}

	return
}

// Set wraps the Config.Set method.
func (ev *Events) Set(key string, value interface{}) error {
	return ev.config.Set(key, value)
}

// CreateProducer creates a producer.
func (ev *Events) CreateProducer() error {
	if ev.producer != nil {
		return errors.New("producer already exists")
	}

	producer, err := gnsq.NewProducer(ev.nsq, ev.config)
	if err != nil {
		return err
	}

	ev.producer = producer
	return nil
}

// CreateConsumer creates a new consumer.
func (ev *Events) CreateConsumer(ch string, handler gnsq.Handler) error {
	c, err := gnsq.NewConsumer(ev.topic, ch, ev.config)

	if err != nil {
		return err
	}

	c.AddHandler(handler)

	if err = c.ConnectToNSQD(ev.nsq); err != nil {
		return err
	}

	ev.consumer = append(ev.consumer, c)

	return nil
}

// Publish a message.
func (ev *Events) Publish(message []byte) error {
	return ev.producer.Publish(ev.topic, message)
}

// Close the event.
func (ev *Events) Close() {
	if ev.producer != nil {
		ev.producer.Stop()
		ev.producer = nil
	}

	for _, c := range ev.consumer {
		c.DisconnectFromNSQD(ev.nsq)
	}
}
