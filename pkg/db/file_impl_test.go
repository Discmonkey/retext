package db

import (
	"io/ioutil"
	"os"
	"testing"
)

func createTestDir() string {
	testDirName, _ := ioutil.TempDir("", "retext")
	return testDirName
}

// TestFSBackend covers all the file interface methods
func TestFileStore(t *testing.T) {
	testDirName := createTestDir()

	fileBackend := &DevFileBackend{}

	err := fileBackend.Init(testDirName)
	if err != nil {
		t.Fatal(err)
	}

	if info, err := os.Stat(testDirName + "/uploadLocation"); err != nil || !info.IsDir() {
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

	if key != files[0].ID {
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

	_ = os.Remove(testDirName)
}

func TestCodeStore(t *testing.T) {
	testDirName := createTestDir()

	codeBackend := &DevCodeBackend{}
	err := codeBackend.Init(testDirName)

	testCodeName := "test"
	firstCodeID, err := codeBackend.CreateCode(testCodeName, 0)
	if err != nil {
		t.Fatalf("failed to save code: %s", err)
	}

	firstCodeMain, err := codeBackend.GetCodeContainer(firstCodeID)
	if err != nil {
		t.Fatalf("failed to get code: %s", err)
	}
	if firstCodeMain.Codes[0].Name != testCodeName {
		t.Fatalf("code came back with unexpected name: %s", err)
	}
	_, err = codeBackend.GetCode(1000)
	if err == nil {
		t.Fatal("non-existent codes should return an error")
	}
	// test creating a subcode
	testSubCodeName := "subcode 1 1"
	_, err = codeBackend.CreateCode(testSubCodeName, firstCodeMain.Main)
	if err != nil {
		t.Fatalf("unable to create a subcode: %s", err)
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
	testFileName := "test1.txt"
	err = codeBackend.CodifyText(firstCodeID, testFileName, testText, anchor, lastWord)
	if err != nil {
		t.Fatalf("failed to codify text: %s", err)
	}
	firstCode, err := codeBackend.GetCode(firstCodeID)
	if err != nil || len(firstCode.Texts) == 0 {
		t.Fatalf("failed to codify text: %s", err)
	}

	codes, err := codeBackend.Codes()
	if err != nil {
		t.Fatalf("failed to get list of codes: %s", err)
	}
	//TODO: update the # used in this len() comparison if you change the number
	// of created codes
	if len(codes) != 1 {
		numCodes := len(codes)
		t.Fatalf("incorrect number of codes; got: %d", numCodes)
	}
	_ = os.Remove("/tmp/filetest")

	// second start-up tests "cache path"
	err = codeBackend.Init(testDirName)
	if err != nil {
		t.Fatalf("failed to load cached codes: %s", err)
	}

	_ = os.Remove(testDirName)
}
