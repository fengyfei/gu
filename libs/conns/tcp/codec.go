/*
 * MIT License
 *
 * Copyright (c) 2018 SmartestEE Co., Ltd.
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
 *     Initial: 2018/01/08        Feng Yifei
 */

package tcp

import (
	"bytes"
	"encoding/binary"
)

func encode(message *Message) ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.BigEndian, message.len)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, message.id)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, message.payload)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func decode(payload []byte) (Message, error) {
	var message Message

	buf := bytes.NewReader(payload)

	err := binary.Read(buf, binary.BigEndian, &message.len)
	if err != nil {
		return nil, err
	}

	err = binary.Read(buf, binary.BigEndian, &message.id)
	if err != nil {
		return nil, err
	}

	err = binary.Read(buf, binary.BigEndian, &message.payload)
	if err != nil {
		return nil, err
	}

	return message, nil
}
