package jsmodules_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/godepbot/depbot/internal/jsmodules"
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

func Test_JsSingleDependency(t *testing.T) {

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

	dependencies, err := jsmodules.FindDependencies(file.Name())
	if err != nil {
		t.Fatalf("got an error but should be nil, error : %v ", err.Error())
		return
	}

	if len(dependencies) != 5 {
		t.Fatalf("got %v, but was expected %v", len(dependencies), 5)
		return
	}

	if dependencies[0].Name != "Js" {
		t.Fatalf("Got %v, but was expected %v", dependencies[0].Name, "Js")
	}

	if dependencies[1].Name != "@fortawesome/fontawesome-free" {
		t.Fatalf("Got %v, but was expected %v", dependencies[1].Name, "@fortawesome/fontawesome-free")
	}

	if dependencies[3].Version != "^7.16.0" {
		t.Fatalf("Got %v, but was expected %v", dependencies[1].Version, "^7.16.0")
	}

}

func Test_JsMultipleDependencies(t *testing.T) {
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

	errWriteFile = ioutil.WriteFile(newDirectoriesPath+"/package.json", []byte(fileContent), 0644)
	if errWriteFile != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	dependencies, errFindDep := jsmodules.FindDependencies(tmpDir)
	if errFindDep != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	if len(dependencies) != 15 {
		t.Fatalf("got %v, but was expected %v", len(dependencies), 15)
		return
	}

}

func Test_JsNoDependency(t *testing.T) {
	tmp := t.TempDir()

	dependencies, err := jsmodules.FindDependencies(tmp)

	if err != nil {
		t.Fatalf("Error finding the dependencies : %v ", err.Error())
		return
	}

	if len(dependencies) > 0 {
		t.Fatalf("got %v, but was expected %v", len(dependencies), 0)
	}

}
