package postgres_backend

import (
	"github.com/discmonkey/retext/pkg/store"
	"testing"
)

func TestFileStorePostgresBackend(t *testing.T) {
	testDirName := store.CreateTestDir()

	fileBackend, err := NewStore(testDirName)

	if err != nil {
		t.Fatal(err)
	}

	store.StubTestStore(t, fileBackend, testDirName)
}
