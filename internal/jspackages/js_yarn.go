package jspackages

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/godepbot/depbot"
	"github.com/godepbot/depbot/internal"
)

func FindYarnDependencies(wd string) (depbot.Dependencies, error) {
	pths := internal.PathsFor(wd, jsYarnLockFile)

	dependencies := depbot.Dependencies{}

	depRegex, err := regexp.Compile(`(?:"?([^\s]+?)@)`)
	if err != nil {
		return dependencies, fmt.Errorf("error compiling regexp: %w", err)
	}

	versionRegex, err := regexp.Compile(`"(\d.+?)"`)
	if err != nil {
		return dependencies, fmt.Errorf("error compiling regexp: %w", err)
	}

	for _, p := range pths {
		relPath, _ := filepath.Rel(wd, p)
		if relPath == "" {
			relPath = jsYarnLockFile
		}

		openFile, err := ioutil.ReadFile(p)
		if err != nil {
			return dependencies, fmt.Errorf("error reading dependency file '%v': %w", p, err)
		}

		dependencies = append(dependencies, depbot.Dependency{
			File:     relPath,
			Kind:     depbot.DependencyKindLanguage,
			Language: depbot.DependencyLanguageJs,
			Name:     jsDependencyNameJs,
		})

		rawFile := strings.Split(string(openFile), "\n")
		version := ""
		name := ""

		for _, line := range rawFile {
			if strings.Contains(line, "lockfile") {
				sLine := strings.Split(line, " ")
				dependencies = append(dependencies, depbot.Dependency{
					File:     relPath,
					Kind:     depbot.DependencyKindTool,
					Language: depbot.DependencyLanguageJs,
					Name:     jsDependencyNameYARN,
					Version:  sLine[len(sLine)-1],
				})
				continue
			}

			result := depRegex.FindStringSubmatch(line)
			// position 0 is the whole match for the regexp
			// position 1 is the clean string
			if len(result) >= 2 {
				name = result[1]
			}

			if strings.Contains(line, "version") {
				version = versionRegex.FindString(line)
			}

			if version != "" && name != "" {
				dependencies = append(dependencies, depbot.Dependency{
					File:     relPath,
					Kind:     depbot.DependencyKindLibrary,
					Language: depbot.DependencyLanguageJs,
					Name:     strings.ReplaceAll(name, "\"", ""),
					Version:  strings.ReplaceAll(version, "\"", ""),
				})
				name = ""
				version = ""
			}
		}
	}

	return dependencies, nil
}
