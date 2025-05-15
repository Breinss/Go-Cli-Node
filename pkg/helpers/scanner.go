package helpers

import (
    "fmt"
    "os"
    "path/filepath"
    "sort"
    "sync"
)

// NodeModuleInfo holds information about a node_modules directory
type NodeModuleInfo struct {
    Path string
    Size int64
}

// ScanNodeModules scans the filesystem starting from the given root path
// and returns information about all node_modules directories found
func ScanNodeModules(rootPath string) ([]NodeModuleInfo, error) {
    var results []NodeModuleInfo
    var mutex sync.Mutex
    var wg sync.WaitGroup

    err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            // Skip directories we can't access
            return nil
        }

        if info.IsDir() && info.Name() == "node_modules" {
            wg.Add(1)
            go func(p string) {
                defer wg.Done()
                
                size, err := getDirSize(p)
                if err != nil {
                    fmt.Printf("Error calculating size for %s: %v\n", p, err)
                    return
                }
                
                mutex.Lock()
                results = append(results, NodeModuleInfo{
                    Path: p,
                    Size: size,
                })
                mutex.Unlock()
            }(path)
            
            // Skip descending into node_modules directories to avoid counting nested ones
            return filepath.SkipDir
        }
        
        return nil
    })

    wg.Wait()
    
    // Sort results by size (largest first)
    sort.Slice(results, func(i, j int) bool {
        return results[i].Size > results[j].Size
    })
    
    return results, err
}

// getDirSize calculates the total size of a directory recursively
func getDirSize(path string) (int64, error) {
    var size int64
    
    err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
        if err != nil {
            // Skip files/dirs we can't access
            return nil
        }
        
        if !info.IsDir() {
            size += info.Size()
        }
        
        return nil
    })
    
    return size, err
}

// FormatSize converts byte size to a human-readable format
func FormatSize(bytes int64) string {
    const (
        KB = 1024
        MB = KB * 1024
        GB = MB * 1024
    )
    
    switch {
    case bytes >= GB:
        return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
    case bytes >= MB:
        return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
    case bytes >= KB:
        return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
    default:
        return fmt.Sprintf("%d bytes", bytes)
    }
}
