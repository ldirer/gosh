package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func main() {
	// Loading configuration files. Ended up not really using it.
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
	goshLoop()

	// terminate/cleanup
}


type GoshState struct {
	currentDirectory string
}

var state = GoshState{currentDirectory: "./"}


type CommandRegister struct {
	Commands map[string] func(args ...string) error
}

var commandRegister = CommandRegister{Commands: map[string] func(args ...string) error{
	"cd": cd,
}}


const HOME = "/home/laurent"

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func NewError(s string) error {
	return &errorString{s: s}
}



// `args` should not contain the name of the command in this case.
func cd (args ...string) error {
	// All but first arg are ignored
	destination := HOME
	if len(args) > 0 {
		destination = args[0]
	}

	destination, absErr := filepath.Abs(destination)
	if absErr != nil {
		return absErr
	}


	dst, err := os.Stat(destination)
	if os.IsNotExist(err) {
		return err
	}

	if !dst.IsDir() {
		return NewError(fmt.Sprintf("Target is not a directory: %s", dst.Name()))
	}

	// chdir changes current directory... But I dont know how to retrieve that!
	// We need to chdir so that absolute path computation works.
	os.Chdir(destination)
	state.currentDirectory = destination

	fmt.Println(state.currentDirectory)
	return nil
}

func getNil() error {
	return nil
}

const beer = "\U0001f37a"

func prompt() string {
	absPath, _ := filepath.Abs(state.currentDirectory)
	return fmt.Sprintf("%s %s -> ", beer, absPath)
}

/**
1. Read command.
2. Parse it.
3. Run it.
 */
func goshLoop() {
	fmt.Println("IN LOOP")
	fmt.Println("Gimme a command")
	for {
		// this print is racing against the output of the previous command. Tis a bit weird.
		// We could wait on the forked process... but I'm not sure this is how a shell should work.
		// --> Though this is what all shells do I guess.
		fmt.Print(prompt())
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		// For debugging. Cant debug with stdin cuz debugger and program fight for it.
		//input := "cd /home/laurent"
		//err := getNil()

		if err != nil {
			log.Fatalln(err)
		}

		parsedArgs := parseCommand(input)
		if len(parsedArgs) == 0 {
			continue
		}

		runErr := goshRun(parsedArgs)
		if runErr != nil {
			log.Printf("Error executing command %s", parsedArgs)
			// no reason to crash here, we want our shell to outlive typos and mistakes.
			log.Printf(runErr.Error())
		}
	}
}

func goshRun(argv []string) error {
	procAttr := syscall.ProcAttr{
		Dir: state.currentDirectory,
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
	}

	handlerFunction, isBuiltin := commandRegister.Commands[argv[0]]
	if isBuiltin {
		return handlerFunction(argv[1:]...)
	}

	// ForkExec expects the full path to the executable.
	binary, lookErr := exec.LookPath(argv[0])
	if lookErr != nil {
		return lookErr
	}
	//fmt.Println("about to ForkExec")
	_, err := syscall.ForkExec(binary, argv, &procAttr)

	if err != nil {
		return err
	}

	//fmt.Printf("Child process pid: %d\n", _)
	return err
}

func parseCommand(s string) []string {
	// First argument is the name of the program
	argv := []string{s}
	argv = strings.Fields(s)
	return argv
}

