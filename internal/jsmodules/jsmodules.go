package jsmodules

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/godepbot/depbot"
)

const (
	jsDependencyFile string = "package.json"
	dependencyNameJs string = "Js"
)

type PackageJson struct {
	License         string            `json:"license" db:"-"`
	Main            string            `json:"main" db:"-"`
	Name            string            `json:"name" db:"-"`
	Repository      string            `json:"repository" db:"-"`
	Scripts         map[string]string `json:"scripts" db:"-"`
	Version         string            `json:"version" db:"-"`
	Dependencies    map[string]string `json:"dependencies" db:"-"`
	DevDependencies map[string]string `json:"devDependencies" db:"-"`
}

func FindDependencies(wd string) (depbot.Dependencies, error) {
	pths := []string{}

	filepath.WalkDir(wd, func(path string, d fs.DirEntry, err error) error {
		if strings.Contains(path, jsDependencyFile) {
			pths = append(pths, path)
		}

		return nil
	})

	dependencies := depbot.Dependencies{}
	for _, p := range pths {
		openFile, err := ioutil.ReadFile(p)
		if err != nil {
			fmt.Println("Error reading the file.")
			return dependencies, err
		}

		packageJson := PackageJson{}
		errU := json.Unmarshal(openFile, &packageJson)
		if errU != nil {
			fmt.Println("Error mashal the file.", errU)
			return dependencies, errU
		}

		dependencies = append(dependencies, depbot.Dependency{
			File:    jsDependencyFile,
			Version: packageJson.Version,
			Name:    dependencyNameJs,
			Kind:    depbot.DependencyKindLanguage,
		})

		for d := range packageJson.Dependencies {
			dependencies = append(dependencies, depbot.Dependency{
				File:    jsDependencyFile,
				Name:    d,
				Version: packageJson.Dependencies[d],
				Kind:    depbot.DependencyKindLibrary,
			})
		}

		for d := range packageJson.DevDependencies {
			dependencies = append(dependencies, depbot.Dependency{
				File:    jsDependencyFile,
				Name:    d,
				Version: packageJson.DevDependencies[d],
				Kind:    depbot.DependencyKindLibrary,
			})
		}
	}

	return dependencies, nil
}
