package main

import (
	"encoding/binary"
	"fmt"
	"gost-ling/internal/lab3"
	"time"
)

func main() {
	seed := []byte{
		0xA1, 0xB2, 0xC3, 0xD4, 0xE5, 0xF6, 0x17, 0x28,
		0xA1, 0xB2, 0xC3, 0xD4, 0xE5, 0xF6, 0x17, 0x28,
		0xA1, 0xB2, 0xC3, 0xD4, 0xE5, 0xF6, 0x17, 0x28,
		0xA1, 0xB2, 0xC3, 0xD4, 0xE5, 0xF6, 0x17, 0x28,
	}
	binary.BigEndian.PutUint64(seed, uint64(time.Now().UnixMilli()))

	random := lab3.New(seed)

	fmt.Printf("Next random value: %d\n", random.Next())
	fmt.Printf("Next random value: %d\n", random.Next())
	fmt.Printf("Next random value: %d\n", random.Next())
	fmt.Printf("Next random value: %d\n", random.Next())
	fmt.Printf("Next random value: %d\n", random.Next())
	fmt.Printf("Next random value: %d\n", random.Next())
	fmt.Printf("Next random value: %d\n", random.Next())
	fmt.Printf("Next random value: %d\n", random.Next())
	fmt.Printf("Next random value: %d\n", random.Next())
	fmt.Printf("Next random value: %d\n", random.Next())
	fmt.Printf("Next random value: %d\n", random.Next())
	fmt.Printf("Next random value: %d\n", random.Next())
}
