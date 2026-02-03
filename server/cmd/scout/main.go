package main

import (
	"os"

	"docker-scout/internal/app"
)

func main() {
	if isCLIMode(os.Args[1:]) {
		app.RunCLI()
		return
	}
	app.RunServer()
}

func isCLIMode(args []string) bool {
	for _, arg := range args {
		if arg == "--cli" || arg == "cli" {
			return true
		}
	}
	return false
}
