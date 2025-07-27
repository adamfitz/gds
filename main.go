package main

import (
	"flag"
	"fmt"
	"gds/parser"
	"os"
	"path/filepath"
	"sort"
)

func main() {

	type dirInfo struct {
		Name string
		Size int64
	}

	var dirs []dirInfo

	targetDir := flag.String("t", "", "target directory")
	flag.Parse()

	fmt.Printf("Calculating disk usage for: %s\n", *targetDir)

	// Get immediate subdirectories and their sizes
	entries, err := os.ReadDir(*targetDir)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			dirPath := filepath.Join(*targetDir, entry.Name())
			size, err := parser.DirSize(dirPath)
			if err != nil {
				fmt.Printf("Error calculating size for %s: %v\n", dirPath, err)
				continue
			}
			dirs = append(dirs, dirInfo{Name: entry.Name(), Size: size})
		}
	}

	// Sort by size descending (largest first)
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Size > dirs[j].Size
	})

	// Print sorted results
	for _, d := range dirs {
		fmt.Printf("%-20s %s\n", parser.FormatBytes(d.Size), d.Name)
	}

	// Optionally, get the total size of the target directory itself
	totalTargetSize, err := parser.DirSize(*targetDir)
	if err != nil {
		fmt.Printf("Error calculating total size for %s: %v\n", *targetDir, err)
	} else {
		fmt.Printf("\nTotal size of %s: %s\n", *targetDir, parser.FormatBytes(totalTargetSize))
	}
}
