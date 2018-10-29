package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func e(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	command()
	output()
}

func output() {
	count := exec.Command("../count/count")
	stdout, _ := count.StdoutPipe()

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Printf("(stdout)%s\n", scanner.Text())
		}
	}()

	err := count.Run()
	e(err)
}

func command() {
	if len(os.Args) == 1 {
		return
	}

	cmd := exec.Command(os.Args[1], os.Args[2])
	err := cmd.Run()
	e(err)

	state := cmd.ProcessState
	fmt.Printf("%s\n", state.String())
	fmt.Printf(" Pid: %d\n", state.Pid())
	fmt.Printf(" System: %v\n", state.SystemTime())
	fmt.Printf(" User: %v\n", state.UserTime())

}
