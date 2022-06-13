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

	fileNames := []string{
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

	for index, fileName := range fileNames {
		errWriteFile := ioutil.WriteFile(directoryPaths[index]+fileName, []byte("hello depbot"), 0644)
		if errWriteFile != nil {
			t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		}
	}

	paths := internal.PathsFor(tmpDir, fileNames...)
	for _, path := range paths {

		if strings.Contains(path, "node_modules") {
			t.Fatalf("expected path to no contains 'node_modules' folder, but got %v", path)
		}
	}

}

func Test_PathContainsFolder_HelperMethod(t *testing.T) {

	tcases := []struct {
		path                string
		folderName          string
		shouldContainFolder bool
	}{
		{
			path:                "/depbot/node_module",
			folderName:          "node_module",
			shouldContainFolder: true,
		},
		{
			path:                "/depbot/nodes_mosdules",
			folderName:          "node_module",
			shouldContainFolder: false,
		},
		{
			path:                "/depbot/other/basis_path",
			folderName:          "other",
			shouldContainFolder: true,
		},
	}

	for _, tcase := range tcases {
		if internal.PathContainsFolder(tcase.path, tcase.folderName) != tcase.shouldContainFolder {
			t.Fatalf("expected condition to not be: %v", internal.PathContainsFolder(tcase.path, tcase.folderName))
		}
	}
}
