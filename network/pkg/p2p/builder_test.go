package p2p

import (
	"encoding/gob"
	"github.com/izveigor/p2p-words/network/pkg/config"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateFiles(t *testing.T) {
	defer func() {
		WriteTXTFile = WriteTXTFileFunction
		WriteGOBFile = WriteGOBFileFunction
		CreateNewBook = CreateNewBookFunction
		UserLibrary = Library{[]*Book{}}
	}()

	var nameArgument string = "Имя"
	var textArgument []byte = []byte("Текст")

	var (
		WriteTXTFileWasCalled  bool = false
		WriteGOBFileWasCalled  bool = false
		CreateNewBookWasCalled bool = false
	)

	CreateNewBook = func(LemmatizersServiceClient Client, tree *Tree, name, text string) (*Book, []*TreeElement) {
		CreateNewBookWasCalled = true
		assert.Equal(t, name, nameArgument)
		assert.Equal(t, text, string(textArgument))
		return &Book{
			Name:            name,
			CharactersCount: 50,
			WordsCount:      20,
			Tree:            tree,
		}, []*TreeElement{}
	}

	WriteTXTFile = func(name string, text []byte) {
		WriteTXTFileWasCalled = true
		assert.Equal(t, name, config.Config.PathToData+nameArgument)
		assert.Equal(t, text, textArgument)
	}

	WriteGOBFile = func(name string, book *Book, elements []*TreeElement) {
		WriteGOBFileWasCalled = true
		assert.Equal(t, name, config.Config.PathToData+nameArgument)
	}

	UserLibrary = Library{[]*Book{}}
	CreateFiles(nameArgument, textArgument)

	assert.Equal(t, UserLibrary.Books[0].Name, nameArgument)
	assert.Equal(t, UserLibrary.Books[0].CharactersCount, 50)
	assert.Equal(t, UserLibrary.Books[0].WordsCount, 20)
	assert.IsType(t, Tree{}, *UserLibrary.Books[0].Tree)

	if !(CreateNewBookWasCalled && WriteTXTFileWasCalled && WriteGOBFileWasCalled) {
		t.Fail()
	}
}

func TestReadFile(t *testing.T) {
	defer func() {
		OpenFileByPath = OpenFileByPathFunction
		UserLibrary = Library{[]*Book{}}
		if err := os.RemoveAll("temp"); err != nil {
			t.Fatal(err)
		}
	}()
	book := Book{
		Name:            "Имя",
		CharactersCount: 50,
		WordsCount:      20,
		Tree: &Tree{Root: NewNode(&TreeElement{
			Word:      "Слово",
			Sentences: []string{"Предложение"},
		})},
	}
	OpenFileByPath = func(path string) (*os.File, error) {
		assert.Equal(t, path, "path")
		if err := os.Mkdir("temp", 644); err != nil {
			t.Fatal(err)
		}
		file, err := os.Create("temp/book.gob")
		if err != nil {
			t.Fatal(err)
		}

		enc := gob.NewEncoder(file)

		bookInformation := BookInformation{
			Name:            book.Name,
			CharactersCount: book.CharactersCount,
			WordsCount:      book.WordsCount,
			TreeElements: []*TreeElement{
				{
					Word:      "Слово",
					Sentences: []string{"Предложение"},
				},
			},
		}
		if err := enc.Encode(&bookInformation); err != nil {
			t.Fatal(err)
		}

		file, err = os.Open("temp/book.gob")
		if err != nil {
			t.Fatal(err)
		}

		return file, nil
	}

	UserLibrary = Library{[]*Book{}}
	ReadFile("path")
	assert.Equal(t, UserLibrary.Books[0].Name, book.Name)
	assert.Equal(t, UserLibrary.Books[0].CharactersCount, book.CharactersCount)
	assert.Equal(t, UserLibrary.Books[0].WordsCount, book.WordsCount)
	assert.Equal(t, UserLibrary.Books[0].Tree.Root.Key.Word, book.Tree.Root.Key.Word)
	assert.Equal(t, UserLibrary.Books[0].Tree.Root.Key.Sentences, book.Tree.Root.Key.Sentences)
}

func TestReadFiles(t *testing.T) {
	defer func() {
		ReadDir = ReadDirFunction
		ReadFile = ReadFileFunction
		if err := os.RemoveAll("temp"); err != nil {
			t.Fatal(err)
		}
	}()
	ReadDir = func(dirname string) ([]fs.FileInfo, error) {
		if err := os.Mkdir("temp", 0644); err != nil {
			t.Fatal(err)
		}
		var files [3]string = [3]string{
			"temp/first.gob",
			"temp/text.txt",
			"temp/second.gob",
		}
		for _, filename := range files {
			if _, err := os.Create(filename); err != nil {
				t.Fatal(err)
			}
		}

		return ioutil.ReadDir("temp")
	}
	ReadFile = func(path string) {
		assert.Equal(t, filepath.Ext(path), ".gob")
	}
	ReadFiles()
}
