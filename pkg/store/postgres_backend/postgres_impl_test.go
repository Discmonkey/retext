package postgres_backend

import (
	"fmt"
	"github.com/discmonkey/retext/pkg/store"
	"testing"
	"time"
)

func TestFileStorePostgresBackend(t *testing.T) {
	projects, files, _, _, err := setup()

	if err != nil {
		t.Fatal(err)
	}

	projectId, err := createTestProject(projects)

	if err != nil {
		t.Fatal(err)
	}

	store.StubTestStore(t, files, projectId)
}

func setup() (ProjectStore, FileStore, CodeStore, InsightStore, error) {
	conn, err := GetConnection()
	if err != nil {
		return ProjectStore{}, FileStore{}, CodeStore{}, InsightStore{}, err
	}

	testDirName := store.CreateTestDir()
	codeBackend := NewCodeStore(conn)
	fileBackend := NewFileStore(testDirName, conn)
	projectBackend := NewProjectStore(conn)
	insightsBackend := NewInsightStore(conn)

	return projectBackend, fileBackend, codeBackend, insightsBackend, nil
}

func TestCodeStorePostgresBackend(t *testing.T) {

	projects, files, codes, insights, err := setup()
	fatalIf(err, t)

	projectId, err := createTestProject(projects)
	fatalIf(err, t)

	store.StubTestCodeStore(t, codes, files, insights, projectId)

}

func fatalIf(err error, t *testing.T) {
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestCodeStoreList(t *testing.T) {
	projects, files, codes, _, err := setup()
	fatalIf(err, t)

	projectId, err := createTestProject(projects)
	fatalIf(err, t)

	_, err = files.UploadFile("test.txt", []byte("TestCodeStoreList"), projectId, store.SourceFile)
	fatalIf(err, t)

	containerId, err := codes.CreateContainer(projectId)
	fatalIf(err, t)

	containerId2, err := codes.CreateContainer(projectId)
	fatalIf(err, t)

	assertContainerID := func(codeId store.CodeId, containerId int64) {
		code, err := codes.GetCode(codeId)
		if err != nil {
			t.Fail()
			return
		}

		if code.Container != containerId {
			t.Fail()
		}
	}
	id1, err := codes.CreateCode("test1", containerId)

	fatalIf(err, t)
	assertContainerID(id1, containerId)

	id2, err := codes.CreateCode("test2", containerId2)
	fatalIf(err, t)
	assertContainerID(id2, containerId2)

	id3, err := codes.CreateCode("test3", containerId)
	fatalIf(err, t)
	assertContainerID(id3, containerId)

	containers, err := codes.GetContainers(projectId)

	if len(containers) < 2 {
		t.Fail()
	}

	for _, container := range containers {
		fmt.Println(container)
	}
}

func createTestProject(projectStore store.ProjectStore) (store.ProjectId, error) {
	return projectStore.CreateProject(fmt.Sprint("test", time.Now()), "test",
		int(time.Now().Month()), time.Now().Year())
}
