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

type FSBackend struct {
	dirLocation string
	catMutex    sync.RWMutex
	catFileLoc  string
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
		cats := Categories{}
		err = jsonToFile(F.catFileLoc, cats)
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

func (F *FSBackend) getCategoriesFromFile() (Categories, error) {
	var cats Categories
	err := jsonFromFile(F.catFileLoc, &cats)
	return cats, err
}

func (F *FSBackend) writeCategoriesToFile(cats Categories) error {
	err := jsonToFile(F.catFileLoc, cats)
	return err
}

func (F *FSBackend) CreateCategory(name string) (CategoryID, error) {
	// todo: category names should be unique (across documents?) either way,
	//  to be enforced by the db via constraints (check those errors)
	F.catMutex.Lock()
	defer F.catMutex.Unlock()
	cats, err := F.getCategoriesFromFile()
	if err != nil {
		return 0, err
	}

	newId := len(cats) + 1

	// since the ID is just going to be the name (until there's a db providing AI,
	//  use the name as the ID. Therefor, no point in check if the name already exists
	newCat := Category{Name: name, ID: newId, Texts: []DocumentText{}}

	cats[newId] = newCat

	err = F.writeCategoriesToFile(cats)
	if err != nil {
		return 0, err
	}

	return newCat.ID, nil
}

func (F *FSBackend) CategorizeText(categoryID CategoryID, documentID FileID, text string) error {
	F.catMutex.Lock()
	defer F.catMutex.Unlock()
	cats, err := F.getCategoriesFromFile()
	if err != nil {
		return err
	}

	if _, ok := cats[categoryID]; ok == false {
		return errors.New(fmt.Sprintf("No category found with ID: %d", categoryID))
	}

	var cat = cats[categoryID]
	cat.Texts = append(cat.Texts, DocumentText{
		DocumentID: documentID,
		Text:       text,
	})
	cats[categoryID] = cat

	err = F.writeCategoriesToFile(cats)
	return err
}

func (F *FSBackend) GetCategory(categoryID CategoryID) (Category, error) {
	F.catMutex.RLock()
	defer F.catMutex.RUnlock()
	cats, err := F.getCategoriesFromFile()
	if err != nil {
		return Category{}, err
	}

	if cat, ok := cats[categoryID]; ok {
		return cat, nil
	}

	return Category{}, errors.New(fmt.Sprintf("No category found with ID: %d", categoryID))
}

func (F *FSBackend) Categories() ([]CategoryID, error) {
	F.catMutex.RLock()
	defer F.catMutex.RUnlock()
	currentCats, err := F.getCategoriesFromFile()

	if err != nil {
		return nil, err
	}

	listCats := make([]CategoryID, 0)

	for _, v := range currentCats {
		listCats = append(listCats, v.ID)
	}

	return listCats, nil
}

var _ Store = &FSBackend{}
