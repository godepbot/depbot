// package gomodules does the analysis of Go modules and returns its findings.
package gomodules

import "github.com/godepbot/depbot"

// Find walks the directory three and looks for go.mod files
// to then parse dependencies and return them.
func FindDependencies(wd string) ([]depbot.Dependency, error) {
	return []depbot.Dependency{}, nil
}

// import (
// 	"depbot/app/models"
// 	"fmt"
// 	"io/ioutil"
// 	"path/filepath"

// 	"golang.org/x/mod/modfile"
// )

// // CanParse a path file, returns true if the file
// // can be parsed by the Parse method, this is useful
// // when looking at a directory for dependency files.
// func CanParse(path string) bool {
// 	return filepath.Base(path) == "go.mod"
// }

// func Parse(path string) ([]models.Dependency, error) {
// 	file, err := ioutil.ReadFile(path)
// 	if err != nil {
// 		return nil, fmt.Errorf("could not read %s: %w\n", path, err)
// 	}

// 	f, err := modfile.Parse(path, file, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("could not read %v: %w\n", path, err)
// 	}

// 	deps := []models.Dependency{}
// 	for _, req := range f.Require {
// 		dep := models.Dependency{
// 			Library: models.Library{
// 				Name:     req.Mod.Path,
// 				Language: models.LanguageGo,
// 			},

// 			Version: req.Mod.Version,
// 			Source:  path,
// 		}

// 		deps = append(deps, dep)
// 	}

// 	return deps, nil
// }
