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
 *     Initial: 2018/04/20        Li Zebang
 */

package increment

import (
	"errors"
	"sync"
)

type Increment struct {
	IncrementKey   string
	updateKey      string
	path           string
	bucket         string
	signalCh       chan *Signal
	singalChStatus bool
	done           chan struct{}
	once           sync.Once
}

type Signal struct {
	IncrementKey string
	Done         bool
}

var (
	ErrSignalChClosed = errors.New("singal channel is closed")
)

const (
	incrementKey = "IncrementKey"
)

func NewIncrement(path, bucket string) *Increment {
	return &Increment{
		path:     path,
		bucket:   bucket,
		signalCh: make(chan *Signal, 0),
		done:     make(chan struct{}, 0),
		once:     sync.Once{},
	}
}

func (i *Increment) StartIncrement() (err error) {
	db, err := open(i.path)
	if err != nil {
		return err
	}

	value, err := get(db, []byte(i.bucket), []byte(incrementKey))
	if err != nil {
		return err
	}
	i.IncrementKey = string(value)

	go func() {
		defer func() {
			close(i.signalCh)
			err = set(db, []byte(i.bucket), []byte(incrementKey), []byte(i.updateKey))
			db.Close()
		}()
		for {
			select {
			case signal := <-i.signalCh:
				i.updateKey = signal.IncrementKey
				if signal.Done {
					i.Close()
				}
			case <-i.done:
				return
			}
		}
	}()

	return nil
}

func (i *Increment) SendSignal(signal *Signal) error {
	if i.singalChStatus {
		return ErrSignalChClosed
	}
	i.singalChStatus = signal.Done
	i.signalCh <- signal
	return nil
}

func (i *Increment) Close() {
	i.once.Do(func() {
		close(i.done)
	})
}
