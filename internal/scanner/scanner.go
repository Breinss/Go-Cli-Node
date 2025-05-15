package scanner

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/Breinss/Go-Cli-Node/pkg/helpers"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/charmbracelet/log"
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

// Scanner represents a node_modules scanner
type Scanner struct {
	rootPath  string
	results   []helpers.NodeModuleInfo
	totalSize int64
	duration  time.Duration
}

// New creates a new scanner instance
func New(rootPath string) *Scanner {
	return &Scanner{
		rootPath: rootPath,
	}
}

// ConfirmScan asks for user confirmation before scanning
func (s *Scanner) ConfirmScan() (bool, error) {
	var confirmed bool
	if err := huh.NewConfirm().
		Title("Are you sure you want to scan for node_modules?").
		Description("This operation may take some time depending on the directory size.").
		Affirmative("Yes, scan").
		Negative("No, cancel").
		Value(&confirmed).
		Run(); err != nil {
		return false, fmt.Errorf("confirmation error: %v", err)
	}

	return confirmed, nil
}

// Scan performs the scanning operation
func (s *Scanner) Scan() error {
	log.Info(fmt.Sprintf("Scanning for node_modules directories from: %s", s.rootPath))

	var scanErr error
	startTime := time.Now()

	// Define action function for the spinner
	action := func() {
		s.results, scanErr = helpers.ScanNodeModules(s.rootPath)
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

	s.duration = time.Since(startTime)

	// Calculate total size from results
	for _, result := range s.results {
		s.totalSize += result.Size
	}

	return nil
}

// DisplayResults shows all scan results
func (s *Scanner) DisplayResults() {
	log.Info(fmt.Sprintf("Found %d node_modules directories in %s", len(s.results), s.duration))
	log.Info(fmt.Sprintf("Total space used: %s", helpers.FormatSize(s.totalSize)))

	s.DisplayLargestDirectories(30)
	s.DisplaySizeDistribution()
}

// AskForPruning asks the user if they want to prune node_modules directories
func (s *Scanner) AskForPruning() error {
	if len(s.results) == 0 {
		return nil
	}

	var wantToPrune bool
	if err := huh.NewConfirm().
		Title("Would you like to prune node_modules directories?").
		Description(fmt.Sprintf("You can free up to %s of disk space", helpers.FormatSize(s.totalSize))).
		Affirmative("Yes, let me select directories").
		Negative("No, keep them all").
		Value(&wantToPrune).
		Run(); err != nil {
		return fmt.Errorf("pruning confirmation error: %v", err)
	}

	if wantToPrune {
		return s.PruneNodeModules()
	}
	
	return nil
}

// DisplayLargestDirectories shows the largest directories as a tree
func (s *Scanner) DisplayLargestDirectories(topCount int) {
	if len(s.results) < topCount {
		topCount = len(s.results)
	}

	fmt.Println()
	fmt.Println(headerStyle.Render("Largest node_modules Directories"))
	fmt.Println()

	// Setup tree styles
	enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("63")).MarginRight(1)
	rootStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("35")).Bold(true)
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	// Create tree with root path
	t := tree.Root(s.rootPath)

	// Add the top N directories to the tree
	for i := 0; i < topCount; i++ {
		// Make path relative to root
		relPath, err := filepath.Rel(s.rootPath, s.results[i].Path)
		if err != nil {
			relPath = s.results[i].Path
		}

		// Add path segments as a nested structure with size info
		label := fmt.Sprintf("%s (%s)", relPath, helpers.FormatSize(s.results[i].Size))

		// Add the directory to the tree
		t.Child(label)
	}

	// Apply styling to the tree
	t.Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(enumeratorStyle).
		RootStyle(rootStyle).
		ItemStyle(itemStyle)

	// Print the tree
	fmt.Println(t)
}

