package main

import (
	"os"
	"runtime"

	commands "./commands"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}

}
