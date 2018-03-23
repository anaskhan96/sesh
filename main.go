package main

import (
	"bufio"
	"bytes"
	"errors"
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
	CONFIG    string
	aliases   map[string]string
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
	sesh_config()
}

func sesh_config() {
	CONFIG = fmt.Sprintf("%s/%s", os.Getenv("HOME"), ".seshrc")
	aliases = make(map[string]string)
	aliases["~"] = os.Getenv("HOME")
	if _, err := os.Stat(CONFIG); err == nil {
		f, _ := os.OpenFile(CONFIG, os.O_RDONLY, 0666)
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			parseLine(scanner.Text())
		}
	}
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
		args, err := parseLine(line)
		if err != nil {
			fmt.Printf(ERRFORMAT, err.Error())
		}
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

func parseLine(line string) ([]string, error) {
	/* Need to include whitespaces as single token inside double quotes */
	args := strings.Fields(line)
	if args[0] == "alias" {
		for _, i := range args[1:] {
			aliasArgs := strings.Split(i, "=")
			if len(aliasArgs) != 2 {
				return nil, errors.New("Wrong alias key and value format in config")
			}
			aliases[aliasArgs[0]] = aliasArgs[1]
		}
		return args, nil
	}
	if args[0] == "export" {
		exportArgs := strings.Split(args[1], "=")
		if len(exportArgs) != 2 {
			return nil, errors.New("Wrong export format in config")
		}
		os.Setenv(exportArgs[0], exportArgs[1])
		return args, nil
	}
	// replace if an alias
	for i, arg := range args {
		if val, ok := aliases[arg]; ok {
			args[i] = val
		}
	}
	// replace if an environment variable
	for i, arg := range args {
		if arg[0] == '$' {
			args[i] = os.Getenv(arg[1:])
		}
	}
	return args, nil
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

func aliasing(line string) string {
	for key, value := range aliases {
		if strings.Contains(line, key) {
			line = strings.Replace(line, key, value, -1)
		}
	}
	return line
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
