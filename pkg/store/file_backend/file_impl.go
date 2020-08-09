package file_backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/discmonkey/retext/pkg/store"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
)

type CodeCache struct {
	Flat      store.CodeMap
	ParentMap store.CodeParentIdMap
}
type DevFileBackend struct {
	dirLocation string
}
type DevCodeBackend struct {
	codeMutex   sync.RWMutex
	codeFileLoc string
	cache       CodeCache
}

func (F *DevFileBackend) Init(pathToDir string) error {
	if _, err := os.Stat(pathToDir); os.IsNotExist(err) {
		err = os.Mkdir(pathToDir, 0766)

		if err != nil {
			return err
		}
	}

	uploadDirLocation := path.Join(pathToDir, "uploadLocation")
	if _, err := os.Stat(uploadDirLocation); os.IsNotExist(err) {
		err = os.Mkdir(uploadDirLocation, 0766)

		if err != nil {
			return err
		}
	}

	F.dirLocation = uploadDirLocation

	return nil
}
func (F *DevCodeBackend) Init(pathToDir string) error {
	if _, err := os.Stat(pathToDir); os.IsNotExist(err) {
		err = os.Mkdir(pathToDir, 0766)

		if err != nil {
			return err
		}
	}

	//create file to store codes, if it doesn't already exist
	F.codeFileLoc = path.Join(pathToDir, "codes.json")
	if _, err := os.Stat(F.codeFileLoc); os.IsNotExist(err) {
		F.cache = CodeCache{
			Flat:      store.CodeMap{},
			ParentMap: store.CodeParentIdMap{},
		}
		err = jsonToFile(F.codeFileLoc, F.cache)
		if err != nil {
			return err
		}

	} else {
		F.cache, err = F.getCodesFromFile()

		if err != nil {
			return err
		}
	}

	return nil
}

func (F *DevFileBackend) UploadFile(filename string, contents []byte) (store.File, error) {
	filepath := path.Join(F.dirLocation, filename)
	err := ioutil.WriteFile(filepath, contents, 0644)

	if err != nil {
		return store.File{}, err
	}

	return store.File{Id: -1, Type: store.SourceFile, Name: filename}, nil
}

func (F *DevFileBackend) GetFile(id store.FileId) ([]byte, error) {
	// shit's hacky ha
	filepath := path.Join(F.dirLocation, fmt.Sprintf("%d", id))

	return ioutil.ReadFile(filepath)
}

func (F *DevFileBackend) Files() ([]store.File, error) {
	osFiles, err := ioutil.ReadDir(F.dirLocation)
	if err != nil {
		return nil, err
	}

	files := make([]store.File, len(osFiles))

	for i, f := range osFiles {
		fType := store.SourceFile
		if strings.HasSuffix(f.Name(), "xlsx") {
			fType = store.DemoFile
		}

		files[i] = store.File{
			Id:   i,
			Type: fType,
		}
	}

	return files, nil
}

func jsonFromFile(filename string, i interface{}) (err error) {
	jsonFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonFile, &i)

	return err
}
func jsonToFile(filename string, i interface{}) error {
	newJson, err := json.Marshal(i)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, newJson, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (F *DevCodeBackend) getCodesFromFile() (CodeCache, error) {
	var cache CodeCache
	err := jsonFromFile(F.codeFileLoc, &cache)
	return cache, err
}

func (F *DevCodeBackend) writeCodesToFile(cache CodeCache) error {
	err := jsonToFile(F.codeFileLoc, cache)
	return err
}

func (F *DevCodeBackend) CreateCode(name string, parentCodeId store.CodeId) (store.CodeId, error) {
	// todo: code names should be unique (per project, probably?)
	F.codeMutex.Lock()
	defer F.codeMutex.Unlock()

	newId := len(F.cache.Flat) + 1

	newCode := store.Code{Name: name, Id: newId, Texts: []store.DocumentText{}}

	if parentCodeId == 0 {
		F.cache.ParentMap[newId] = []store.CodeId{newId}
	} else {
		if subCodes, ok := F.cache.ParentMap[parentCodeId]; ok {
			F.cache.ParentMap[parentCodeId] = append(subCodes, newId)
		} else {
			return 0, errors.New(fmt.Sprintf("CodeId not found.Id: %d", parentCodeId))
		}
	}

	F.cache.Flat[newId] = newCode

	err := F.writeCodesToFile(F.cache)
	if err != nil {
		return 0, err
	}

	return newCode.Id, nil
}

func (F *DevCodeBackend) CodifyText(codeId store.CodeId, documentId store.FileId, text string, firstWord store.WordCoordinate, lastWord store.WordCoordinate) error {
	F.codeMutex.Lock()
	defer F.codeMutex.Unlock()

	cache, err := F.getCodesFromFile()
	if err != nil {
		return err
	}

	if _, ok := cache.Flat[codeId]; ok == false {
		return errors.New(fmt.Sprintf("No code found withId: %d", codeId))
	}

	var code = cache.Flat[codeId]
	code.Texts = append(code.Texts, store.DocumentText{
		DocumentId: documentId,
		Text:       text,
		FirstWord:  firstWord,
		LastWord:   lastWord,
	})
	cache.Flat[codeId] = code

	err = F.writeCodesToFile(cache)

	F.cache = cache
	return err
}

func (F *DevCodeBackend) GetCode(codeId store.CodeId) (store.Code, error) {
	F.codeMutex.RLock()
	defer F.codeMutex.RUnlock()

	if code, ok := F.cache.Flat[codeId]; ok {
		return code, nil
	}
	return store.Code{}, errors.New(fmt.Sprintf("No code found withId: %d", codeId))
}

func (F *DevCodeBackend) GetCodeContainer(codeId store.CodeId) (store.CodeContainer, error) {
	F.codeMutex.RLock()
	defer F.codeMutex.RUnlock()

	if code, ok := F.cache.Flat[codeId]; ok {
		codeContainer := store.CodeContainer{
			Main:  code.Id,
			Codes: make([]store.Code, len(F.cache.ParentMap[code.Id])),
		}
		for i, subCodeId := range F.cache.ParentMap[code.Id] {
			codeContainer.Codes[i] = F.cache.Flat[subCodeId]
		}
		return codeContainer, nil
	}

	return store.CodeContainer{}, errors.New(fmt.Sprintf("No code found withId: %d", codeId))
}

func (F *DevCodeBackend) Codes() ([]store.CodeId, error) {
	F.codeMutex.RLock()
	defer F.codeMutex.RUnlock()
	cache, err := F.getCodesFromFile()

	if err != nil {
		return nil, err
	}

	codeList := make([]store.CodeId, len(cache.ParentMap))
	i := 0
	for mainId := range cache.ParentMap {
		codeList[i] = mainId
		i++
	}

	return codeList, nil
}
