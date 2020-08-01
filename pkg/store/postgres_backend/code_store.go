package postgres_backend

import "github.com/discmonkey/retext/pkg/store"

type CodeStore struct {
	db connection
}

func (c CodeStore) CreateCode(name string, ParentCodeID store.CodeID) (store.CodeID, error) {
	panic("implement me")
}

func (c CodeStore) CodifyText(codeID store.CodeID, documentID store.FileID, text string, firstWord store.WordCoordinate, lastWord store.WordCoordinate) error {
	panic("implement me")
}

func (c CodeStore) GetCode(codeID store.CodeID) (store.Code, error) {
	panic("implement me")
}

func (c CodeStore) GetCodeContainer(codeID store.CodeID) (store.CodeContainer, error) {
	panic("implement me")
}

func (c CodeStore) Codes() ([]store.CodeID, error) {
	panic("implement me")
}

var _ store.CodeStore = CodeStore{}
