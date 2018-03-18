package main

import (
	"fmt"
	"os"
	"strings"
)

var builtins = map[string]func([]string) int{
	"exit": sesh_exit,
	"cd":   sesh_cd,
	"help": sesh_help,
	//"pwd":  sesh_pwd,
	//"echo": sesh_echo,
	"history": sesh_history,
}

func sesh_exit(args []string) int {
	fmt.Println("Exiting shell")
	return 0
}

func sesh_help(args []string) int {
	fmt.Println("sesh -- simple elegant shell by anaskhan96")
	return 1
}

func sesh_history(args []string) int {
	for _, i := range HISTMEM {
		fmt.Println(i)
	}
	return 1
}

func sesh_cd(args []string) int {
	if len(args) == 0 {
		fmt.Printf(ERRFORMAT, "Please provide a path to change directory to")
	} else if len(args) > 1 {
		fmt.Printf(ERRFORMAT, "Too many args for changing directory")
	} else {
		err := os.Chdir(args[0])
		if err != nil {
			fmt.Printf(ERRFORMAT, err.Error())
			return 2
		}
		wd, err := os.Getwd()
		wdSlice := strings.Split(wd, "/")
		os.Setenv("CWD", wdSlice[len(wdSlice)-1])
	}
	return 1
}

/* Commenting extra builtins for now
func sesh_pwd(args []string) int {
	if len(args) != 0 {
		fmt.Printf(ERRFORMAT, "pwd expects 0 args")
		return 2
	}
	dir, _ := os.Getwd()
	absPath, _ := filepath.Abs(dir)
	fmt.Println(absPath)
	return 1
}

func sesh_echo(args []string) int {
	for _, i := range args {
		fmt.Printf("%s ", i)
	}
	fmt.Printf("\n")
	return 1
}*/
