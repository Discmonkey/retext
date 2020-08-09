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
	ParentMap store.CodeParentIDMap
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
			ParentMap: store.CodeParentIDMap{},
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

	return store.File{ID: -1, Type: store.SourceFile, Name: filename}, nil
}

func (F *DevFileBackend) GetFile(id store.FileID) ([]byte, error) {
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
			ID:   i,
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

func (F *DevCodeBackend) CreateCode(name string, parentCodeID store.CodeID) (store.CodeID, error) {
	// todo: code names should be unique (per project, probably?)
	F.codeMutex.Lock()
	defer F.codeMutex.Unlock()

	newID := len(F.cache.Flat) + 1

	newCode := store.Code{Name: name, ID: newID, Texts: []store.DocumentText{}}

	if parentCodeID == 0 {
		F.cache.ParentMap[newID] = []store.CodeID{newID}
	} else {
		if subCodes, ok := F.cache.ParentMap[parentCodeID]; ok {
			F.cache.ParentMap[parentCodeID] = append(subCodes, newID)
		} else {
			return 0, errors.New(fmt.Sprintf("CodeID not found. ID: %d", parentCodeID))
		}
	}

	F.cache.Flat[newID] = newCode

	err := F.writeCodesToFile(F.cache)
	if err != nil {
		return 0, err
	}

	return newCode.ID, nil
}

func (F *DevCodeBackend) CodifyText(codeID store.CodeID, documentID store.FileID, text string, firstWord store.WordCoordinate, lastWord store.WordCoordinate) error {
	F.codeMutex.Lock()
	defer F.codeMutex.Unlock()

	cache, err := F.getCodesFromFile()
	if err != nil {
		return err
	}

	if _, ok := cache.Flat[codeID]; ok == false {
		return errors.New(fmt.Sprintf("No code found with ID: %d", codeID))
	}

	var code = cache.Flat[codeID]
	code.Texts = append(code.Texts, store.DocumentText{
		DocumentID: documentID,
		Text:       text,
		FirstWord:  firstWord,
		LastWord:   lastWord,
	})
	cache.Flat[codeID] = code

	err = F.writeCodesToFile(cache)

	F.cache = cache
	return err
}

func (F *DevCodeBackend) GetCode(codeID store.CodeID) (store.Code, error) {
	F.codeMutex.RLock()
	defer F.codeMutex.RUnlock()

	if code, ok := F.cache.Flat[codeID]; ok {
		return code, nil
	}
	return store.Code{}, errors.New(fmt.Sprintf("No code found with ID: %d", codeID))
}

func (F *DevCodeBackend) GetCodeContainer(codeID store.CodeID) (store.CodeContainer, error) {
	F.codeMutex.RLock()
	defer F.codeMutex.RUnlock()

	if code, ok := F.cache.Flat[codeID]; ok {
		codeContainer := store.CodeContainer{
			Main:  code.ID,
			Codes: make([]store.Code, len(F.cache.ParentMap[code.ID])),
		}
		for i, subCodeID := range F.cache.ParentMap[code.ID] {
			codeContainer.Codes[i] = F.cache.Flat[subCodeID]
		}
		return codeContainer, nil
	}

	return store.CodeContainer{}, errors.New(fmt.Sprintf("No code found with ID: %d", codeID))
}

func (F *DevCodeBackend) Codes() ([]store.CodeID, error) {
	F.codeMutex.RLock()
	defer F.codeMutex.RUnlock()
	cache, err := F.getCodesFromFile()

	if err != nil {
		return nil, err
	}

	codeList := make([]store.CodeID, len(cache.ParentMap))
	i := 0
	for mainID := range cache.ParentMap {
		codeList[i] = mainID
		i++
	}

	return codeList, nil
}
