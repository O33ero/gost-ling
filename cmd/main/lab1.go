package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"gost-ling/internal/control"
	"gost-ling/internal/lab1"
	"os"
	"sync"
	"time"
)

func main() {
	control.NewAccessControl()
	ec := control.NewExecuteControl()
	defer ec.Wait()

	iv := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x10, 0x21, 0x32, 0x43, 0x54, 0x65, 0x76, 0x87,
	}
	key := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27,
		0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
	}
	plaintext := []byte{
		// block 0
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,

		// block 1
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,

		// block 2
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27,
		0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
	}
	fmt.Printf("Plain: %s\n", hex.EncodeToString(plaintext))

	cipher := lab1.NewOfb(key)
	defer cipher.Close()
	//// Encrypt
	//encryptBlocks := cipher.Encrypt(plaintext, iv)
	//for i, block := range encryptBlocks {
	//	fmt.Printf("Encrypt [%d]: %s\n", i, hex.EncodeToString(block[:]))
	//}
	//
	//// Decrypt
	//var ciphertext []byte
	//for _, block := range encryptBlocks {
	//	ciphertext = append(ciphertext, block[:]...)
	//}
	//
	//plainBlock := cipher.Decrypt(ciphertext, iv)
	//for i, block := range plainBlock {
	//	fmt.Printf("Decrypt [%d]: %s\n", i, hex.EncodeToString(block[:]))
	//}

	b, err := os.ReadFile("xoroshiro_1000mb.bin")
	if err != nil {
		panic("failed to read file: " + err.Error())
	}

	var wg sync.WaitGroup
	start := time.Now().UnixMilli()
	for i := 0; i < len(b); i += 16 {
		wg.Add(1)
		go func(part int) {
			defer wg.Done()
			encryptBlocks := cipher.Encrypt(b[part:part+16], iv)
			var ciphertext []byte
			for _, block := range encryptBlocks {
				ciphertext = append(ciphertext, block[:]...)
			}

			decryptBlocks := cipher.Decrypt(ciphertext[:], iv)
			var decrypt []byte
			for _, block := range decryptBlocks {
				decrypt = append(decrypt, block[:]...)
			}

			if !bytes.Equal(b[part:part+16], decrypt[:]) {
				panic("incorrect decrypt")
			}
		}(i)
	}
	wg.Wait()
	fmt.Printf("Complete in %d msec.\n", time.Now().UnixMilli()-start)
}
