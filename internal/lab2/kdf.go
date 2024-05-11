package lab2

import (
	"crypto/hmac"
	"fmt"
	"hash"
	"runtime"
)

type KDF struct {
	h hash.Hash
}

func NewKDF(key []byte) *KDF {
	return &KDF{hmac.New(newHash, key)}
}

func (kdf *KDF) Close() {
	kdf.h.Reset()
	runtime.GC()
	fmt.Printf("Clear mem [KDF]: %p\n", &kdf)
}

func (kdf *KDF) Derive(label, seed []byte, r int) (res []byte) {
	if r < 0 || r > 4 {
		panic("R should be between 1 and 4 inclusive")
	}

	for i := 1; i <= r; i++ {
		kdf.h.Write([]byte{byte(i)})
		kdf.h.Write(label)
		kdf.h.Write([]byte{0x00})
		kdf.h.Write(seed)
		kdf.h.Write([]byte{0x01})
		kdf.h.Write([]byte{0x00})
	}

	res = kdf.h.Sum(nil)
	kdf.h.Reset()
	return res
}

func newHash() hash.Hash {
	return NewHash()
}
