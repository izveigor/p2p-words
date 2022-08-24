package p2p

import (
	"encoding/gob"
	"log"
	"path/filepath"
	"sync"

	"github.com/izveigor/p2p-words/network/pkg/config"
)

func CreateFilesFunction(name string, text []byte) {
	WriteTXTFile(config.Config.PathToData+name, text)

	book, elements := CreateNewBook(
		LemmatizersServiceClient,
		&Tree{nil},
		name,
		string(text),
	)
	UserLibrary.Books = append(UserLibrary.Books, book)

	WriteGOBFile(config.Config.PathToData+name, book, elements)
}

func ReadFileFunction(path string) {
	file, err := OpenFileByPath(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	dec := gob.NewDecoder(file)
	bookInformation := BookInformation{}

	if err = dec.Decode(&bookInformation); err != nil {
		log.Fatal(err)
	}

	tree := &Tree{}
	tree.Create(bookInformation.TreeElements, nil, true)

	var lock sync.Mutex
	lock.Lock()
	UserLibrary.Books = append(UserLibrary.Books, &Book{
		Name:            bookInformation.Name,
		CharactersCount: bookInformation.CharactersCount,
		WordsCount:      bookInformation.WordsCount,
		Tree:            tree,
	})
	lock.Unlock()
}

var ReadFile = ReadFileFunction

func ReadFiles() {
	var path string = config.Config.PathToData
	files, err := ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		var filename = file.Name()
		if filepath.Ext(filename) == ".gob" {
			go ReadFile(path + filename)
		}
	}
}

var CreateFiles = CreateFilesFunction
