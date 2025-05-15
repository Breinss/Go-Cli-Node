package helpers

import "os"

// RemoveDirectory deletes a directory and all its contents
func RemoveDirectory(path string) error {
	return os.RemoveAll(path)
}
