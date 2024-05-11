package lab2

import (
	"encoding/binary"
	"os"
)

func Create10in4Keys(key []byte) {
	iters := 10000
	kdf := NewKDF(key)
	createFile(kdf, iters, "kdf_10000.bin")
}

func Create10in5Keys(key []byte) {
	iters := 100000
	kdf := NewKDF(key)
	createFile(kdf, iters, "kdf_100000.bin")
}

func Create10in6Keys(key []byte) {
	iters := 1000000
	kdf := NewKDF(key)
	createFile(kdf, iters, "kdf_1000000.bin")
}

func createFile(kdf *KDF, iters int, filename string) {
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

	label := []byte{
		0x00, 0x01, 0x02, 0x03,
	}
	for i := 0; i < iters; i++ {
		var seq [8]byte
		binary.LittleEndian.PutUint64(seq[:], uint64(i))

		nextKey := kdf.Derive(label, seq[:], 1)

		_, err := file.Write(nextKey[:])
		if err != nil {
			panic("failed to append to file: " + err.Error())
		}
	}
}
