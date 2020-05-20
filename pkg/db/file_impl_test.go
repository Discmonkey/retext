package db

import (
	"os"
	"testing"
)

// TestFSBackend covers all the interface methods
func TestFSBackend(t *testing.T) {
	testDirName := "/tmp/filetest"

	_ = os.RemoveAll("/tmp/filetest")

	store := &FSBackend{}

	err := store.Init(testDirName)
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

	files, err := store.Files()
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 0 {
		t.Fatal("incorrect number of files returned")
	}

	contents := []byte("hello")
	key, err := store.UploadFile("test1.txt", contents)
	if err != nil {
		t.Fatal(err)
	}

	files, err = store.Files()
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 1 {
		t.Fatal("incorrect number of files returned")
	}

	if key != files[0] {
		t.Fatal("key does not match files scan")
	}

	f, err := store.GetFile(key)
	if err != nil {
		t.Fatal(err)
	}

	for num, b := range f {
		if b != contents[num] {
			t.Fatal("file contents do not match")
		}
	}

	_ = os.Remove("/tmp/filetest")

}
