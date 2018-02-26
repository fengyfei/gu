/*
 * Revision History:
 *     Initial: 2018/02/26        Tong Yuehong
 */

package main

import (
	"fmt"
	"github.com/fengyfei/gu/libs/hash/cityhash"
)

func main() {
	str := []byte("1234567abcde")

	ch64 := cityhash.CityHash64(str, uint32(len(str)))
	fmt.Println("CityHash64:", ch64)

	ch128 := cityhash.CityHash128(str, uint32(len(str)))
	fmt.Println("CityHash128:", ch128)

	chs64 := cityhash.CityHash64WithSeed(str, uint32(len(str)), uint64(1111111))
	fmt.Println("CityHash64WithSeed:", chs64)

	chs128 := cityhash.CityHash128WithSeed(str, uint32(len(str)), cityhash.Uint128{111111, 0xc3a5c85c97cb3127})
	fmt.Println("CityHash128WithSeed:", chs128)
}
