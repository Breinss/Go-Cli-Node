package app

import (
    "flag"
    "fmt"
    "os"

    "github.com/Breinss/Go-Cli-Node/internal/scanner"
    "github.com/charmbracelet/lipgloss"
)

var (
    purple    = lipgloss.Color("99")
    gray      = lipgloss.Color("245")
    lightGray = lipgloss.Color("241")

    headerStyle  = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
    cellStyle    = lipgloss.NewStyle().Padding(0, 1).Width(14)
    oddRowStyle  = cellStyle.Foreground(gray)
    evenRowStyle = cellStyle.Foreground(lightGray)
)

// App represents the application
type App struct {
    rootPath string // Path to start scanning for node_modules
}

// New creates a new instance of the application
func New() *App {
    return &App{}
}

// Run executes the application logic
func (a *App) Run() error {
    // Define command-line flags
	
    scanCmd := flag.NewFlagSet("scan", flag.ExitOnError)
    rootPathPtr := scanCmd.String("path", "/", "Root path to scan for node_modules directories")
    
    // Check which command is being run
    if len(os.Args) < 2 {
        fmt.Println("Expected 'scan' subcommand")
        return nil
    }

    switch os.Args[1] {
    case "scan":
        scanCmd.Parse(os.Args[2:])
        a.rootPath = *rootPathPtr
        return a.scanNodeModules()
    default:
        fmt.Println("Expected 'scan' subcommand")
        return nil
    }
}

// scanNodeModules scans for node_modules directories and displays the results
func (a *App) scanNodeModules() error {
    // Create a new scanner with the root path
    nodeScanner := scanner.New(a.rootPath)
    
    // Ask for confirmation before scanning
    confirmed, err := nodeScanner.ConfirmScan()
    if err != nil {
        return err
    }
    
    if !confirmed {
        fmt.Println("Scan operation cancelled by user.")
        return nil
    }
    
    // Perform the scan
    if err := nodeScanner.Scan(); err != nil {
        return err
    }
    
    // Display the results
    nodeScanner.DisplayResults()
    
    return nil
}