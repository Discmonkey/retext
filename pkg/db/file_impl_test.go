package db

import (
	"io/ioutil"
	"os"
	"testing"
)

// TestFSBackend covers all the interface methods
func TestFSBackend(t *testing.T) {
	testDirName, _ := ioutil.TempDir("", "filetest")

	_ = os.RemoveAll(testDirName)

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
	testFileName := "test1.txt"
	key, err := store.UploadFile(testFileName, contents)
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

	//todo: if your "categoriesFile" has categories added by using the website, the `len(cats) != 1` test will fail.
	testCategoryName := "test"
	c, err := store.CreateCategory(testCategoryName)
	if err != nil {
		t.Fatalf("failed to save category: %s", err)
	}

	c2, err := store.GetCategory(c)
	if err != nil {
		t.Fatalf("failed to get category: %s", err)
	}
	if c2.Name != testCategoryName {
		t.Fatalf("category came back with unexpected name: %s", err)
	}
	//todo: write a test that looks for a non-existent category

	testText := "made up text"
	err = store.CategorizeText(c, testFileName, testText)
	if err != nil {
		t.Fatalf("failed to categorize text: %s", err)
	}
	c2, _ = store.GetCategory(c)
	if len(c2.Texts) == 0 {
		t.Fatal("failed to categorize text")
	}

	cats, err := store.Categories()
	if err != nil {
		t.Fatalf("failed to get list of categories: %s", err)
	}
	//this check for 1 will break if you
	if len(cats) != 1 {
		numCats := string(len(cats))
		t.Fatalf("incorrect number of categories; got: %s", numCats)
	}
	_ = os.Remove("/tmp/filetest")

}