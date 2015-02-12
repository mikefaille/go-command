package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}

func main() {

	/// go get github.com/coreos/etcd//etcdctl
	cmd := exec.Command("etcdctl", "-C", "http://172.17.42.1:4002", "member", "add", "infra5", "http://172.17.42.1:4005") // "/dev/random"
	//	cmd := exec.Command("ls", "/")
	randomBytes := &bytes.Buffer{}
	cmd.Stdout = randomBytes

	// Start command asynchronously
	err := cmd.Start()
	printError(err)

	// Create a ticker that outputs elapsed time
	ticker := time.NewTicker(time.Second)
	go func(ticker *time.Ticker) {
		now := time.Now()
		for _ = range ticker.C {
			printOutput(
				[]byte(fmt.Sprintf("%s", time.Since(now))),
			)
		}
	}(ticker)

	// Create a timer that will kill the process
	timer := time.NewTimer(time.Second * 4)
	go func(timer *time.Timer, ticker *time.Ticker, cmd *exec.Cmd) {
		for _ = range timer.C {
			err := cmd.Process.Signal(os.Kill)
			printError(err)
			ticker.Stop()
		}
	}(timer, ticker, cmd)

	// Only proceed once the process has finished
	cmd.Wait()
	//	fmt.Println(string(randomBytes.Bytes()))
	printOutput(
		randomBytes.Bytes(),
	)
}
