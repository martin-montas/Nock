package cmd

import (
	"fmt"
	"os"
)

func Execute() {
	command := os.Args[1]
	if command == "version" || command == "-v" || command == "--version" {
		fmt.Println("Nock", Version)
		return
	}
	if command == "" {
		fmt.Println("Expected 'dir' , 'crawl' or 'version' subcommand")
		return
	}

	d, ok := registry[command]
	if !ok {
		fmt.Println("Unknown command:", command)
		return
	}
	Parse(d)
}
