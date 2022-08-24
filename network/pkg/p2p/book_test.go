package p2p

import (
	"github.com/izveigor/p2p-words/network/pkg/p2p/pb"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	"github.com/izveigor/p2p-words/network/pkg/mock"
	"testing"
)

func MockCreateTreeElements(lemmatizedTest *pb.LemmatizerResponse) (elements []*TreeElement, charactersCount int, wordsCount int) {
	elements = []*TreeElement{
		{Word: "Слово", Sentences: []string{"Предложение"}},
	}
	charactersCount = 50
	wordsCount = 100
	return
}

func TestCreateNewBook(t *testing.T) {
	ctrl := gomock.NewController(t)

	CreateTreeElements = MockCreateTreeElements
	defer func() {
		CreateTreeElements = CreateTreeElementsFunction
		ctrl.Finish()
	}()

	mockLemmatizersClient := mock.NewMockLemmatizersClient(ctrl)

	mockLemmatizersClient.EXPECT().Lemmatize(
		gomock.Any(),
		gomock.Any(),
	).Return(
		&pb.LemmatizerResponse{
			Words: []*pb.LemmatizedWord{
				{Word: "Слово", Sentence: "Предложение", Id: 1},
			},
		}, nil,
	)
	book, elements := CreateNewBook(Client{mockLemmatizersClient}, &Tree{nil}, "Название", "Текст")

	assert.Equal(t, elements[0].Word, "Слово")
	assert.Equal(t, elements[0].Sentences, []string{"Предложение"})

	assert.Equal(t, book.Name, "Название")
	assert.Equal(t, book.CharactersCount, 50)
	assert.Equal(t, book.WordsCount, 100)
	assert.Equal(t, book.Tree.Root.Key.Word, "Слово")
	assert.Equal(t, book.Tree.Root.Key.Sentences, []string{"Предложение"})
}

func TestBookSearchSentences(t *testing.T) {
	book := &Book{
		Name:            "Название",
		CharactersCount: 200,
		WordsCount:      100,
		Tree: &Tree{
			Root: NewNode(&TreeElement{
				Word:      "Слово",
				Sentences: []string{"Предложение"},
			}),
		},
	}
	partSentences := make(chan []string)
	go book.SearchSentences("Слово", partSentences)
	var sentences []string = <-partSentences
	assert.Equal(t, sentences, []string{"Предложение"})
}

func TestLibrarySearchSentences(t *testing.T) {
	firstBook := &Book{
		Name:            "Название",
		CharactersCount: 200,
		WordsCount:      100,
		Tree: &Tree{
			Root: NewNode(&TreeElement{
				Word:      "Слово",
				Sentences: []string{"Предложение"},
			}),
		},
	}
	secondBook := &Book{
		Name:            "Название",
		CharactersCount: 200,
		WordsCount:      100,
		Tree: &Tree{
			Root: NewNode(&TreeElement{
				Word:      "Другое",
				Sentences: []string{"Другое предложение"},
			}),
		},
	}

	UserLibrary.Books = append(UserLibrary.Books, firstBook)
	UserLibrary.Books = append(UserLibrary.Books, secondBook)

	sentences := UserLibrary.SearchSentences("Слово")
	assert.Equal(t, sentences, []string{"Предложение"})
}

func TestLibraryGetBooksInformation(t *testing.T) {
	defer func() {
		UserLibrary = Library{[]*Book{}}
	}()

	book := &Book{
		Name:            "Название",
		CharactersCount: 50,
		WordsCount:      20,
		Tree:            &Tree{nil},
	}
	UserLibrary = Library{[]*Book{}}
	UserLibrary.Books = append(UserLibrary.Books, book)
	response := UserLibrary.GetBooksInformation()
	assert.Equal(t, response, &pb.GetBooksInformationResponse{
		Number: int32(1),
		Books: []*pb.Book{
			{
				Name:            "Название",
				CharactersCount: int32(50),
				WordsCount:      int32(20),
			},
		},
	})
}
