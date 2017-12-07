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
 *     Initial: 2017/12/07        Feng Yifei
 */

package godash

import (
	"reflect"
)

// Each iterates a map, slice or map and do iterator(key, value)
func Each(collection, iterator interface{}) {
	each(collection, iterator)
}

func each(collection, handler interface{}) {
	size, iterator := keyValueIterator(collection)

	if iterator == nil || size == 0 {
		return
	}

	handlerValue := reflect.ValueOf(handler)

	for i := 0; i < size; i++ {
		k, v := iterator(i)

		handlerValue.Call([]reflect.Value{k, v})
	}
}

func keyValueIterator(collection interface{}) (int, func(int) (reflect.Value, reflect.Value)) {
	if collection == nil {
		return 0, nil
	}

	value := reflect.ValueOf(collection)
	switch value.Kind() {
	case reflect.Array:
	case reflect.Slice:
		return value.Len(), func(i int) (reflect.Value, reflect.Value) {
			return reflect.ValueOf(i), value.Index(i)
		}
	case reflect.Map:
		keys := value.MapKeys()

		return len(keys), func(i int) (reflect.Value, reflect.Value) {
			return keys[i], value.MapIndex(keys[i])
		}
	}

	return 0, nil
}
