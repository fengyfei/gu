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
 *     Initial: 2017/09/21        Jia Chenhui
 */

package command

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type PingSt struct {
	SendPk   string
	RecvPk   string
	LossPk   string
	MinDelay string
	AvgDelay string
	MaxDelay string
}

func PingLinux(addr string) PingSt {
	var ps PingSt

	cmd := exec.Command("ping", "-w", "60", "-i", "3", "-c", "20", addr)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("StdoutPipe returned error: %s", err.Error())
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)

	ps.SendPk = "20"
	ps.RecvPk = "0"
	ps.LossPk = "100"
	ps.MinDelay = "0"
	ps.AvgDelay = "0"
	ps.MaxDelay = "0"

	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || err2 == io.EOF {
			break
		}

		if strings.Contains(line, "packets transmitted") {
			pkg := strings.Fields(line)
			ps.SendPk = pkg[0]
			ps.RecvPk = pkg[3]
			ps.LossPk = strings.Split(pkg[5], "%")[0]
		}

		if strings.Contains(line, "rtt min/avg/max/mdev") {
			rttList := strings.Fields(line)
			rtt := strings.Split(rttList[3], "/")
			ps.MinDelay = rtt[0]
			ps.AvgDelay = rtt[1]
			ps.MaxDelay = rtt[2]
		}
	}

	cmd.Wait()
	return ps
}

func PingMacOS(addr string) PingSt {
	var ps PingSt

	cmd := exec.Command("ping", "-i", "3", "-c", "20", addr)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("StdoutPipe returned error: %s", err.Error())
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)

	ps.SendPk = "20"
	ps.RecvPk = "0"
	ps.LossPk = "100"
	ps.MinDelay = "0"
	ps.AvgDelay = "0"
	ps.MaxDelay = "0"

	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || err2 == io.EOF {
			break
		}

		if strings.Contains(line, "packets transmitted") {
			pkg := strings.Fields(line)
			ps.SendPk = pkg[0]
			ps.RecvPk = pkg[3]
			ps.LossPk = strings.Split(pkg[6], "%")[0]
		}

		if strings.Contains(line, "round-trip min/avg/max/stddev") {
			rttList := strings.Fields(line)
			rtt := strings.Split(rttList[3], "/")
			ps.MinDelay = rtt[0]
			ps.AvgDelay = rtt[1]
			ps.MaxDelay = rtt[2]
		}
	}

	cmd.Wait()
	return ps
}
