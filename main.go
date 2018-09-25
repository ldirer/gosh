package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func main() {
	// Loading configuration files
	filename := "./gosh.rc"

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		// log is a nifty little utility which can also output
		// in this case, a fatal log will halt the program
		log.Fatalln("Error reading file", filename)
	}

	// content returned as []byte not string
	// so must cast to string first and then can display
	fmt.Println(string(content))

	// Main shell loop
	gosh_loop()

	// terminate/cleanup
}


type GoshState struct {
	currentDirectory string
}

var state = GoshState{currentDirectory: "./"}

/**
1. Read command.
2. Parse it.
3. Run it.
 */
func gosh_loop() {
	fmt.Println("IN LOOP")
	fmt.Println("Gimme a command")
	for {
		//_, err := fmt.Scanf("%s", &arg)
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')

		if err != nil {
			log.Fatalln("Unexpected error getting command from stdin.")
		}

		parsedArgs := parseCommand(input)

		runErr := goshRun(parsedArgs)
		if runErr != nil {
			log.Printf("Error executing command %s", parsedArgs)
			log.Fatalln(runErr)
		}

	}

}
func goshRun(argv []string) error {
	procAttr := syscall.ProcAttr{
		Dir: state.currentDirectory,
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
	}
	// ForkExec expects the full path to the executable.
	binary, lookErr := exec.LookPath(argv[0])
	if lookErr != nil {
		panic(lookErr)
	}
	childPid, err := syscall.ForkExec(binary, argv, &procAttr)

	if err != nil {
		return err
	}

	fmt.Printf("Child process pid: %d\n", childPid)
	return err
}

func parseCommand(s string) []string {
	// First argument is the name of the program
	argv := []string{s}
	argv = strings.Fields(s)
	fmt.Println(argv)
	return argv
}

