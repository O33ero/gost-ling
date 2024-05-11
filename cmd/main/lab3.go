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

	generator := lab3.New(seed)

	fmt.Printf("Next generator value: %d\n", generator.Next())
	fmt.Printf("Next generator value: %d\n", generator.Next())
	fmt.Printf("Next generator value: %d\n", generator.Next())
	fmt.Printf("Next generator value: %d\n", generator.Next())
	fmt.Printf("Next generator value: %d\n", generator.Next())
	fmt.Printf("Next generator value: %d\n", generator.Next())
	fmt.Printf("Next generator value: %d\n", generator.Next())
	fmt.Printf("Next generator value: %d\n", generator.Next())
	fmt.Printf("Next generator value: %d\n", generator.Next())
	fmt.Printf("Next generator value: %d\n", generator.Next())
	fmt.Printf("Next generator value: %d\n", generator.Next())
	fmt.Printf("Next generator value: %d\n", generator.Next())

	start := time.Now().UnixMilli()
	lab3.Create1MbFile(generator) // ~ 280 msec
	fmt.Printf("Elapsed [1mb]: %d msec.\n", time.Now().UnixMilli()-start)

	start = time.Now().UnixMilli()
	lab3.Create100MbFile(generator) // ~ 28000 msec (28 sec)
	fmt.Printf("Elapsed [100mb]: %d msec.\n", time.Now().UnixMilli()-start)

	start = time.Now().UnixMilli()
	lab3.Create1000MbFile(generator) // ~ 280000 msec (300 sec)
	fmt.Printf("Elapsed [1000mb]: %d msec.\n", time.Now().UnixMilli()-start)

	start = time.Now().UnixMilli()
	lab3.Create1000Values(generator) // ~ 3 msec
	fmt.Printf("Elapsed [1000val]: %d msec.\n", time.Now().UnixMilli()-start)

	start = time.Now().UnixMilli()
	lab3.Create10000Values(generator) // ~ 22 msec
	fmt.Printf("Elapsed [10000val]: %d msec.\n", time.Now().UnixMilli()-start)

	lab3.CreateNValuesInBinaryFormatFile(generator, 1000)
}
