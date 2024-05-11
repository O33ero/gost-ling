package lab3

import (
	"encoding/binary"
	"fmt"
	"os"
)

func Create1MbFile(rand *XoroShiro256StarStar) {
	iters := 131072 // 8 byte * 131072 = 1mb
	createFile(rand, iters, "xoroshiro_1mb.bin")
}

func Create100MbFile(rand *XoroShiro256StarStar) {
	iters := 13107200 // 8 byte * 13107200 = 100mb
	createFile(rand, iters, "xoroshiro_100mb.bin")
}

func Create1000MbFile(rand *XoroShiro256StarStar) {
	iters := 131072000 // 8 byte * 131072000 = 1000mb
	createFile(rand, iters, "xoroshiro_1000mb.bin")
}

func Create1000Values(rand *XoroShiro256StarStar) {
	iters := 1000 // 10^3
	createFile(rand, iters, "xoroshiro_1000values.bin")
}

func Create10000Values(rand *XoroShiro256StarStar) {
	iters := 10000 // 10^4
	createFile(rand, iters, "xoroshiro_10000values.bin")
}

func CreateNValuesInBinaryFormatFile(rand *XoroShiro256StarStar, n int) {
	createBinaryFormatFile(rand, n, fmt.Sprintf("xoroshiro_%dbin.txt", n))
}

func createBinaryFormatFile(rand *XoroShiro256StarStar, iters int, filename string) {
	err := os.Truncate(filename, 0)
	if err != nil {
		panic("failed to clear existed file")
	}
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic("failed to create file: " + err.Error())
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	for i := 0; i < iters; i++ {
		var nextBatch [8]byte
		value := rand.Next()
		binary.LittleEndian.PutUint64(nextBatch[:], value)

		for _, bt := range nextBatch {
			b := fmt.Sprintf("%08b", bt)
			_, err := file.WriteString(b)
			if err != nil {
				panic("failed to append to file: " + err.Error())
			}
		}
	}
}

func createFile(rand *XoroShiro256StarStar, iters int, filename string) {
	err := os.Truncate(filename, 0)
	if err != nil {
		panic("failed to clear existed file")
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic("failed to create file: " + err.Error())
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	for i := 0; i < iters; i++ {
		var nextBatch [8]byte
		value := rand.Next()
		binary.LittleEndian.PutUint64(nextBatch[:], value)

		_, err := file.Write(nextBatch[:])
		if err != nil {
			panic("failed to append to file: " + err.Error())
		}
	}
}
