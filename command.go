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
		fmt.Print("error!!!")
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}

type StdMsg struct {
	stdout, stderr *bytes.Buffer
}

func main() {
	machineName := "infra5"
	/// go get github.com/coreos/etcd//etcdctl
	cmd := exec.Command("etcdctl", "-C", "http://172.17.42.1:4002", "member", "add", machineName, "http://172.17.42.1:4005") // "/dev/random"
	//	cmd := exec.Command("ls", "/")

	stdmsg := StdMsg{}
	stdmsg.stderr = &bytes.Buffer{}
	stdmsg.stdout = &bytes.Buffer{}
	cmd.Stdout = stdmsg.stderr
	cmd.Stderr = stdmsg.stdout
	// Start command asynchronously
	err := cmd.Start()

	// Print command error when command is never exec
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

	outmsg := string(stdmsg.stdout.Bytes())
	if strings.Contains(outmsg, "Added member named "+machineName) {
		printOutput(
			stdmsg.stdout.Bytes(),
		)

	} else {

		printOutput(stdmsg.stderr.Bytes())

	}

}
