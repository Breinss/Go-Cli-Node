package app

import (
    "flag"
    "fmt"
    "os"
    "time"

    "github.com/charmbracelet/lipgloss"
    "github.com/Breinss/Go-Cli-Node/pkg/helpers/scanner"
	"github.com/charmbracelet/huh/spinner"
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
    fmt.Printf("Scanning for node_modules directories from: %s\n", a.rootPath)
    
    var results []scanner.NodeModulesInfo
    var scanErr error
    startTime := time.Now()
    
    // Define action function for the spinner
    action := func() {
        results, scanErr = helpers.ScanNodeModules(a.rootPath)
    }
    
    // Run the spinner with our scanning action
    if err := spinner.New().
        Title("Scanning for node_modules... This may take some time").
        Action(action).
        Run(); err != nil {
        return fmt.Errorf("spinner error: %v", err)
    }
    
    // After spinner completes, check if scan had an error
    if scanErr != nil {
        return fmt.Errorf("error during scan: %v", scanErr)
    }
    
    duration := time.Since(startTime)
    
    fmt.Printf("\nFound %d node_modules directories in %s\n\n", len(results), duration)
    
    // Calculate total size from results
    var totalSize int64
    for _, result := range results {
        totalSize += result.Size
    }
    
    fmt.Printf("\nTotal space used: %s\n", helpers.FormatSize(totalSize))
    
    return nil
}