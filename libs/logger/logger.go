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
 *     Initial: 2017/10/28        Feng Yifei
 */

package logger

import (
	"os"

	"github.com/op/go-logging"
)

const guModuleID = "gu"

var (
	guLogger         *logging.Logger
	backend          *logging.LogBackend
	format           logging.Formatter
	backendFormatter logging.Backend
)

func init() {
	guLogger = logging.MustGetLogger(guModuleID)
	backend = logging.NewLogBackend(os.Stderr, "", 0)
	format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{id:03x}%{color:reset} %{message}`,
	)
	backendFormatter = logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
}

// Debug prints debug messages.
func Debug(v ...interface{}) {
	if len(v) > 1 {
		guLogger.Debug(v[0].(string), v[1:]...)
	} else {
		guLogger.Debug("Debug: %s", v...)
	}
}

// Info prints normal messages.
func Info(v ...interface{}) {
	if len(v) > 1 {
		guLogger.Info(v[0].(string), v[1:]...)
	} else {
		guLogger.Info("Info: %s", v...)
	}
}

// Warn prints warning messages.
func Warn(v ...interface{}) {
	if len(v) > 1 {
		guLogger.Warning(v[0].(string), v[1:]...)
	} else {
		guLogger.Warning("Warning: %s", v...)
	}
}

// Error prints error.
func Error(v ...interface{}) {
	if len(v) > 1 {
		guLogger.Error(v[0].(string), v[1:]...)
	} else {
		guLogger.Error("Error: %s", v...)
	}
}
