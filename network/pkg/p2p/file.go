package p2p

import (
	"encoding/gob"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
)

type BookInformation struct {
	Name            string
	CharactersCount int
	WordsCount      int
	TreeElements    []*TreeElement
}

func WriteTXTFileFunction(name string, text []byte) {
	err := ioutil.WriteFile(name[:len(name)-4]+".txt", text, 0644)

	if err != nil {
		log.Fatal(err)
	}
}

func WriteGOBFileFunction(name string, book *Book, elements []*TreeElement) {
	file, _ := os.Create(name[:len(name)-4] + ".gob")
	defer file.Close()

	bookInformation := BookInformation{
		Name:            book.Name,
		CharactersCount: book.CharactersCount,
		WordsCount:      book.WordsCount,
		TreeElements:    elements,
	}

	enc := gob.NewEncoder(file)
	if err := enc.Encode(&bookInformation); err != nil {
		log.Fatal(err)
	}
}

func OpenFileByPathFunction(name string) (*os.File, error)  { return os.Open(name) }
func ReadDirFunction(dirname string) ([]fs.FileInfo, error) { return ioutil.ReadDir(dirname) }

var (
	WriteTXTFile   = WriteTXTFileFunction
	WriteGOBFile   = WriteGOBFileFunction
	OpenFileByPath = OpenFileByPathFunction
	ReadDir        = ReadDirFunction
)
