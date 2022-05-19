// package gomodules does the analysis of Go modules and returns its findings.
package gomodules

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/godepbot/depbot"
	"golang.org/x/mod/modfile"
)

const (
	goDependencyFile string = "go.mod"
	dependencyNameGo string = "Go"
)

// Find walks the directory three and looks for go.mod files
// to then parse dependencies and return them.
func FindDependencies(wd string) (depbot.Dependencies, error) {
	pths := []string{}

	filepath.WalkDir(wd, func(path string, d fs.DirEntry, err error) error {
		if strings.Contains(path, goDependencyFile) {
			pths = append(pths, path)
		}

		return nil
	})

	dependencies := depbot.Dependencies{}

	for _, p := range pths {
		d, err := ioutil.ReadFile(p)

		if err != nil {
			fmt.Println("Error reading the file.")
			return dependencies, err
		}

		f, err := modfile.Parse(p, d, nil)

		if err != nil {
			fmt.Println("Error parsing the file.")
			return dependencies, err
		}

		dependencies = append(dependencies, depbot.Dependency{
			File:    goDependencyFile,
			Version: f.Go.Version,
			Name:    dependencyNameGo,
			Kind:    depbot.DependencyKindLanguage,
		})

		for _, r := range f.Require {
			dependencies = append(dependencies, depbot.Dependency{
				File:    goDependencyFile,
				Name:    r.Mod.Path,
				Version: r.Mod.Version,
				Kind:    depbot.DependencyKindLibrary,
				Direct:  !r.Indirect,
			})
		}
	}

	if len(deps) > 0 {
		fmt.Println("Total Dependencies found for this project:", len(deps))
	} else {
		fmt.Println("No Go Dependencies were found for this project.")
	}

	for _, d := range deps {
		fmt.Println(d.Name, d.Version)
	}

	return dependencies, nil
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