// DisplaySizeDistribution shows the distribution of directories by size
func (s *Scanner) DisplaySizeDistribution() {
	fmt.Println()
	fmt.Println(headerStyle.Render("Size Distribution"))
	fmt.Println()

	// Define size ranges
	ranges := []struct {
		name  string
		min   int64
		max   int64
		count int
		total int64
	}{
		{"Huge (> 1GB)", 1024 * 1024 * 1024, int64(^uint(0) >> 1), 0, 0},
		{"Large (100MB-1GB)", 100 * 1024 * 1024, 1024 * 1024 * 1024, 0, 0},
		{"Medium (10MB-100MB)", 10 * 1024 * 1024, 100 * 1024 * 1024, 0, 0},
		{"Small (1MB-10MB)", 1024 * 1024, 10 * 1024 * 1024, 0, 0},
		{"Tiny (< 1MB)", 0, 1024 * 1024, 0, 0},
	}

	// Count directories in each range
	for _, result := range s.results {
		for i := range ranges {
			if result.Size >= ranges[i].min && result.Size < ranges[i].max {
				ranges[i].count++
				ranges[i].total += result.Size
				break
			}
		}
	}

	// Display distribution
	fmt.Printf("%s %s %s\n",
		cellStyle.Bold(true).Width(20).Render("Size Category"),
		cellStyle.Bold(true).Render("Count"),
		cellStyle.Bold(true).Render("Total Size"))

	for i, r := range ranges {
		rowStyle := oddRowStyle
		if i%2 == 0 {
			rowStyle = evenRowStyle
		}

		fmt.Printf("%s %s %s\n",
			rowStyle.Width(20).Render(r.name),
			rowStyle.Render(fmt.Sprintf("%d", r.count)),
			rowStyle.Render(helpers.FormatSize(r.total)))
	}
}

// PruneNodeModules allows user to select and delete node_modules directories
func (s *Scanner) PruneNodeModules() error {
	if len(s.results) == 0 {
		return fmt.Errorf("no node_modules directories found to prune")
	}

	// Create options for selection form
	options := make([]huh.Option[string], len(s.results))
	for i, result := range s.results {
		// Create a label with path and size
		label := fmt.Sprintf("%s (%s)", result.Path, helpers.FormatSize(result.Size))
		options[i] = huh.NewOption(label, result.Path)
	}

	// Create multi-select form for choosing directories to prune
	var selectedPaths []string
	form := huh.NewMultiSelect[string]().
		Title("Select node_modules directories to prune").
		Options(options...).
		Value(&selectedPaths)

	if err := form.Run(); err != nil {
		return fmt.Errorf("selection error: %v", err)
	}

	if len(selectedPaths) == 0 {
		log.Info("No directories selected for pruning")
		return nil
	}

	// Calculate total size to be pruned
	var totalSizeToPrune int64
	for _, path := range selectedPaths {
		for _, result := range s.results {
			if result.Path == path {
				totalSizeToPrune += result.Size
				break
			}
		}
	}

	// Ask for confirmation before pruning
	var confirmed bool
	confirmForm := huh.NewConfirm().
		Title("Are you sure you want to prune these directories?").
		Description(fmt.Sprintf("This will permanently delete %d directories (total %s)", 
			len(selectedPaths), helpers.FormatSize(totalSizeToPrune))).
		Affirmative("Yes, delete them").
		Negative("No, cancel").
		Value(&confirmed)

	if err := confirmForm.Run(); err != nil {
		return fmt.Errorf("confirmation error: %v", err)
	}

	if !confirmed {
		log.Info("Pruning cancelled")
		return nil
	}

	// Perform pruning with spinner
	var pruneErr error
	action := func() {
		pruneErr = s.executeNodeModulesPruning(selectedPaths)
	}

	if err := spinner.New().
		Title("Pruning node_modules directories...").
		Action(action).
		Run(); err != nil {
		return fmt.Errorf("spinner error: %v", err)
	}

	if pruneErr != nil {
		return fmt.Errorf("pruning error: %v", pruneErr)
	}

	log.Info(fmt.Sprintf("Successfully pruned %d directories, freed %s of disk space", 
		len(selectedPaths), helpers.FormatSize(totalSizeToPrune)))
	
	return nil
}

// executeNodeModulesPruning performs the actual deletion of directories
func (s *Scanner) executeNodeModulesPruning(paths []string) error {
	for _, path := range paths {
		log.Debug(fmt.Sprintf("Removing directory: %s", path))
		if err := helpers.RemoveDirectory(path); err != nil {
			return fmt.Errorf("failed to remove %s: %v", path, err)
		}
	}
	return nil
}
