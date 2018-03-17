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
	status := 1

	for status != 0 {
		symbol := "\u2713"
		if status == 2 {
			symbol = "\u2715"
		}
		fmt.Printf("sesh 🔥  %s %s ", CWD, symbol)
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
	/* Need to add proper method to support all kinds of tokenising
	For now, returning tokens split with any whitespace as a delimiter */
	return strings.Fields(line)
}

func launch(args []string) int {
	// Spawning and executing a process
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = nil // making sure the command uses the current process' environment
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Printf(ERRFORMAT, err.Error())
		return 2
	}
	fmt.Printf(out.String())
	return 1
}

func execute(args []string) int {
	if len(args) == 0 {
		return 1
	}
	for k, v := range builtins {
		if args[0] == k {
			return v(args[1:])
		}
	}
	return launch(args)
}
