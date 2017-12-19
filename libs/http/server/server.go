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
 *     Initial: 2017/12/18        Feng Yifei
 */

package server

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	errNoRouter = errors.New("Entrypoint requires a router")
)

// Entrypoint represents a http server.
type Entrypoint struct {
	configuration *Configuration
	tlsConfig     *TLSConfiguration
	server        *http.Server
	listener      net.Listener
	stop          chan bool
	signals       chan os.Signal
}

// NewEntrypoint creates a new Entrypoint.
func NewEntrypoint(conf *Configuration, tlsConf *TLSConfiguration) *Entrypoint {
	return &Entrypoint{
		configuration: conf,
		tlsConfig:     tlsConf,
		stop:          make(chan bool, 1),
		signals:       make(chan os.Signal, 1),
	}
}

// Prepare the entrypoint for serving requests.
func (ep *Entrypoint) prepare(router http.Handler) error {
	var (
		err       error
		tlsConfig *tls.Config
		listener  net.Listener
	)

	if tlsConfig, err = ep.createTLSConfig(); err != nil {
		return err
	}

	listener, err = net.Listen("tcp", ep.configuration.Address)
	if err != nil {
		return err
	}

	ep.listener = listener
	ep.server = &http.Server{
		Addr:      ep.configuration.Address,
		Handler:   router,
		TLSConfig: tlsConfig,
	}

	return nil
}

// Create the TLS Configuration for the http server.
func (ep *Entrypoint) createTLSConfig() (*tls.Config, error) {
	if ep.tlsConfig == nil {
		return nil, nil
	}

	return nil, nil
}

func (ep *Entrypoint) startServer() error {
	if ep.tlsConfig != nil {
		return ep.server.ServeTLS(ep.listener, "", "")
	}

	return ep.server.Serve(ep.listener)
}

// Start the entrypoint.
func (ep *Entrypoint) Start(router http.Handler) error {
	if router == nil {
		return errNoRouter
	}

	if err := ep.prepare(router); err != nil {
		return err
	}

	ep.configureSignals()

	go ep.listenSignals()
	go ep.startServer()

	return nil
}

// Wait until stop channel emits a value.
func (ep *Entrypoint) Wait() {
	<-ep.stop
}

// Stop the http server.
func (ep *Entrypoint) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	// graceful shutdown
	if err := ep.server.Shutdown(ctx); err != nil {
		ep.server.Close()
	}
	cancel()

	close(ep.stop)
}
