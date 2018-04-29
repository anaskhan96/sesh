package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var builtins = map[string]func([]string) int{
	"exit":    sesh_exit,
	"cd":      sesh_cd,
	"help":    sesh_help,
	"history": sesh_history,
	"walk":    sesh_walk,
	"show":    sesh_show,
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

func sesh_walk(args []string) int {
	var dir string
	if len(args) == 0 || args[0] == "." {
		dir, _ = filepath.Abs("")
	} else if args[0] == ".." {
		currDir, _ := filepath.Abs("")
		dir = filepath.Dir(currDir)
	} else {
		dir, _ = filepath.Abs(args[0])
	}
	if fi, err := os.Stat(dir); err == nil {
		if fi.Mode().IsDir() {
			return traverse(dir)
		}
		fmt.Printf(ERRFORMAT, "Not a directory")
		return 2
	}
	fmt.Printf(ERRFORMAT, "Invalid path")
	return 2
}

func sesh_show(args []string) int {
	prefix := ""
	if len(args) > 1 {
		fmt.Printf(ERRFORMAT, "wrong usage of show")
		return 2
	} else if len(args) == 1 {
		prefix = args[0]
	}
	dirs := strings.Split(os.Getenv("PATH"), ":")
	commands := make([]string, 0, 10)
	for _, dir := range dirs {
		files, _ := ioutil.ReadDir(dir)
		for _, file := range files {
			if strings.HasPrefix(file.Name(), prefix) {
				commands = append(commands, file.Name())
			}
		}
	}
	for _, command := range commands {
		fmt.Printf("%s\t", command)
	}
	fmt.Println()
	return 1
}

func traverse(dir string) int {
	dashes, _ := "|", filepath.Base(dir)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		name := filepath.Base(path)
		// TODO: don't show hidden files and directories in tree
		/*if (name != "." && name != "..") && name[0] == '.' {
			return filepath.SkipDir
		}*/
		if info.IsDir() {
			dashes += "--"
		}
		fmt.Printf("%s %s\n", dashes, name)
		return nil
	})
	return 1
}
