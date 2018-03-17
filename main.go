package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const TOKDELIM = " \t\r\n\a"
const ERRFORMAT = "sesh: %s\n"

var (
	CWD      string
	HISTSIZE = 25
	HISTFILE = ".sesh_history"
	HISTMEM  []string
)

func main() {
	/* Importing config */

	sesh_loop()

	exit()
}

func sesh_loop() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf(ERRFORMAT, err.Error())
	}
	wdSlice := strings.Split(wd, "/")
	CWD = wdSlice[len(wdSlice)-1]

	history := make([]string, 0)
	history = initHistory(history)
	status, historyCount := 1, 0
	reader := bufio.NewReader(os.Stdin)
	//go listenForKeyPress()

	for status != 0 {
		symbol := "\u2713"
		if status == 2 {
			symbol = "\u2715"
		}
		fmt.Printf("sesh 🔥  %s %s ", CWD, symbol)
		line, _ := reader.ReadString('\n')
		line = line[:len(line)-1]
		args := splitIntoTokens(line)
		status = execute(args)
		if status == 1 {
			/* Store line in history */
			if historyCount == HISTSIZE {
				history = history[1:]
			}
			history = append(history, line)
			historyCount++
		}
	}
	// Reversing the history slice
	last := len(history) - 1
	for i := 0; i < len(history)/2; i++ {
		history[i], history[last-i] = history[last-i], history[i]
	}
	HISTMEM = history
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

func initHistory(history []string) []string {
	if _, err := os.Stat(HISTFILE); err == nil {
		f, _ := os.OpenFile(HISTFILE, os.O_RDONLY, 0666)
		defer f.Close()
		/* Read file and store each line in history slice */
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			text := scanner.Text()
			history = append(history, text)
		}
		// Reversing the history slice
		last := len(history) - 1
		for i := 0; i < len(history)/2; i++ {
			history[i], history[last-i] = history[last-i], history[i]
		}
	}
	return history
}

func exit() {
	f, err := os.OpenFile(HISTFILE, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf(ERRFORMAT, err.Error())
	}
	defer f.Close()
	for _, i := range HISTMEM {
		f.Write([]byte(i))
		f.Write([]byte("\n"))
	}
	os.Exit(0)
}
