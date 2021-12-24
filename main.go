package main

import (
	"os"

	"github.com/jgarrone/foaas-limiting/cmd"
)

func main() {
	if err := cmd.RootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
