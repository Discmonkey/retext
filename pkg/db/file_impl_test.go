package db

import (
	"io/ioutil"
	"os"
	"testing"
)

// TestFSBackend covers all the file interface methods
func TestFileStore(t *testing.T) {
	testDirName, _ := ioutil.TempDir("", "retext")

	_ = os.RemoveAll(testDirName)

	fileBackend := &FSBackendFile{}

	err := fileBackend.Init(testDirName)
	if err != nil {
		t.Fatal(err)
	}

	if info, err := os.Stat(testDirName); err != nil || !info.IsDir() {
		if err != nil {
			t.Fatal(err)
		} else {
			t.Fatal("directory not created properly")
		}
	}

	files, err := fileBackend.Files()
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 0 {
		t.Fatal("incorrect number of files returned")
	}

	contents := []byte("hello")
	testFileName := "test1.txt"
	key, err := fileBackend.UploadFile(testFileName, contents)
	if err != nil {
		t.Fatal(err)
	}

	files, err = fileBackend.Files()
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 1 {
		t.Fatal("incorrect number of files returned")
	}

	if key != files[0] {
		t.Fatal("key does not match files scan")
	}

	f, err := fileBackend.GetFile(key)
	if err != nil {
		t.Fatal(err)
	}

	for num, b := range f {
		if b != contents[num] {
			t.Fatal("file contents do not match")
		}
	}
}
