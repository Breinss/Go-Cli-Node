package main

import (
	"fmt"
	"os"
	"github.com/Breinss/Go-Cli-Node/internal/app"
)

func main() {
	app := app.New()
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
