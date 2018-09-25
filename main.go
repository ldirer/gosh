package main

import (
	"io/ioutil"
	"log"
	"fmt"
	"math"
	"os/exec"
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

/**
1. Read command.
2. Parse it.
3. Run it.

 */
func gosh_loop() {
	//currentDirectory := "./"
	fmt.Println("IN LOOP")
	arg := ""
	fmt.Println("Gimme a command")
	_, err := fmt.Scanf("%s", &arg)

	if err != nil {
		log.Fatalln("Unexpected error getting command from stdin.")
	}



	parsedArgs, parseError := parseCommand(arg)
	if parseError != nil {
		panic(parseError)
	}

	runErr := goshRun(parsedArgs)
	if runErr != nil {
		log.Printf("Error executing command %s", parsedArgs)
		log.Fatalln(runErr)
	}

}
func goshRun(argv []string) error {
	procAttr := syscall.ProcAttr{Dir:"./"}
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


	//if command == "ls" {
	//	fileInfos, err := ioutil.ReadDir("./")
	//	if err == nil {
	//		for _, fileInfo := range fileInfos {
	//			dirSuffix := ""
	//			if fileInfo.IsDir() {
	//				dirSuffix = "/"
	//			}
	//			fmt.Printf("%-10s %-15s\n", formatSize(fileInfo.Size()), fileInfo.Name() + dirSuffix)
	//		}
	//	} else {
	//		return err
	//	}
	//}
	//if command == "cd"  {
	//
	//}
	//return fmt.Errorf("unknown command %s", command)
}
func parseCommand(s string) ([]string, error) {
	// First argument is the name of the program
	argv := []string{s}
	return argv, nil
}

func formatSize(nBytes int64) string {
	//units := []string{"KB", "MB", "GB"}

	kiloBytes := float64(nBytes) / math.Pow(2, 8)
	megaBytes := float64(nBytes) / math.Pow(2, 16)
	gigaBytes := float64(nBytes) / math.Pow(2, 24)

	if gigaBytes >= 1 {
		return fmt.Sprintf("%.2fGB", gigaBytes)
	}
	if megaBytes >= 1 {
		return fmt.Sprintf("%.2fMB", megaBytes)
	}
	if kiloBytes >= 1 {
		return fmt.Sprintf("%.2fKB", kiloBytes)
	}
	return fmt.Sprintf("%d", nBytes)
}
