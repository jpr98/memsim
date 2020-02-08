package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jpr98/memsim/cpu"
	"github.com/jpr98/memsim/mem"
)

var comp cpu.CPU
var reader *bufio.Reader
var breaking *bool

func main() {
	filename := flag.String("filename", "input.txt", "the name of the file to read instructions from")
	debug := flag.Bool("debug", false, "debug mode, pass true to see memory allocation")
	policy := flag.String("policy", "FIFO", "sets the replacement policy to be used when swapping pages")
	breaking = flag.Bool("breaking", false, "*experimental* asks user to continue or quit execution if there is an error")
	server := flag.Bool("server", false, "if true a server with live monitoring will start")
	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil {
		errorMessage := fmt.Sprintf("File %s couldn't be opened. Verify it's existance.", *filename)
		fmt.Println(errorMessage)
		fmt.Println("Stopping execution")
		os.Exit(1)
		return
	}
	defer file.Close()

	MMU, err := mem.NewMMU(2048, 4096, 16, *policy)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Stopping execution")
		os.Exit(1)
	}
	comp = cpu.New(*MMU)

	if *breaking {
		fmt.Println("")
		fmt.Println("* You are using breaking, which is an experimental feature, if you decide to continue after an error anything could go wrong *")
		fmt.Println("")
	}

	c := make(chan bool)
	if *server {
		go startServer(c)
	}

	reader = bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		err := parseCommand(scanner.Text())
		if err != nil {
			fmt.Printf("Error parsing instruction: %s\n", err.Error())
			if cont := askBreak(); cont {
				continue
			}
		}
		comp.Print(*debug, os.Stdout)
	}

	fmt.Println("Done!")
	if *server {
		<-c
	}
}

func startServer(c chan bool) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		comp.Print(true, w)
	})

	http.ListenAndServe(":8000", nil)
}

func parseCommand(cmdStr string) error {
	cmd := strings.Fields(cmdStr)
	start := time.Now()
	switch strings.ToUpper(cmd[0]) {
	case "P":
		if len(cmd) < 3 {
			return cmdArgsError("P", 3)
		}
		handleCreateProcess(cmd)
	case "A":
		if len(cmd) < 4 {
			return cmdArgsError("A", 4)
		}
		handleAccess(cmd)
	case "L":
		if len(cmd) < 2 {
			return cmdArgsError("L", 2)
		}
		handleClear(cmd)
	case "C":
		if len(cmd) < 1 {
			return cmdArgsError("C", 2)
		}
		handleComment(cmd)
	case "F":
		if len(cmd) > 1 {
			if confirmIntention("Did you mean F (finalize)?") {
				handleFinalize()
			} else {
				return cmdArgsError("F", 1)
			}
		} else {
			handleFinalize()
		}
	case "E":
		if len(cmd) > 1 {
			if confirmIntention("Did you mean E (end)?") {
				handleEnd()
			} else {
				return cmdArgsError("E", 1)
			}
		} else {
			handleEnd()
		}
	default:
		fmt.Println("Invalid command")
	}
	fmt.Println(time.Since(start))
	return nil
}

// askBreak checks if user wants to continue even though an error ocurred, something can go wrong if they decide to continue
func askBreak() bool {
	if *breaking {
		fmt.Print("Do you wish to continue with next instruction? [y/n] ")
		resp, _ := reader.ReadString('\n')
		resp = strings.ToLower(resp)
		resp = strings.TrimSpace(resp)
		if resp == "n" {
			fmt.Println("Stopping execution")
			os.Exit(1)
			return false
		}
		return true
	}
	fmt.Println("Stopping execution")
	os.Exit(1)
	return false
}

func confirmIntention(message string) bool {
	fmt.Printf("%s [y/n] ", message)
	resp, _ := reader.ReadString('\n')
	resp = strings.ToLower(resp)
	resp = strings.TrimSpace(resp)
	if resp == "n" {
		return false
	}
	return true
}

func handleCreateProcess(cmd []string) {
	size, err := strconv.Atoi(cmd[1])
	if err != nil {
		fmt.Printf("Argument %s needs to be a number\n", cmd[1])
		if cont := askBreak(); cont {
			return
		}
	}
	pid := cmd[2]

	fmt.Printf("Loading PID: %s size: %d\n", pid, size)
	err = comp.CreateProcess(pid, size)
	if err != nil {
		fmt.Println(err.Error())
		if cont := askBreak(); cont {
			return
		}
	}
}

func handleAccess(cmd []string) {
	address, err := strconv.Atoi(cmd[1])
	if err != nil {
		fmt.Printf("Argument %s needs to be a number\n", cmd[1])
		if cont := askBreak(); cont {
			return
		}
	}
	pid := cmd[2]
	modify, err := strconv.ParseBool(cmd[3])
	if err != nil {
		fmt.Printf("Argument %s needs to be a boolean\n", cmd[3])
		if cont := askBreak(); cont {
			return
		}
	}

	fmt.Printf("Accessing PID: %s address: %d modify: %t\n", pid, address, modify)
	add, err := comp.AccessProcess(pid, address)
	if err != nil {
		fmt.Println(err.Error())
		if cont := askBreak(); cont {
			return
		}
	}
	fmt.Printf("found at real address %d\n", add)
}

func handleClear(cmd []string) {
	pid := cmd[1]

	fmt.Printf("Clearing PID: %s\n", pid)
	err := comp.DeleteProcess(pid)
	if err != nil {
		fmt.Println(err.Error())
		if cont := askBreak(); cont {
			return
		}
	}
}

func handleComment(cmd []string) {
	if len(cmd) == 0 {
		fmt.Println("")
		return
	}
	comment := strings.Join(cmd[1:], " ")
	fmt.Printf("Comment: %s\n", comment)
}

func handleFinalize() {
	fmt.Println("Finalized this sequence of instructions")
	fmt.Println("Reseting system")
	fmt.Println("----------------------------------------------------")
}

func handleEnd() {
	fmt.Println("End. Thanks for using the program!")
}
