package jspackages

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/godepbot/depbot"
	"github.com/godepbot/depbot/internal"
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

func FindPackageDependencies(wd string) (depbot.Dependencies, error) {

	allPths := internal.PathsFor(wd, []string{jsPackageLockFile, jsPackageFile}...)
	dirTree := map[string][]string{}

	for _, path := range allPths {
		dir := filepath.Dir(path)
		dirTree[dir] = append(dirTree[dir], path)
	}

	pths := []string{}
	for _, files := range dirTree {
		existLockFileInTree := false
		packageFiles := []string{}
		for _, file := range files {
			if filepath.Base(file) == jsPackageFile {
				packageFiles = append(packageFiles, file)
			}
			if filepath.Base(file) == jsPackageLockFile {
				existLockFileInTree = true
			}
		}

		if !existLockFileInTree {
			pths = append(pths, packageFiles...)
		}
	}

	dependencies := depbot.Dependencies{}
	for _, p := range pths {
		relPath, _ := filepath.Rel(wd, p)
		if relPath == "" {
			relPath = jsPackageFile
		}

		openFile, err := ioutil.ReadFile(p)
		if err != nil {
			return dependencies, fmt.Errorf("error reading dependency file '%v': %w", p, err)
		}

		packageJson := PackageJson{}
		errU := json.Unmarshal(openFile, &packageJson)
		if errU != nil {
			return dependencies, fmt.Errorf("error parsing dependency file '%v': %w", p, errU)
		}
		dependencies = append(dependencies, packageDependencies(packageJson, relPath)...)
	}

	return dependencies, nil
}

func packageDependencies(p PackageJson, file string) depbot.Dependencies {
	dependencies := depbot.Dependencies{
		{
			File:     file,
			Kind:     depbot.DependencyKindTool,
			Language: depbot.DependencyLanguageJs,
			Name:     jsDependencyNameNPM,
		},
		{
			File:     file,
			Kind:     depbot.DependencyKindLanguage,
			Language: depbot.DependencyLanguageJs,
			Name:     jsDependencyNameJs,
			Version:  p.Version,
		},
	}

	for d := range p.Dependencies {
		dependencies = append(dependencies, depbot.Dependency{
			Direct:   true,
			Kind:     depbot.DependencyKindLibrary,
			Language: depbot.DependencyLanguageJs,
			File:     file,
			Name:     d,
			Version:  p.Dependencies[d],
		})
	}

	for d := range p.DevDependencies {
		dependencies = append(dependencies, depbot.Dependency{
			File:     file,
			Kind:     depbot.DependencyKindLibrary,
			Language: depbot.DependencyLanguageJs,
			Name:     d,
			Version:  p.DevDependencies[d],
		})
	}

	return dependencies
}
