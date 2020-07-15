package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

type CSCache struct {
	Flat      CategoryMap
	ParentMap CategoryParentIDMap
}
type FSBackendFile struct {
	dirLocation string
}
type FSBackendCategory struct {
	catMutex   sync.RWMutex
	catFileLoc string
	cache      CSCache
}

func (F *FSBackendFile) Init(pathToDir string) error {
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
func (F *FSBackendCategory) Init(pathToDir string) error {
	if _, err := os.Stat(pathToDir); os.IsNotExist(err) {
		err = os.Mkdir(pathToDir, 0766)

		if err != nil {
			return err
		}
	}

	//create file to store categories, if it doesn't already exist
	F.catFileLoc = path.Join(pathToDir, "cats.json")
	if _, err := os.Stat(F.catFileLoc); os.IsNotExist(err) {
		F.cache = CSCache{
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

func (F *FSBackendFile) UploadFile(filename string, contents []byte) (FileID, error) {
	filepath := path.Join(F.dirLocation, filename)
	err := ioutil.WriteFile(filepath, contents, 0644)

	if err != nil {
		return "", err
	}

	return filename, nil
}

func (F *FSBackendFile) GetFile(id FileID) ([]byte, error) {
	// shit's hacky ha
	filepath := path.Join(F.dirLocation, id)

	return ioutil.ReadFile(filepath)
}

func (F *FSBackendFile) Files() ([]string, error) {
	files, err := ioutil.ReadDir(F.dirLocation)
	if err != nil {
		return nil, err
	}
	var names []string

	for _, f := range files {
		names = append(names, f.Name())
	}

	return names, nil
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

func (F *FSBackendCategory) getCategoriesFromFile() (CSCache, error) {
	var cache CSCache
	err := jsonFromFile(F.catFileLoc, &cache)
	return cache, err
}

func (F *FSBackendCategory) writeCategoriesToFile(cache CSCache) error {
	err := jsonToFile(F.catFileLoc, cache)
	return err
}

func (F *FSBackendCategory) CreateCategory(name string, parentCategoryID CategoryID) (CategoryID, error) {
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

func (F *FSBackendCategory) CategorizeText(categoryID CategoryID, documentID FileID, text string, firstWord WordCoordinate, lastWord WordCoordinate) error {
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

func (F *FSBackendCategory) GetCategory(categoryID CategoryID) (Category, error) {
	F.catMutex.RLock()
	defer F.catMutex.RUnlock()

	if cat, ok := F.cache.Flat[categoryID]; ok {
		return cat, nil
	}
	return Category{}, errors.New(fmt.Sprintf("No category found with ID: %d", categoryID))
}

func (F *FSBackendCategory) GetCategoryMain(categoryID CategoryID) (CategoryMain, error) {
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

func (F *FSBackendCategory) Categories() ([]CategoryID, error) {
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
