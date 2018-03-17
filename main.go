package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const TOKDELIM = " \t\r\n\a"
const ERRFORMAT = "sesh: %s\n"

var CWD string

func main() {
	/* Importing config */

	sesh_loop()

	os.Exit(0)
}

func sesh_loop() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf(ERRFORMAT, err.Error())
	}
	wdSlice := strings.Split(wd, "/")
	CWD = wdSlice[len(wdSlice)-1]

	reader := bufio.NewReader(os.Stdin)
	status := true

	for status {
		fmt.Printf("sesh 🔥  %s -> ", CWD)
		line, err := reader.ReadString('\n')
		line = line[:len(line)-1]
		if err != nil {
			log.Fatalf("sesh: %s\n", err.Error())
		}
		args := splitIntoTokens(line)
		status = execute(args)
	}
}

func splitIntoTokens(line string) []string {
	/* Need to add proper method to support all kinds of tokenising */
	return strings.Fields(line)
}

func launch(args []string) {
	// Spawning and executing a process
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = nil // making sure the command uses the current process' environment
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Printf(ERRFORMAT, err.Error())
	} else {
		fmt.Printf(out.String())
	}
}

func execute(args []string) bool {
	if len(args) == 0 {
		return true
	}
	for k, v := range builtins {
		if args[0] == k {
			return v(args[1:])
		}
	}
	launch(args)
	return true
}
