// Scanner package is in charge of identifying dependencies
// within a repository. It would look at different kind of
// dependency files and identify the dependencies within them.
package scanner

import (
	"fmt"
	"io/fs"
	"log"
)

// Scan a repository for files that match any of the dependency
// parsers, it will call parser.CanParse() to check if the file
// is a dependency file that could be parsed.
func Scan(fldr fs.FS, filterer func(string) bool) ([]string, error) {
	files := []string{}
	err := fs.WalkDir(fldr, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if d.IsDir() {
			return nil
		}

		if !filterer(path) {
			return nil
		}

		files = append(files, path)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking the folder while scanning: %w", err)
	}

	return files, nil
}
