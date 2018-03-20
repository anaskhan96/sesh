package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"
)

const (
	TOKDELIM  = " \t\r\n\a"
	ERRFORMAT = "sesh: %s\n"
)

var (
	HISTSIZE  = 25
	HISTFILE  string
	HISTMEM   []string
	HISTCOUNT int
	HISTLINE  string
)

func main() {
	sesh_setup()

	sesh_loop()

	exit()
}

func sesh_setup() {
	os.Clearenv()
	os.Setenv("SHELL", "/bin/sh")

	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf(ERRFORMAT, err.Error())
	}
	wdSlice := strings.Split(wd, "/")
	os.Setenv("CWD", wdSlice[len(wdSlice)-1])

	os.Setenv("PATH", "/bin:/usr/bin:/usr/local/bin/")
	currUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	os.Setenv("USER", currUser.Username)
	os.Setenv("HOME", currUser.HomeDir)
	HISTFILE = fmt.Sprintf("%s/%s", os.Getenv("HOME"), ".sesh_history")

	/* Importing config */
	//sesh_config()
}

func sesh_loop() {
	HISTMEM = initHistory(HISTMEM)
	status := 1
	reader := bufio.NewReader(os.Stdin)
	//go listenForKeyPress()

	for status != 0 {
		symbol := "\u2713"
		if status == 2 {
			symbol = "\u2715"
		}
		fmt.Printf("sesh 🔥  %s %s ", os.Getenv("CWD"), symbol)
		line, _ := reader.ReadString('\n')
		line = line[:len(line)-1]
		HISTLINE = line
		args := parseLine(line)
		status = execute(args)
		if status == 1 {
			/* Store line in history */
			if HISTCOUNT == HISTSIZE {
				HISTMEM = HISTMEM[1:]
				HISTCOUNT = 0
			}
			HISTMEM = append([]string{HISTLINE}, HISTMEM...)
			HISTCOUNT++
		}
	}
}

func parseLine(line string) []string {
	line = strings.Replace(line, "~", os.Getenv("HOME"), -1)
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
	timestamp := time.Now().String()
	if err := cmd.Run(); err != nil {
		fmt.Printf(ERRFORMAT, err.Error())
		return 2
	}
	fmt.Printf(out.String())
	HISTLINE = fmt.Sprintf("%d::%s::%s", cmd.Process.Pid, timestamp, HISTLINE)
	return 1
}

func execute(args []string) int {
	if len(args) == 0 {
		return 1
	}
	for k, v := range builtins {
		if args[0] == k {
			timestamp := time.Now().String()
			HISTLINE = fmt.Sprintf("%d::%s::%s", os.Getpid(), timestamp, HISTLINE)
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
			history = append(history, scanner.Text())
			HISTCOUNT++
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
