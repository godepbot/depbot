// package gomodules does the analysis of Go modules and returns its findings.
package gomodules

import "github.com/godepbot/depbot"

// Find walks the directory three and looks for go.mod files
// to then parse dependencies and return them.
func FindDependencies(wd string) ([]depbot.Dependency, error) {
	return []depbot.Dependency{}, nil
}

// d, err := os.TempDir()
// Write go.mod en directory
// Write package/go.mod
// Write package/p2/go.mod
// Check if no go.mod
//
// import (
// 	"fmt"
// 	"io/ioutil"
// 	"path/filepath"
//
// 	"golang.org/x/mod/modfile"
// )
//
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
//
// 			Version: req.Mod.Version,
// 			Source:  path,
// 		}

// 		deps = append(deps, dep)
// 	}

// 	return deps, nil
// }
