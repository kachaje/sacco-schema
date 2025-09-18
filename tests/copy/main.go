package main

import (
	"flag"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	var rootFolder string = filepath.Join("..", "..")

	flag.StringVar(&rootFolder, "f", rootFolder, "root folder")

	copyFile := func(src, dst string) error {
		sourceFile, err := os.Open(src)
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		destinationFile, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer destinationFile.Close()

		_, err = io.Copy(destinationFile, sourceFile)
		if err != nil {
			return err
		}

		os.Remove(src)

		return nil
	}

	err := filepath.WalkDir(rootFolder, func(pathName string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if dir.IsDir() {
			if regexp.MustCompile(`\.git*|\.vscode|tests|\.[A-Za-z]+`).MatchString(pathName) {
				return filepath.SkipDir
			}
		}

		filename := filepath.Base(pathName)
		folderName := filepath.Dir(pathName)

		if strings.HasSuffix(filename, "_test.go") {
			return copyFile(pathName, filepath.Join(rootFolder, "tests", filename))
		} else if regexp.MustCompile(`fixtures`).MatchString(folderName) {
			return copyFile(pathName, filepath.Join(rootFolder, "tests", "fixtures", filename))
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}
