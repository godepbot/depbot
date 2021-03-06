package internal

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func PathsFor(wd string, fileNames ...string) []string {
	pths := []string{}

	filepath.WalkDir(wd, func(path string, d fs.DirEntry, err error) error {
		for _, fileName := range fileNames {
			if PathContainsFolder(path, "node_modules") {
				continue
			}

			if filepath.Base(path) == fileName {
				pths = append(pths, path)
			}
		}

		return nil
	})

	return pths
}

func PathContainsFolder(path, folderName string) bool {
	for _, name := range strings.Split(path, string(filepath.Separator)) {
		if name == folderName {
			return true
		}
	}

	return false
}
