package uspscan

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func searchPaths(servicePath string) ([]string, error) {
	re := regexp.MustCompile("'.+'|\".+\"|\\S+")
	parts := re.FindAllString(servicePath, -1)
	if len(parts) == 0 {
		return nil, fmt.Errorf("no path elements in %s", servicePath)
	}
	paths := make([]string, 0, len(parts))
	for i := 1; i <= len(parts); i++ {
		restoredPath := strings.Join(parts[:i], " ")
		restoredPath = strings.Replace(restoredPath, "\"", "", -1)
		restoredPath = strings.Replace(restoredPath, "'", "", -1)
		paths = append(paths, restoredPath)
	}
	return paths, nil
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if _, err := os.Stat(path + ".exe"); os.IsNotExist(err) {
			return false
		}

	}
	return true
}

func fileCreationPossible(path string) bool {
	file, err := os.OpenFile(path, os.O_CREATE, 0666)
	if err != nil {
		if os.IsPermission(err) {
			return false
		}
		fmt.Printf("WARNING: non-permission error: %v\n", err)
		return false
	}
	file.Close()
	err = os.Remove(path)
	if err != nil {
		panic(fmt.Sprintf("could not remove the created file: %v", path))
	}
	return true
}
