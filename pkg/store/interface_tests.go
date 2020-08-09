package store

import (
	"io/ioutil"
	"os"
	"testing"
)

func CreateTestDir() string {
	testDirName, _ := ioutil.TempDir("", "retext")

	return testDirName
}

// TestFSBackend covers all the file interface methods
func StubTestStore(t *testing.T, fileBackend FileStore, testDirName string) {

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

	initialLength := len(files)

	contents := []byte("hello")
	testFileName := "test1.txt"
	file, err := fileBackend.UploadFile(testFileName, contents)
	if err != nil {
		t.Fatal(err)
	}

	files, err = fileBackend.Files()
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != initialLength+1 {
		t.Fatal("incorrect number of files returned")
	}

	f, err := fileBackend.GetFile(file.ID)
	if err != nil {
		t.Fatal(err)
	}

	for num, b := range f {
		if b != contents[num] {
			t.Fatal("file contents do not match")
		}
	}
}

func StubTestCodeStore(t *testing.T, codeBackend CodeStore, fileBackend FileStore) {

	testCodeName := "test"
	someBytes := []byte("hello")
	testFile, err := fileBackend.UploadFile("temp", someBytes)
	if err != nil {
		t.Fatalf("failed to upload file")
	}

	initial, err := codeBackend.GetContainers()
	if err != nil {
		t.Fatalf("could not query for containers")
	}

	initialLength := len(initial)

	containerID, err := codeBackend.CreateContainer()
	if err != nil {
		t.Fatalf(err.Error())
	}

	firstCodeID, err := codeBackend.CreateCode(testCodeName, containerID)
	if err != nil {
		t.Fatalf("failed to save code: %s", err)
	}

	firstCodeMain, err := codeBackend.GetContainer(containerID)
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
	_, err = codeBackend.CreateCode(testSubCodeName, containerID)
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

	err = codeBackend.CodifyText(firstCodeID, testFile.ID, testText, anchor, lastWord)
	if err != nil {
		t.Fatalf("failed to codify text: %s", err)
	}
	firstCode, err := codeBackend.GetCode(firstCodeID)
	if err != nil || len(firstCode.Texts) == 0 {
		t.Fatalf("failed to codify text: %s", err)
	}

	containers, err := codeBackend.GetContainers()
	if err != nil {
		t.Fatalf("failed to get list of codes: %s", err)
	}
	//TODO: update the # used in this len() comparison if you change the number
	// of created codes
	if len(containers)-initialLength != 1 {
		numCodes := len(containers)
		t.Fatalf("incorrect number of codes; got: %d", numCodes)
	}
	_ = os.Remove("/tmp/filetest")

	// second start-up tests "cache path"

}
