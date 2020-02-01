package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jpr98/memsim/mem"

	"github.com/jpr98/memsim/cpu"
)

var comp cpu.CPU

func main() {
	filename := flag.String("filename", "input.txt", "the name of the file to read instructions from")
	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil {
		errorMessage := fmt.Sprintf("File %s couldn't be opened. Verify it's existance.", *filename)
		panic(errorMessage)
	}

	scanner := bufio.NewScanner(file)

	comp = cpu.New(
		mem.New(1204, 16),
		mem.New(1204, 16),
	)

	for scanner.Scan() {
		err := parseCommand(scanner.Text())
		if err != nil {
			panic("received error")
		}
	}

	fmt.Println("Done!")
}

func parseCommand(cmdStr string) error {
	cmd := strings.Split(cmdStr, " ")

	switch cmd[0] {
	case "P":
		if len(cmd) != 3 {
			return cmdArgsError("P", 3)
		}
		handleProcess(cmd)
	case "A":
		if len(cmd) != 4 {
			return cmdArgsError("A", 4)
		}
		handleAccess(cmd)
	case "L":
		if len(cmd) != 2 {
			return cmdArgsError("L", 2)
		}
		handleClear(cmd)
	case "C":
		// TODO: this will be a comment, we need to join the words to form the sentence
		if len(cmd) < 2 {
			return cmdArgsError("C", 2)
		}
		handleComment(cmd)
	case "F":
		if len(cmd) > 1 {
			return cmdArgsError("F", 1)
		}
		handleFinalize()
	case "E":
		if len(cmd) > 1 {
			return cmdArgsError("E", 1)
		}
		handleEnd()
	default:
		fmt.Println("Invalid command")
	}

	return nil
}

func handleProcess(cmd []string) {
	size, _ := strconv.Atoi(cmd[1]) //FIXME: Handle error
	pid := cmd[2]

	fmt.Printf("Loading PID: %s size: %d\n", pid, size)
	err := comp.CreateProcess(pid, size)
	if err != nil {
		panic(err)
	}
}

func handleAccess(cmd []string) {
	address, _ := strconv.Atoi(cmd[1]) //FIXME: Handle error
	pid := cmd[2]
	modify, _ := strconv.ParseBool(cmd[3]) //FIXME: Handle error

	fmt.Printf("Accessing PID: %s address: %d modify: %t\n", pid, address, modify)
	add, err := comp.AccessProcess(pid, address)
	if err != nil {
		panic(err)
	}
	fmt.Printf("found at real address %d\n", add)
}

func handleClear(cmd []string) {
	pid := cmd[1]

	fmt.Printf("Clearing PID: %s\n", pid)
}

func handleComment(cmd []string) {
	comment := strings.Join(cmd[1:], " ")
	fmt.Printf("Comment %s\n", comment)
}

func handleFinalize() {
	fmt.Println("Finalized this sequence of instructions")
}

func handleEnd() {
	fmt.Println("End. Thanks for using the program!")
}
