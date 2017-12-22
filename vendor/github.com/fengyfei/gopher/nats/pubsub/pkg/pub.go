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
 *     Initial: 2017/10/22        Jia Chenhui
 */

package pkg

// PubObject represents the nats object to be sent.
type PubObject struct {
	subject string
	data    []byte
}

// NewPubObj create a PubObject.
func NewPubObj(subj string, data []byte) *PubObject {
	return &PubObject{
		subject: subj,
		data:    data,
	}
}

// Pub publish objects on the supplied connection.
// nc can be nil.
func (po *PubObject) Pub(nc *NatsConn) error {
	var (
		err error
	)

	if nc == nil {
		nc, err = NewConn("")
		if err != nil {
			return err
		}
	}

	defer nc.conn.Flush()

	err = nc.conn.Publish(po.subject, po.data)
	if err != nil {
		return err
	}

	return nil
}
