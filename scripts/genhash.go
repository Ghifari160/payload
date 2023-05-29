package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Cannot get working directory: %v\n", err)
		os.Exit(1)
	}

	skipExt := make(map[string]bool)
	skipExt[".sha1"] = true
	skipExt[".sha256"] = true
	skipExt[".md5"] = true

	skipDir := make(map[string]bool)
	skipDir[".git"] = true
	skipDir[".gitattributes"] = true
	skipDir[".gitignore"] = true

	skipDir[".github"] = true

	skipDir["scripts"] = true

	filepath.WalkDir(cwd, func(path string, info fs.DirEntry, err error) error {
		if path == cwd {
			return nil
		}

		if path != cwd && info.IsDir() {
			if skipDir[info.Name()] {
				return filepath.SkipDir
			}
		}

		if skipExt[filepath.Ext(path)] ||
			filepath.Dir(path) == cwd {
			return nil
		}

		fmt.Printf("  %s", path)

		err = hash(path)
		if err != nil {
			fmt.Printf("    ERROR:\n%v\n", err)
		} else {
			fmt.Println("    SUCCESS")
		}

		return nil
	})
}

// hash runs hashing utilities against filePath.
func hash(filePath string) error {
	err := hasher("sha1sum", filePath, filePath+".sha1")
	if err != nil {
		return err
	}

	err = hasher("sha256sum", filePath, filePath+".sha256")
	if err != nil {
		return err
	}

	err = hasher("md5sum", filePath, filePath+".md5")
	if err != nil {
		return err
	}

	return nil
}

// hasher runs the provided utility against filePath.
// hasher saves the output of the utility to outPath.
func hasher(utility, filePath, outPath string) error {
	sumfile, err := os.Create(outPath)
	if err != nil {
		return err
	}

	cmd := exec.Command(utility, filepath.Base(filePath))
	cmd.Dir = filepath.Dir(filePath)

	sum, err := cmd.Output()
	if err != nil {
		return err
	}

	_, err = sumfile.Write(sum)
	if err != nil {
		return err
	}

	return nil
}
