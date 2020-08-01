package file_backend

import (
	"github.com/discmonkey/retext/pkg/store"
	"os"
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

func TestCodeStoreFileBackend(t *testing.T) {
	testDirName := store.CreateTestDir()

	codeBackend := &DevCodeBackend{}
	err := codeBackend.Init(testDirName)

	store.StubTestCodeStore(t, codeBackend)

	err = codeBackend.Init(testDirName)
	if err != nil {
		t.Fatalf("failed to load cached codes: %s", err)
	}

	_ = os.Remove(testDirName)
}
