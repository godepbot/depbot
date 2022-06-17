package jspackages_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/godepbot/depbot"
	"github.com/godepbot/depbot/internal/jspackages"
)

var fileContent = `
	{
		"license": "MIT",
		"main": "index.js",
		"name": "buffalo",
		"repository": "repo",
		"scripts": {
			"dev": "webpack --watch",
			"build": "webpack --mode production --progress"
		},
		"version": "1.0.0",
		"dependencies": {
			"@fortawesome/fontawesome-free": "^5.15.4",
			"@hotwired/stimulus": "^3.0.1"
		},
		"devDependencies": {
			"@babel/cli": "^7.16.0",
			"webpack-manifest-plugin": "^4.0.2"
		}
	}  
`

var secondFileContent = `
	{
		"license": "MIT",
		"main": "index.js",
		"name": "buffalo",
		"repository": "repo",
		"scripts": {
			"dev": "webpack --watch",
			"build": "webpack --mode production --progress"
		},
		"version": "1.0.0",
		"dependencies": {
			"depbot": "^500",
			"milo": "^3.0.1"
		}
	}  
`

func Test_Package_SingleDependency(t *testing.T) {

	file, err := os.Create(t.TempDir() + "/" + "package.json")

	if err != nil {
		t.Fatalf("got an error but should be nil, error: %v ", err)
		return
	}

	errWriteFile := ioutil.WriteFile(file.Name(), []byte(fileContent), 0644)
	if errWriteFile != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	dependencies, err := jspackages.FindPackageDependencies(file.Name())
	if err != nil {
		t.Fatalf("got an error but should be nil, error : %v ", err.Error())
		return
	}

	if len(dependencies) != 6 {
		t.Fatalf("got %v, but was expected %v", len(dependencies), 6)
		return
	}

	for _, dependency := range dependencies {
		if dependency.Language != depbot.DependencyLanguageJs {
			t.Fatalf("got %v, but was expected %v", dependency.Language, depbot.DependencyLanguageJs)
		}
	}

	tcases := []struct {
		name    string
		version string
		exist   bool
	}{
		{
			name:    "@fortawesome/fontawesome-free",
			version: "^5.15.4",
			exist:   true,
		},
		{
			name:    "@hotwired/stimulus",
			version: "^3.0.1",
			exist:   true,
		},
		{
			name:    "cli",
			version: "^7.16.0",
			exist:   false,
		},
	}

	for _, tcase := range tcases {
		var exist bool
		for _, d := range dependencies {
			if d.Name == tcase.name && d.Version == tcase.version {
				exist = true
			}
		}

		if exist != tcase.exist {
			complement := "exist"
			if !tcase.exist {
				complement = "no exist"
			}
			t.Fatalf("expected %v with version %v to %v", tcase.name, tcase.version, complement)
		}
	}
}

func Test_Package_MultipleDependencies(t *testing.T) {
	tmpDir := t.TempDir()

	errWriteFile := ioutil.WriteFile(tmpDir+"/package.json", []byte(fileContent), 0644)
	if errWriteFile != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	newDirectoriesPath := tmpDir + "/package/v2"
	packagePath := tmpDir + "/package"

	err := os.MkdirAll(newDirectoriesPath, os.ModePerm)

	if err != nil {
		t.Fatalf("got an error but should be nil, error : %v ", err)
		return
	}

	errWriteFile = ioutil.WriteFile(packagePath+"/package.json", []byte(fileContent), 0644)
	if errWriteFile != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	errWriteFile = ioutil.WriteFile(newDirectoriesPath+"/package.json", []byte(secondFileContent), 0644)
	if errWriteFile != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	dependencies, errFindDep := jspackages.FindPackageDependencies(tmpDir)
	if errFindDep != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errFindDep.Error())
		return
	}

	tcases := []struct {
		name    string
		version string
		exist   bool
	}{
		{
			name:    "@fortawesome/fontawesome-free",
			version: "^5.15.4",
			exist:   true,
		},
		{
			name:    "@hotwired/stimulus",
			version: "^3.0.1",
			exist:   true,
		},
		{
			name:    "cli",
			version: "^7.16.0",
			exist:   false,
		},
		{
			name:    "depbot",
			version: "^7.16.0",
			exist:   false,
		},
		{
			name:    "depbot",
			version: "^500",
			exist:   true,
		},
	}

	for _, tcase := range tcases {
		var exist bool
		for _, d := range dependencies {
			if d.Name == tcase.name && d.Version == tcase.version {
				exist = true
			}
		}

		if exist != tcase.exist {
			complement := "exist"
			if !tcase.exist {
				complement = "no exist"
			}
			t.Fatalf("expected %v with version %v to %v", tcase.name, tcase.version, complement)
		}
	}

}

