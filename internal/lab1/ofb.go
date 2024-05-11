package lab1

import (
	"fmt"
	"runtime"
)

type Ofb struct {
	cipher *Cipher
}

func (c *Ofb) Close() {
	c.cipher.Close()
	runtime.GC()
	fmt.Printf("Clear mem [Ofb]: %p\n", &c)
}

func NewOfb(key []byte) *Ofb {
	return &Ofb{
		cipher: NewCipher(key),
	}
}

func (c *Ofb) Encrypt(plaintext []byte, iv []byte) [][BlockSize]byte {
	var result [][BlockSize]byte

	plainBlock := iv
	for i := 0; i < len(plaintext); i += BlockSize {
		iterBlock := c.cipher.Encrypt(plainBlock)

		var encryptBlock [BlockSize]byte
		xor(encryptBlock[:], plaintext[i:i+16], iterBlock[:])
		result = append(result, encryptBlock)

		plainBlock = iterBlock[:]
	}

	return result
}

func (c *Ofb) Decrypt(ciphertext []byte, iv []byte) [][BlockSize]byte {
	var result [][BlockSize]byte

	cipherBlock := iv
	for i := 0; i < len(ciphertext); i += BlockSize {
		iterBlock := c.cipher.Encrypt(cipherBlock)

		var plainBlock [BlockSize]byte
		xor(plainBlock[:], ciphertext[i:i+16], iterBlock[:])
		result = append(result, plainBlock)

		cipherBlock = iterBlock[:]
	}

	return result
}
