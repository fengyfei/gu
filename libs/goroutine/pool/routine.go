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
 *     Initial: 2017/12/16        Feng Yifei
 */

package pool

import (
	"context"
	"runtime/debug"
	"sync"

	"github.com/fengyfei/gu/libs/logger"
)

// routine is a goroutine with a stop channel.
type routine struct {
	goroutine func(chan bool)
	stop      chan bool
}

// Pool is a pool of go routines
type Pool struct {
	routines   []routine
	waitGroup  sync.WaitGroup
	lock       sync.Mutex
	baseCtx    context.Context
	baseCancel context.CancelFunc
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewPool creates a Pool
func NewPool(parentCtx context.Context) *Pool {
	baseCtx, baseCancel := context.WithCancel(parentCtx)
	ctx, cancel := context.WithCancel(baseCtx)
	return &Pool{
		baseCtx:    baseCtx,
		baseCancel: baseCancel,
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Ctx returns main context
func (p *Pool) Ctx() context.Context {
	return p.baseCtx
}

// Go starts a recoverable goroutine, and can be stopped with stop chan
func (p *Pool) Go(goroutine func(stop chan bool)) {
	p.lock.Lock()

	newRoutine := routine{
		goroutine: goroutine,
		stop:      make(chan bool, 1),
	}
	p.routines = append(p.routines, newRoutine)
	p.waitGroup.Add(1)

	Go(func() {
		goroutine(newRoutine.stop)
		p.waitGroup.Done()
	})
	p.lock.Unlock()
}

// Start starts all stopped routines
func (p *Pool) Start() {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.ctx, p.cancel = context.WithCancel(p.baseCtx)
	for i := range p.routines {
		p.waitGroup.Add(1)
		p.routines[i].stop = make(chan bool, 1)
		Go(func() {
			p.routines[i].goroutine(p.routines[i].stop)
			p.waitGroup.Done()
		})
	}
}

// Stop stops all started routines, waiting for their termination
func (p *Pool) Stop() {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.cancel()
	for _, routine := range p.routines {
		routine.stop <- true
	}
	p.waitGroup.Wait()
	for _, routine := range p.routines {
		close(routine.stop)
	}
}

// Cleanup releases resources used by the pool, and should be called when the pool will no longer be used
func (p *Pool) Cleanup() {
	p.Stop()
	p.lock.Lock()
	defer p.lock.Unlock()
	p.baseCancel()
}

// Go starts a recoverable goroutine
func Go(goroutine func()) {
	GoWithRecover(goroutine, defaultRecoverGoroutine)
}

// GoWithRecover starts a recoverable goroutine using given customRecover() function
func GoWithRecover(goroutine func(), customRecover func(err interface{})) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				customRecover(err)
			}
		}()
		goroutine()
	}()
}

func defaultRecoverGoroutine(err interface{}) {
	logger.Error(err.(error))
	debug.PrintStack()
}