func Test_Files_With_Similar_Names(t *testing.T) {
	tmpDir := t.TempDir()

	errWriteFile := ioutil.WriteFile(tmpDir+"/package.json.Mac", []byte(fileContent), 0644)
	if errWriteFile != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	newDirectoriesPath := tmpDir + "/package/v2"
	packagePath := tmpDir + "/package"

	err := os.MkdirAll(newDirectoriesPath, os.ModePerm)

	if err != nil {
		t.Fatalf("got an error but should be nil, error : %v ", err)
		return
	}

	errWriteFile = ioutil.WriteFile(packagePath+"/Karo.package.json", []byte(fileContent), 0644)
	if errWriteFile != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	errWriteFile = ioutil.WriteFile(newDirectoriesPath+"/package.json", []byte(secondFileContent), 0644)
	if errWriteFile != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	dependencies, errFindDep := jspackages.FindPackageDependencies(tmpDir)
	if errFindDep != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errFindDep.Error())
		return
	}

	if len(dependencies) != 4 {
		t.Fatalf("got %v, but was expected %v", len(dependencies), 4)
		return
	}
}

func Test_Package_NoDependency(t *testing.T) {
	tmp := t.TempDir()

	dependencies, err := jspackages.FindPackageDependencies(tmp)

	if err != nil {
		t.Fatalf("error finding the dependencies : %v ", err.Error())
		return
	}

	if len(dependencies) > 0 {
		t.Fatalf("got %v, but was expected %v", len(dependencies), 0)
	}

}

func Test_Package_Dep_Relative_Path_Check(t *testing.T) {
	tmpDir := t.TempDir()

	packagePath := tmpDir + "/routecheck"

	err := os.MkdirAll(packagePath, os.ModePerm)
	if err != nil {
		t.Fatalf("got an error but should be nil, error : %v ", err)
		return
	}

	errWriteFile := ioutil.WriteFile(packagePath+"/package.json", []byte(fileContent), 0644)
	if errWriteFile != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	dependencies, errFindDep := jspackages.FindPackageDependencies(tmpDir)
	if errFindDep != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errFindDep.Error())
		return
	}

	for _, dependency := range dependencies {
		if dependency.File != "routecheck/package.json" {
			t.Fatalf("file expected to be in 'routecheck/package.json'")
		}
	}
}

func Test_Package_NoAnalize_If_Packagelock_Exist_At_Same_Folder_Level(t *testing.T) {
	tmpDir := t.TempDir()

	directoryPaths := []string{
		tmpDir + "/package/",
		tmpDir + "/package/v2/",
		tmpDir + "/package/v3/",
	}

	for _, path := range directoryPaths {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			t.Fatalf("got an error but should be nil, error : %v ", err)
			return
		}
	}

	files := []string{
		tmpDir + "/package/package.json",
		tmpDir + "/package/v2/package-lock.json",
		tmpDir + "/package/v2/package.json",
		tmpDir + "/package/v3/package.json",
		tmpDir + "/package/v3/package-lock.json",
	}

	for _, fileName := range files {
		errWriteFile := ioutil.WriteFile(fileName, []byte(secondFileContent), 0644)
		if errWriteFile != nil {
			t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		}
	}

	dependencies, errFindDep := jspackages.FindPackageDependencies(tmpDir)
	if errFindDep != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errFindDep.Error())
		return
	}

	if len(dependencies) != 4 {
		t.Fatalf("got %v, but was expected %v", len(dependencies), 0)
		return
	}

	for _, dependency := range dependencies {
		if dependency.File != "package/package.json" {
			t.Fatalf("got %v, but was expected %v", dependency.File, "package/package.json")
			return
		}
	}
}
