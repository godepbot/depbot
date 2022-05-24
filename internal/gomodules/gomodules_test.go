package gomodules_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/godepbot/depbot/internal/gomodules"
)

var fileContent = `
	module github.com/godepbot/depbot

	go 1.18
	
	require (
		github.com/gobuffalo/plush/v4 v4.1.9
		golang.org/x/mod v0.5.1 // indirect
		golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898 // indirect
	)
`

func Test_SingleDependency(t *testing.T) {

	file, err := os.Create(t.TempDir() + "/" + "go.mod")

	if err != nil {
		t.Fatalf("got an error but should be nil, error: %v ", err)
		return
	}

	errWriteFile := ioutil.WriteFile(file.Name(), []byte(fileContent), 0644)
	if errWriteFile != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	analysis, err := gomodules.FindDependencies(file.Name())
	if err != nil {
		t.Fatalf("got an error but should be nil, error : %v ", err.Error())
		return
	}

	if len(analysis.Dependencies) != 4 {
		t.Fatalf("got %v, but was expected %v", len(analysis.Dependencies), 4)
		return
	}

	if analysis.Dependencies[0].Name != "Go" {
		t.Fatalf("Got %v, but was expected %v", analysis.Dependencies[0].Name, "github.com/gobuffalo/plush/v4")
	}

	if analysis.Dependencies[1].Name != "github.com/gobuffalo/plush/v4" {
		t.Fatalf("Got %v, but was expected %v", analysis.Dependencies[1].Name, "github.com/gobuffalo/plush/v4")
	}

	if analysis.Dependencies[3].Version != "v0.0.0-20191011141410-1b5146add898" {
		t.Fatalf("Got %v, but was expected %v", analysis.Dependencies[1].Version, "v0.0.0-20191011141410-1b5146add898")
	}

}

func Test_MultipleDependencies(t *testing.T) {
	tmpDir := t.TempDir()

	errWriteFile := ioutil.WriteFile(tmpDir+"/go.mod", []byte(fileContent), 0644)
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

	errWriteFile = ioutil.WriteFile(packagePath+"/go.mod", []byte(fileContent), 0644)
	if errWriteFile != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	errWriteFile = ioutil.WriteFile(newDirectoriesPath+"/go.mod", []byte(fileContent), 0644)
	if errWriteFile != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	analysis, errFindDep := gomodules.FindDependencies(tmpDir)
	if errFindDep != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	if len(analysis.Dependencies) != 12 {
		t.Fatalf("got %v, but was expected %v", len(analysis.Dependencies), 12)
		return
	}

}

func Test_NoDependency(t *testing.T) {
	tmp := t.TempDir()

	analysis, err := gomodules.FindDependencies(tmp)

	if err != nil {
		t.Fatalf("Error finding the dependencies : %v ", err.Error())
		return
	}

	if len(analysis.Dependencies) > 0 {
		t.Fatalf("got %v, but was expected %v", len(analysis.Dependencies), 0)
	}

}
