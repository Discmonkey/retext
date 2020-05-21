package db

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type FSBackend struct {
	dirLocation string
}

func (F *FSBackend) Init(pathToDir string) error {
	if _, err := os.Stat(pathToDir); os.IsNotExist(err) {
		err = os.Mkdir(pathToDir, 0766)

		if err != nil {
			return err
		}
	}

	F.dirLocation = pathToDir

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

//used for storing things in a file (remove if DB starts being used)
//create a category-json file in the OS's temp directory
var ex, _ = os.Executable()
var exPath = filepath.Dir(ex)
var categoriesFile = exPath + "/cat.json"

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

func getCategoriesFromFile() (Categories, error) {
	var cat Categories
	err := jsonFromFile(categoriesFile, &cat)
	return cat, err
}

func (F *FSBackend) CreateCategory(name string) (Category, error) {
	// todo: use a db when worried about things like "concurrent requests"
	//  category names should be unique (across documents?) either way,
	//  to be enforced by the db via constraints (check those errors)
	cats, err := getCategoriesFromFile()
	if err != nil {
		if os.IsNotExist(err) { // no file? no problem, we're about to write to it
			cats = Categories{Categories: map[string]Category{}}
		} else {
			return Category{}, err
		}
	}

	// since the ID is just going to be the name (until there's a db providing AI,
	//  use the name as the ID. Therefor, no point in check if the name already exists
	newCat := Category{Name: name, ID: name, Texts: []DocumentText{}}

	cats.Categories[name] = newCat

	err = jsonToFile(categoriesFile, cats)
	if err != nil {
		return Category{}, err
	}

	return newCat, nil
}

func (F *FSBackend) CategorizeText(categoryID CategoryID, documentID FileID, text string) error {
	cats, err := getCategoriesFromFile()
	if err != nil {
		return err
	}

	if _, ok := cats.Categories[categoryID]; ok == false {
		return errors.New("No category found with ID: " + categoryID)
	}

	var cat = cats.Categories[categoryID]
	cat.Texts = append(cat.Texts, DocumentText{
		DocumentID: documentID,
		Text:       text,
	})
	cats.Categories[categoryID] = cat

	err = jsonToFile(categoriesFile, cats)
	return err
}

func (F *FSBackend) GetCategory(categoryID CategoryID) (Category, error) {
	cats, err := getCategoriesFromFile()
	if err != nil {
		return Category{}, err
	}

	if cat, ok := cats.Categories[categoryID]; ok {
		return cat, nil
	}

	// todo: is this how "no category found" should be returned?
	return Category{}, nil
}

func (F *FSBackend) Categories() ([]CategoryID, error) {
	currentCats, err := getCategoriesFromFile()

	if err != nil {
		return nil, err
	}

	listCats := make([]CategoryID, 0)

	for _, v := range currentCats.Categories {
		listCats = append(listCats, v.ID)
	}

	return listCats, nil
}

var _ Store = &FSBackend{}
