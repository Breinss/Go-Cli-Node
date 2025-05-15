package app
// import "github.com/charmbracelet/log"
import (
	"flag"
	"fmt"
)

// App represents the application
type App struct {
	// Add fields as needed
}

// New creates a new instance of the application
func New() *App {
	return &App{}
}

// Run executes the application logic
func (a *App) Run() error {
	// Parse command-line flags
	flag.Parse()

	log.Info("Hello from the CLI app!")
	return nil
}
