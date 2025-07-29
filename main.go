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

	// Define flags
	sortDescending := flag.Bool("d", false, "sort by size descending")
	sortAscending := flag.Bool("a", false, "sort by size ascending")

	// Expand grouped flags like -ad into -a -d
	var expandedArgs []string
	for _, arg := range os.Args {
		if len(arg) > 2 && arg[0] == '-' && arg[1] != '-' {
			// e.g., -ad becomes -a -d
			for _, ch := range arg[1:] {
				expandedArgs = append(expandedArgs, "-"+string(ch))
			}
		} else {
			expandedArgs = append(expandedArgs, arg)
		}
	}
	os.Args = expandedArgs

	flag.Parse()

	// Ensure only one sort flag is used
	if *sortDescending && *sortAscending {
		fmt.Fprintln(os.Stderr, "Error: cannot use both -a (ascending) and -d (descending) at the same time.")
		os.Exit(1)
	}

	// Check number of positional args (target directory)
	args := flag.Args()
	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "Error: only one target directory may be specified.")
		os.Exit(1)
	}
	targetDir := "."
	if len(args) == 1 {
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
			size, err := parser.WalkDirSize(dirPath)
			if err != nil {
				fmt.Printf("Error calculating size for %s: %v\n", dirPath, err)
				continue
			}
			dirs = append(dirs, dirInfo{Name: entry.Name(), Size: size})
		}
	}

	// Apply sorting
	if *sortDescending {
		sort.Slice(dirs, func(i, j int) bool {
			return dirs[i].Size > dirs[j].Size
		})
	} else {
		// default and -a
		sort.Slice(dirs, func(i, j int) bool {
			return dirs[i].Size < dirs[j].Size
		})
	}

	// Print sorted results
	for _, d := range dirs {
		fmt.Printf("%-20s %s\n", parser.FormatBytes(d.Size), d.Name)
	}

	// Print total size of target dir
	totalTargetSize, err := parser.WalkDirSize(targetDir)
	if err != nil {
		fmt.Printf("Error calculating total size for %s: %v\n", targetDir, err)
	} else {
		fmt.Printf("\n%-20s Total size of: %s\n", parser.FormatBytes(totalTargetSize), targetDir)
	}
}
