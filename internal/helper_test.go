package internal_test

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/godepbot/depbot/internal"
)

func Test_PathsFor_HelperMethod_Ignore_NodeModule(t *testing.T) {

	tmpDir := t.TempDir()

	directoryPaths := []string{
		tmpDir + "/package/",
		tmpDir + "/package/v2/",
		tmpDir + "/package/node_modules/",
	}

	filesNames := []string{
		"yarn.lock",
		"package-lock.json",
		"package.json",
	}

	for _, path := range directoryPaths {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			t.Fatalf("got an error but should be nil, error : %v ", err)
			return
		}
	}

	for index, fileName := range filesNames {
		errWriteFile := ioutil.WriteFile(directoryPaths[index]+fileName, []byte("hello depbot"), 0644)
		if errWriteFile != nil {
			t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		}
	}

	paths := internal.PathsFor(tmpDir, filesNames...)
	for _, path := range paths {

		if strings.Contains(path, "node_modules") {
			t.Fatalf("expected path to no contains 'node_modules' folder, but got %v", path)
		}
	}

}
