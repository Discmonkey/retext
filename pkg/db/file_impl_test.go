package db

import (
	"io/ioutil"
	"os"
	"testing"
)

// TestFSBackend covers all the interface methods
func TestFSBackend(t *testing.T) {
	testDirName, _ := ioutil.TempDir("", "retext")

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

	testCategoryName := "test"
	firstCategoryID, err := store.CreateCategory(testCategoryName, 0)
	if err != nil {
		t.Fatalf("failed to save category: %s", err)
	}

	firstCategoryMain, err := store.GetCategoryMain(firstCategoryID)
	if err != nil {
		t.Fatalf("failed to get category: %s", err)
	}
	if firstCategoryMain.Categories[0].Name != testCategoryName {
		t.Fatalf("category came back with unexpected name: %s", err)
	}
	_, err = store.GetCategory(1000)
	if err == nil {
		t.Fatal("non-existent categories should return an error")
	}
	// test creating a subcategory
	testSubCatName := "subcat 1 1"
	_, err = store.CreateCategory(testSubCatName, firstCategoryMain.Main)
	if err != nil {
		t.Fatalf("unable to create a subcategory: %s", err)
	}

	testText := "made up text"
	anchor := WordCoordinate{
		Paragraph: 1,
		Sentence:  1,
		Word:      1,
	}
	lastWord := WordCoordinate{
		Paragraph: 1,
		Sentence:  1,
		Word:      3,
	}
	err = store.CategorizeText(firstCategoryID, testFileName, testText, anchor, lastWord)
	if err != nil {
		t.Fatalf("failed to categorize text: %s", err)
	}
	firstCategory, err := store.GetCategory(firstCategoryID)
	if err != nil || len(firstCategory.Texts) == 0 {
		t.Fatalf("failed to categorize text: %s", err)
	}

	cats, err := store.Categories()
	if err != nil {
		t.Fatalf("failed to get list of categories: %s", err)
	}
	//TODO: update the # used in this len() comparison if you change the number
	// of created categories
	if len(cats) != 1 {
		numCats := len(cats)
		t.Fatalf("incorrect number of categories; got: %d", numCats)
	}
	_ = os.Remove("/tmp/filetest")

	// second start-up tests "cache path"
	err = store.Init(testDirName)
	if err != nil {
		t.Fatalf("failed to load cached categories: %s", err)
	}
}
