package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const TOKDELIM = " \t\r\n\a"

var CWD string

func main() {
	/* Importing config */

	sesh_loop()

	os.Exit(0)
}

func sesh_loop() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	wdSlice := strings.Split(wd, "/")
	CWD = wdSlice[len(wdSlice)-1]

	reader := bufio.NewReader(os.Stdin)
	status := true

	for status {
		fmt.Printf("sesh 🔥 %s -> ", CWD)
		line, _ := reader.ReadString('\n')
		args := strings.Split(line, TOKDELIM)
		status = execute(args)
	}
}

func execute(args []string) bool {
	return false
}
