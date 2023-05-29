package main

import (
	"errors"
	"fmt"
	"io"
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

	skipExt := make(map[string]bool)

	skipDir := make(map[string]bool)
	skipDir[".git"] = true
	skipDir[".gitattributes"] = true
	skipDir[".gitignore"] = true

	skipDir[".github"] = true

	skipDir["scripts"] = true

	err = os.MkdirAll(buildDir, os.FileMode(0766))
	if err != nil && !errors.Is(err, os.ErrExist) {
		fmt.Printf("Cannot create directory: %v\n", err)
		os.Exit(1)
	}

	filepath.WalkDir(cwd, func(path string, info fs.DirEntry, err error) error {
		if path == cwd {
			return nil
		}

		if info.IsDir() {
			if skipDir[info.Name()] {
				return filepath.SkipDir
			}
		}

		if skipExt[filepath.Ext(path)] ||
			filepath.Dir(path) == cwd {
			return nil
		}

		fmt.Printf("  %s", path)

		err = fileCopy(filepath.Join(buildDir, info.Name()), path)
		if err != nil {
			fmt.Printf("    ERROR:\n%v\n", err)
		} else {
			fmt.Println("    SUCCESS")
		}

		return nil
	})
}

func fileCopy(dest, src string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
