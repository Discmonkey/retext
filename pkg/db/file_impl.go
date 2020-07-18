package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
)

type FSCache struct {
	Flat      CategoryMap
	ParentMap CategoryParentIDMap
}
type FSBackend struct {
	dirLocation string
	catMutex    sync.RWMutex
	catFileLoc  string
	cache       FSCache
}

func (F *FSBackend) Init(pathToDir string) error {
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

	//create file to store categories, if it doesn't already exist
	F.catFileLoc = path.Join(pathToDir, "cats.json")
	if _, err := os.Stat(F.catFileLoc); os.IsNotExist(err) {
		F.cache = FSCache{
			Flat:      CategoryMap{},
			ParentMap: CategoryParentIDMap{},
		}
		err = jsonToFile(F.catFileLoc, F.cache)
		if err != nil {
			return err
		}

	} else {
		F.cache, err = F.getCategoriesFromFile()

		if err != nil {
			return err
		}
	}

	return nil
}

func (F *FSBackend) UploadFile(filename string, contents []byte) (FileID, error) {
	filepath := path.Join(F.dirLocation, filename)
	err := ioutil.WriteFile(filepath, contents, 0644)

	if err != nil {
		return "", err
	}

	return filename, nil
}

func (F *FSBackend) GetFile(id FileID) ([]byte, error) {
	// shit's hacky ha
	filepath := path.Join(F.dirLocation, id)

	return ioutil.ReadFile(filepath)
}

func (F *FSBackend) Files() ([]File, error) {
	osFiles, err := ioutil.ReadDir(F.dirLocation)
	if err != nil {
		return nil, err
	}

	files := make([]File, len(osFiles))

	for i, f := range osFiles {
		fType := SourceFile
		if strings.HasSuffix(f.Name(), "xlsx") {
			fType = DemoFile
		}

		files[i] = File{
			ID:   f.Name(),
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

func (F *FSBackend) getCategoriesFromFile() (FSCache, error) {
	var cache FSCache
	err := jsonFromFile(F.catFileLoc, &cache)
	return cache, err
}

func (F *FSBackend) writeCategoriesToFile(cache FSCache) error {
	err := jsonToFile(F.catFileLoc, cache)
	return err
}

func (F *FSBackend) CreateCategory(name string, parentCategoryID CategoryID) (CategoryID, error) {
	// todo: category names should be unique (per project, probably?)
	F.catMutex.Lock()
	defer F.catMutex.Unlock()

	newID := len(F.cache.Flat) + 1

	newCat := Category{Name: name, ID: newID, Texts: []DocumentText{}}

	if parentCategoryID == 0 {
		F.cache.ParentMap[newID] = []CategoryID{newID}
	} else {
		if subCats, ok := F.cache.ParentMap[parentCategoryID]; ok {
			F.cache.ParentMap[parentCategoryID] = append(subCats, newID)
		} else {
			return 0, errors.New(fmt.Sprintf("CategoryID not found. ID: %d", parentCategoryID))
		}
	}

	F.cache.Flat[newID] = newCat

	err := F.writeCategoriesToFile(F.cache)
	if err != nil {
		return 0, err
	}

	return newCat.ID, nil
}

func (F *FSBackend) CategorizeText(categoryID CategoryID, documentID FileID, text string, firstWord WordCoordinate, lastWord WordCoordinate) error {
	F.catMutex.Lock()
	defer F.catMutex.Unlock()

	cache, err := F.getCategoriesFromFile()
	if err != nil {
		return err
	}

	if _, ok := cache.Flat[categoryID]; ok == false {
		return errors.New(fmt.Sprintf("No category found with ID: %d", categoryID))
	}

	var cat = cache.Flat[categoryID]
	cat.Texts = append(cat.Texts, DocumentText{
		DocumentID: documentID,
		Text:       text,
		FirstWord:  firstWord,
		LastWord:   lastWord,
	})
	cache.Flat[categoryID] = cat

	err = F.writeCategoriesToFile(cache)

	F.cache = cache
	return err
}

func (F *FSBackend) GetCategory(categoryID CategoryID) (Category, error) {
	F.catMutex.RLock()
	defer F.catMutex.RUnlock()

	if cat, ok := F.cache.Flat[categoryID]; ok {
		return cat, nil
	}
	return Category{}, errors.New(fmt.Sprintf("No category found with ID: %d", categoryID))
}

func (F *FSBackend) GetCategoryMain(categoryID CategoryID) (CategoryMain, error) {
	F.catMutex.RLock()
	defer F.catMutex.RUnlock()

	if cat, ok := F.cache.Flat[categoryID]; ok {
		cMain := CategoryMain{
			Main:       cat.ID,
			Categories: make([]Category, len(F.cache.ParentMap[cat.ID])),
		}
		for i, subCatID := range F.cache.ParentMap[cat.ID] {
			cMain.Categories[i] = F.cache.Flat[subCatID]
		}
		return cMain, nil
	}

	return CategoryMain{}, errors.New(fmt.Sprintf("No category found with ID: %d", categoryID))
}

func (F *FSBackend) Categories() ([]CategoryID, error) {
	F.catMutex.RLock()
	defer F.catMutex.RUnlock()
	cache, err := F.getCategoriesFromFile()

	if err != nil {
		return nil, err
	}

	listCats := make([]CategoryID, len(cache.ParentMap))
	i := 0
	for mainID := range cache.ParentMap {
		listCats[i] = mainID
		i++
	}

	return listCats, nil
}

var _ Store = &FSBackend{}
