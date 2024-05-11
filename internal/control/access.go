package control

import (
	"bufio"
	"encoding/hex"
	"flag"
	"gost-ling/internal/lab2"
	"os"
)

type AccessControl struct {
}

func NewAccessControl() *AccessControl {
	pArg := flag.String("p", "", "password")
	flag.Parse()

	if pArg == nil || len(*pArg) == 0 {
		panic("arg -p (password) is not set, but required")
	}
	h := lab2.NewHash()
	h.Write([]byte(*pArg))
	passwd := hex.EncodeToString(h.Sum(nil))

	file, err := os.Open("./access/access.txt")
	if err != nil {
		panic("failed to open file ./access/access.txt")
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if passwd == scanner.Text() {
			return &AccessControl{}
		}
	}
	panic("permission denied. Password incorrect")
}
