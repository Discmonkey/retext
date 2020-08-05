package postgres_backend

import (
	"github.com/discmonkey/retext/pkg/store"
	"testing"
)

func TestFileStorePostgresBackend(t *testing.T) {
	testDirName := store.CreateTestDir()

	fileBackend, err := NewFileStore(testDirName)

	if err != nil {
		t.Fatal(err)
	}

	store.StubTestStore(t, fileBackend, testDirName)
}

func TestCodeStoreFileBackend(t *testing.T) {

	codeBackend, err := NewCodeStore()
	if err != nil {
		t.Fatalf(err.Error())
	}

	store.StubTestCodeStore(t, codeBackend)

}
