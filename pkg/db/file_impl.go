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

type CacheCategory struct {
	*Category
	Subcategories []CategoryID `json:"subcategories"`
}
type CacheCategories = map[CategoryID]*CacheCategory

type FSBackend struct {
	dirLocation string
	catMutex    sync.RWMutex
	catFileLoc  string
	cache       Categories
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
		F.cache = Categories{}
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

func (F *FSBackend) Files() ([]string, error) {
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

/*
	Reads the cache file (converts from "Cache structs" to the "interface structs")
	- The "Cache structs" store the data in a way that reduces file size
	- The "interface structs" store data in a way that is convenient to look up data
*/
func (F *FSBackend) getCategoriesFromFile() (Categories, error) {
	var cacheCats CacheCategories
	retCats := Categories{}
	err := jsonFromFile(F.catFileLoc, &cacheCats)

	if err != nil {
		return nil, err
	}

	for catID, cacheCat := range cacheCats {
		newCat := cacheCat.Category
		//newCat.Subcategories = Subcategories{}
		retCats[catID] = newCat
	}
	for catID, cacheCat := range cacheCats {
		if len(cacheCat.Subcategories) == 0 {
			continue
		}

		for _, subCatID := range cacheCat.Subcategories {
			retCats[catID].Subcategories = append(retCats[catID].Subcategories, retCats[subCatID])
		}
	}

	return retCats, err
}

/*
	Writes to the cache file (converts from "interface structs" to "Cache structs". see getCategoriesFromFile comment)
*/
func (F *FSBackend) writeCategoriesToFile(cats Categories) error {
	cacheCats := CacheCategories{}

	for catID, cat := range cats {
		cacheCat := CacheCategory{
			Category:      cat,
			Subcategories: []CategoryID{},
		}

		for _, subCat := range cat.Subcategories {
			cacheCat.Subcategories = append(cacheCat.Subcategories, subCat.ID)
		}
		cacheCats[catID] = &cacheCat
	}
	err := jsonToFile(F.catFileLoc, cacheCats)
	return err
}

func (F *FSBackend) CreateCategory(name string, parentCategoryID CategoryID) (CategoryID, error) {
	// todo: category names should be unique (per project, probably?)
	F.catMutex.Lock()
	defer F.catMutex.Unlock()

	newId := len(F.cache) + 1

	newCat := Category{Name: name, ID: newId, Texts: []DocumentText{}, Subcategories: []*Category{}}

	if parentCategoryID != 0 {
		if cat, ok := F.cache[parentCategoryID]; ok {
			newCat.IsSub = true
			cat.Subcategories = append(F.cache[parentCategoryID].Subcategories, &newCat)
			F.cache[parentCategoryID] = cat
		}
	}

	F.cache[newId] = &newCat

	err := F.writeCategoriesToFile(F.cache)
	if err != nil {
		return 0, err
	}

	return newCat.ID, nil
}

func (F *FSBackend) CategorizeText(categoryID CategoryID, documentID FileID, text string, firstWord WordCoordinate, lastWord WordCoordinate) error {
	F.catMutex.Lock()
	defer F.catMutex.Unlock()

	if _, ok := F.cache[categoryID]; ok == false {
		return errors.New(fmt.Sprintf("No category found with ID: %d", categoryID))
	}

	var cat = F.cache[categoryID]
	cat.Texts = append(cat.Texts, DocumentText{
		DocumentID: documentID,
		Text:       text,
		FirstWord:  firstWord,
		LastWord:   lastWord,
	})
	F.cache[categoryID] = cat

	err := F.writeCategoriesToFile(F.cache)
	return err
}

func (F *FSBackend) GetCategory(categoryID CategoryID) (Category, error) {
	F.catMutex.RLock()
	defer F.catMutex.RUnlock()

	if cat, ok := F.cache[categoryID]; ok {
		return *cat, nil
	}

	return Category{}, errors.New(fmt.Sprintf("No category found with ID: %d", categoryID))
}

func (F *FSBackend) Categories() ([]CategoryID, error) {
	F.catMutex.RLock()
	defer F.catMutex.RUnlock()

	listCats := make([]CategoryID, 0)

	for _, v := range F.cache {
		listCats = append(listCats, v.ID)
	}

	return listCats, nil
}

var _ Store = &FSBackend{}
