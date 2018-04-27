package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func launch(args []string) int {
	commands, isBackground := make([]*exec.Cmd, 0, 5), false
	if args[len(args)-1] == "&" {
		isBackground = true
	}
	start, cmdInEnd := 0, true
	for i, arg := range args {
		if i == len(args)-1 && cmdInEnd {
			if len(commands) == 0 {
				return launchSimpleCommand(args, isBackground)
			}
			cmd := exec.Command(args[start], args[start+1:]...)
			commands = append(commands, cmd)
		} else if arg == "|" {
			cmd := exec.Command(args[start], args[start+1:i]...)
			commands = append(commands, cmd)
			start = i + 1
		} else if arg == ">" || arg == ">>" {
			cmd := exec.Command(args[start], args[start+1:i]...)
			var f *os.File
			if arg == ">" {
				f, _ = os.OpenFile(args[i+1], os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
			} else {
				f, _ = os.OpenFile(args[i+1], os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
			}
			cmd.Stdout = f
			commands = append(commands, cmd)
			cmdInEnd = false
		} else if arg == "<" {
			f, _ := os.Open(args[i+1])
			if len(commands) > 0 {
				commands[0].Stdin = f
			} else {
				cmd := exec.Command(args[start], args[start+1:i]...)
				cmd.Stdin = f
				commands = append(commands, cmd)
			}
			cmdInEnd = false
		}
	}
	for i := range commands {
		if i != len(commands)-1 {
			if commands[i+1].Stdin == nil {
				commands[i+1].Stdin, _ = commands[i].StdoutPipe()
			}
		} else {
			if commands[i].Stdout == nil {
				commands[i].Stdout = os.Stdout
			}
		}
	}
	for i := len(commands) - 1; i > 0; i-- {
		commands[i].Start()
	}
	timestamp := time.Now().String()
	if !isBackground {
		if err := commands[0].Run(); err != nil {
			fmt.Printf(ERRFORMAT, err.Error())
			return 2
		}
	} else {
		commands[0].Start()
	}
	for i := range commands[1:] {
		commands[i].Wait()
	}
	HISTLINE = fmt.Sprintf("%d::%s::%s", commands[0].Process.Pid, timestamp, HISTLINE)
	return 1
}

func launchSimpleCommand(args []string, isBackground bool) int {
	// Spawning and executing a process
	cmd := exec.Command(args[0], args[1:]...)
	// Setting stdin, stdout, and stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = nil // making sure the command uses the current process' environment
	timestamp := time.Now().String()
	if !isBackground {
		if err := cmd.Run(); err != nil {
			fmt.Printf(ERRFORMAT, err.Error())
			return 2
		}
	} else {
		cmd.Start()
	}
	HISTLINE = fmt.Sprintf("%d::%s::%s", cmd.Process.Pid, timestamp, HISTLINE)
	return 1
}
