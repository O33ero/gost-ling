package control

import (
	"bytes"
	"fmt"
	"gost-ling/internal/lab2"
	"os"
	"sync"
	"time"
)

type ExecuteControl struct {
	ticker *time.Ticker
	target []byte
	hash   *lab2.Hash
	wg     *sync.WaitGroup
}

func NewExecuteControl() *ExecuteControl {
	hsh := lab2.NewHash()
	executable := os.Args[0]
	file, err := os.ReadFile(executable)
	if err != nil {
		panic("failed to get current executable file")
	}
	hsh.Write(file)
	target := hsh.Sum(nil)
	hsh.Reset()

	ticker := time.NewTicker(5 * time.Second)
	ec := &ExecuteControl{
		hash:   hsh,
		target: target,
		ticker: ticker,
		wg:     &sync.WaitGroup{},
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				ec.run()
			}
		}
	}()

	ec.run()
	return ec
}

func (ec *ExecuteControl) run() {
	ec.wg.Add(1)
	fmt.Printf("Run execution control...\n")
	executablePath := os.Args[0]
	file, err := os.ReadFile(executablePath)
	if err != nil {
		panic("failed to get current executable file")
	}
	ec.hash.Write(file)
	target := ec.hash.Sum(nil)
	ec.hash.Reset()

	if !bytes.Equal(target, ec.target) {
		panic("execution control failed")
	} else {
		fmt.Printf("Execution control complete\n")
		ec.wg.Done()
	}

}

func (ec *ExecuteControl) Wait() {
	ec.wg.Wait()
}
