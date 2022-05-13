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
	GoDependencyFile string = ".mod"
	DependencyNameGo string = "Go"
)

// Find walks the directory three and looks for go.mod files
// to then parse dependencies and return them.
func FindDependencies(wd string) (depbot.Dependencies, error) {
	pths := []string{}

	filepath.WalkDir(wd, func(path string, d fs.DirEntry, err error) error {
		if strings.Contains(path, GoDependencyFile) {
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

		dependencies = append(dependencies, languageDependency(f.Go))

		for _, r := range f.Require {
			dependencies = append(dependencies, libraryDependency(r))
		}

		fmt.Println("Dependencies found:", len(f.Require), "for file: ", p)
	}

	return dependencies, nil
}

func languageDependency(g *modfile.Go) depbot.Dependency {
	return depbot.Dependency{
		File:    GoDependencyFile,
		Version: g.Version,
		Name:    DependencyNameGo,
		Kind:    depbot.DependencyKindLanguage,
	}
}

func libraryDependency(r *modfile.Require) depbot.Dependency {
	return depbot.Dependency{
		File:    GoDependencyFile,
		Name:    r.Mod.Path,
		Version: r.Mod.Version,
		Kind:    depbot.DependencyKindLibrary,
		Direct:  !r.Indirect,
	}
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
