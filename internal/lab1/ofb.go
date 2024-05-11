package lab1

import (
	"fmt"
	"runtime"
)

type Ofb struct {
	cipher        *Cipher
	initialVector []byte
}

func (c *Ofb) Close() {
	for i := 0; i < len(c.initialVector); i++ {
		c.initialVector[i] = 0x00
	}
	runtime.GC()
	fmt.Printf("Clear mem [Ofb]: %p\n", &c)
}

func NewOfb(iv []byte, key []byte) *Ofb {
	if len(iv) != BlockSize {
		panic("initial vector should be 16 bytes")
	}

	return &Ofb{
		cipher:        NewCipher(key),
		initialVector: iv,
	}
}

func (c *Ofb) Encrypt(plaintext []byte) [][BlockSize]byte {
	var result [][BlockSize]byte

	plainBlock := c.initialVector
	for i := 0; i < len(plaintext); i += BlockSize {
		iterBlock := c.cipher.Encrypt(plainBlock)

		var encryptBlock [BlockSize]byte
		xor(encryptBlock[:], plaintext[i:i+16], iterBlock[:])
		result = append(result, encryptBlock)

		plainBlock = iterBlock[:]
	}

	return result
}

func (c *Ofb) Decrypt(ciphertext []byte) [][BlockSize]byte {
	var result [][BlockSize]byte

	cipherBlock := c.initialVector
	for i := 0; i < len(ciphertext); i += BlockSize {
		iterBlock := c.cipher.Encrypt(cipherBlock)

		var plainBlock [BlockSize]byte
		xor(plainBlock[:], ciphertext[i:i+16], iterBlock[:])
		result = append(result, plainBlock)

		cipherBlock = iterBlock[:]
	}

	return result
}
