package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const buildDir = "payload"

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Cannot get working directory: %v\n", err)
		os.Exit(1)
	}

	onlyExt := make(map[string]bool)
	onlyExt[".sha1"] = true
	onlyExt[".sha256"] = true
	onlyExt[".md5"] = true

	skipDir := make(map[string]bool)
	skipDir[".git"] = true
	skipDir[".gitattributes"] = true
	skipDir[".gitignore"] = true

	skipDir[".github"] = true

	skipDir[".vscode"] = true

	skipDir[".stfolder"] = true

	skipDir["scripts"] = true

	buildDirAbs, err := filepath.Abs(buildDir)
	if err != nil {
		fmt.Printf("Cannot get absolute path: %v\n", err)
		os.Exit(1)
	}

	remove(buildDirAbs)

	filepath.WalkDir(cwd, func(path string, info fs.DirEntry, err error) error {
		if path == cwd {
			return nil
		}

		if info.IsDir() {
			if skipDir[info.Name()] {
				return filepath.SkipDir
			}
		}

		if !onlyExt[filepath.Ext(path)] ||
			filepath.Dir(path) == cwd {
			return nil
		}

		remove(path)

		return nil
	})
}

func remove(path string) {
	fmt.Printf("  %s", path)

	err := os.RemoveAll(path)
	if err != nil {
		fmt.Printf("    ERROR:\n%v\n", err)
	}

	fmt.Println("    SUCCESS")
}
