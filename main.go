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

	sortDesc := flag.Bool("s", false, "sort by size descending")
	flag.Parse()

	// Get optional positional argument as directory (defaults to ".")
	args := flag.Args()
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	fmt.Printf("Calculating disk usage for: %s\n", targetDir)

	// Get immediate subdirectories and their sizes
	entries, err := os.ReadDir(targetDir)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			dirPath := filepath.Join(targetDir, entry.Name())
			size, err := parser.DirSize(dirPath)
			if err != nil {
				fmt.Printf("Error calculating size for %s: %v\n", dirPath, err)
				continue
			}
			dirs = append(dirs, dirInfo{Name: entry.Name(), Size: size})
		}
	}

	// Sort ascending or descending
	sort.Slice(dirs, func(i, j int) bool {
		if *sortDesc {
			return dirs[i].Size > dirs[j].Size
		}
		return dirs[i].Size < dirs[j].Size
	})

	// Print sorted results
	for _, d := range dirs {
		fmt.Printf("%-20s %s\n", parser.FormatBytes(d.Size), d.Name)
	}

	// Optionally, get the total size of the target directory itself
	totalTargetSize, err := parser.DirSize(targetDir)
	if err != nil {
		fmt.Printf("Error calculating total size for %s: %v\n", targetDir, err)
	} else {
		fmt.Printf("\n%-20s Total size of: %s\n", parser.FormatBytes(totalTargetSize), targetDir)
	}
}
