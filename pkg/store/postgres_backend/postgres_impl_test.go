package postgres_backend

import (
	"fmt"
	"github.com/discmonkey/retext/pkg/store"
	"testing"
)

func TestFileStorePostgresBackend(t *testing.T) {
	testDirName := store.CreateTestDir()

	fileBackend, err := NewFileStore(testDirName)

	if err != nil {
		t.Fatal(err)
	}

	store.StubTestStore(t, fileBackend)
}

func setup() (*CodeStore, *FileStore, error) {
	codeBackend, err := NewCodeStore()
	if err != nil {
		return nil, nil, err
	}

	testDirName := store.CreateTestDir()
	fileBackend, err := NewFileStore(testDirName)
	if err != nil {
		return nil, nil, err
	}

	return codeBackend, fileBackend, nil
}

func TestCodeStorePostgresBackend(t *testing.T) {

	cStore, fStore, err := setup()

	if err != nil {
		t.Fatalf(err.Error())
	}

	store.StubTestCodeStore(t, cStore, fStore)

}

func fatalIf(err error, t *testing.T) {
	if err != nil {
		t.Fatalf(err.Error())
	}
}
func TestCodeStoreList(t *testing.T) {
	cStore, fStore, err := setup()
	fatalIf(err, t)

	_, err = fStore.UploadFile("test.txt", []byte("TestCodeStoreList"))
	fatalIf(err, t)

	containerId, err := cStore.CreateContainer()
	fatalIf(err, t)

	containerId2, err := cStore.CreateContainer()
	fatalIf(err, t)

	assertContainerID := func(codeId store.CodeId, containerId int) {
		code, err := cStore.GetCode(codeId)
		if err != nil {
			t.Fail()
			return
		}

		if code.Container != containerId {
			t.Fail()
		}
	}
	id1, err := cStore.CreateCode("test1", containerId)

	fatalIf(err, t)
	assertContainerID(id1, containerId)

	id2, err := cStore.CreateCode("test2", containerId2)
	fatalIf(err, t)
	assertContainerID(id2, containerId2)

	id3, err := cStore.CreateCode("test3", containerId)
	fatalIf(err, t)
	assertContainerID(id3, containerId)

	containers, err := cStore.GetContainers()

	if len(containers) < 2 {
		t.Fail()
	}

	for _, container := range containers {
		fmt.Println(container)
	}
}
