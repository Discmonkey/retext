package file_backend

import (
	"github.com/discmonkey/retext/pkg/store"
	"testing"
)

func TestFileStoreFileBackend(t *testing.T) {
	testDirName := store.CreateTestDir()
	fileBackend := DevFileBackend{}

	err := fileBackend.Init(testDirName)
	if err != nil {
		t.Fatal(err)
	}

	store.StubTestStore(t, &fileBackend, testDirName)
}
