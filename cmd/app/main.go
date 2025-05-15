package main

import (
	"fmt"
	"os"
	"github.com/Breinss/Go-Cli-Node/internal/app"
	"github.com/Breinss/Go-Cli-Node/internal/scanner"
	"log"
)

func main() {
	app := app.New()
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	rootPath := "/path/to/root" // Replace with actual root path
	scanner := scanner.New(rootPath)

	// Confirm and scan
	confirmed, err := scanner.ConfirmScan()
	if err != nil {
		log.Println("Error confirming scan:", err)
		return
	}

	if confirmed {
		if err := scanner.Scan(); err != nil {
			log.Println("Error during scan:", err)
			return
		}

		scanner.DisplayResults()

		// Add this line to ask for pruning after displaying results
		if err := scanner.AskForPruning(); err != nil {
			log.Println("Error during pruning:", err)
		}
	}
}
