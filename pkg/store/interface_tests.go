package store

import (
	"io/ioutil"
	"log"
	"testing"
)

func CreateTestDir() string {
	testDirName, _ := ioutil.TempDir("", "retext")

	return testDirName
}

// TestFSBackend covers all the file interface methods
func StubTestStore(t *testing.T, fileBackend FileStore, projectId ProjectId) {
	files, err := fileBackend.GetFiles(projectId)
	if err != nil {
		t.Fatal(err)
	}

	initialLength := len(files)

	contents := []byte("hello")
	testFileName := "test1.txt"
	file, err := fileBackend.UploadFile(testFileName, contents, projectId, SourceFile)
	if err != nil {
		t.Fatal(err)
	}

	files, err = fileBackend.GetFiles(projectId)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != initialLength+1 {
		t.Fatal("incorrect number of files returned")
	}

	f, _, err := fileBackend.GetFile(file.Id)
	if err != nil {
		t.Fatal(err)
	}

	for num, b := range f {
		if b != contents[num] {
			t.Fatal("file contents do not match")
		}
	}
}

func StubTestCodeStore(t *testing.T, codeBackend CodeStore, fileBackend FileStore, insightsBackend InsightStore, projectId ProjectId) {

	testCodeName := "test"
	someBytes := []byte("hello")
	testFile, err := fileBackend.UploadFile("temp.txt", someBytes, projectId, SourceFile)
	if err != nil {
		log.Println(err)
		t.Fatalf("failed to upload file")
	}

	initial, err := codeBackend.GetContainers(projectId)
	if err != nil {
		t.Fatalf("could not query for containers")
	}

	initialLength := len(initial)

	containerId, err := codeBackend.CreateContainer(projectId)
	if err != nil {
		t.Fatalf(err.Error())
	}

	firstCodeId, err := codeBackend.CreateCode(testCodeName, containerId)
	if err != nil {
		t.Fatalf("failed to save code: %s", err)
	}

	firstCodeMain, err := codeBackend.GetContainer(containerId)
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
	_, err = codeBackend.CreateCode(testSubCodeName, containerId)
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

	_, err = codeBackend.CodifyText(firstCodeId, testFile.Id, testText, anchor, lastWord)
	if err != nil {
		t.Fatalf("failed to codify text: %s", err)
	}
	firstCode, err := codeBackend.GetCode(firstCodeId)
	if err != nil || len(firstCode.Texts) == 0 {
		t.Fatalf("failed to codify text: %s", err)
	}

	containers, err := codeBackend.GetContainers(projectId)
	if err != nil {
		t.Fatalf("failed to get list of codes: %s", err)
	}
	//TODO: update the # used in this len() comparison if you change the number
	// of created codes
	if len(containers)-initialLength != 1 {
		numCodes := len(containers)
		t.Fatalf("incorrect number of codes; got: %d", numCodes)
	}

	err = codeBackend.DeleteText(firstCode.Texts[0].Id)

	if err != nil {
		t.Fatalf("failed to uncode text: %s", err)
	}

	// insights test
	_, _ = codeBackend.CodifyText(firstCodeId, testFile.Id, testText, anchor, lastWord)
	_, _ = codeBackend.CodifyText(firstCodeId, testFile.Id, testText, anchor, lastWord)
	insightCode, _ := codeBackend.GetCode(firstCodeId)
	textIds := make([]TextId, 0)
	for _, text := range insightCode.Texts {
		textIds = append(textIds, text.Id)
	}

	_, err = insightsBackend.CreateInsight(projectId, "test", textIds)
	if err != nil {
		t.Fatalf("failed to create insight: %s", err)
	}

	insights, err := insightsBackend.GetInsights(projectId)

	if err != nil {
		t.Fatalf("failed to get insights list: %s", err)
	} else if len(insights) != 1 {
		t.Fatalf("failed to get the right number of insights")
	}
}
