/*
 * Revision History:
 *     Initial: 2018/02/26        Shi Ruitao
 */
package main

import (
	"fmt"
	"time"

	farm "github.com/dgryski/go-farm"
)

var res32 uint32
var res64 uint64
var res64lo, res64hi uint64

var buf = []byte("RMVx)@MLxH9M.WeGW-ktWwR3Cy1")

func main() {
	hash()
	fingerprint()
	hashWithSeed()
}

func hash() {
	// Hash32 hashes a byte slice and returns a uint32 hash value
	res32 = farm.Hash32(buf)
	fmt.Println("Hash32: ", res32)

	res64 = farm.Hash64(buf)
	fmt.Println("Hash64: ", res64)

	res64lo, res64hi = farm.Hash128(buf)
	fmt.Println("Hash128: ", res64lo, res64hi)
}

func fingerprint() {
	// Fingerprint32 is a 32-bit fingerprint function for byte-slices
	res32 = farm.Fingerprint32(buf)
	fmt.Println("Fingerprint32: ", res32)
}

func hashWithSeed() {
	// Hash32WithSeed hashes a byte slice and a uint32 seed and returns a uint32 hash value
	res64 = farm.Hash64WithSeed(buf, uint64(time.Now().UnixNano()))
	fmt.Println("Hash32WithSeed: ", res64)
}
