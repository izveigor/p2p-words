package p2p

import (
	"testing"

	"github.com/izveigor/p2p-words/network/pkg/p2p/pb"
	"github.com/stretchr/testify/assert"
)

func TestTreeElement(t *testing.T) {
	var first_sentence string = "Hello, world!"
	tree := TreeElement{
		Word:      "hello",
		Sentences: []string{first_sentence},
	}

	assert.IsType(t, TreeElement{}, tree)
}

func TestTree(t *testing.T) {
	var first_sentence string = "В 1736 году часовню передали Афанасьевскому монастырю по просьбе игумена Ионы" // https://ru.wikipedia.org/wiki/%D0%A6%D0%B5%D1%80%D0%BA%D0%BE%D0%B2%D1%8C_%D0%B8%D0%BA%D0%BE%D0%BD%D1%8B_%D0%91%D0%BE%D0%B6%D0%B8%D0%B5%D0%B9_%D0%9C%D0%B0%D1%82%D0%B5%D1%80%D0%B8_%D0%97%D0%BD%D0%B0%D0%BC%D0%B5%D0%BD%D0%B8%D0%B5_(%D0%AF%D1%80%D0%BE%D1%81%D0%BB%D0%B0%D0%B2%D0%BB%D1%8C)
	var second_sentence string = "Общая протяжённость всех водотоков на территории Санкт-Петербурга " +
		"достигает 282 км, а их водная поверхность составляет около 7 % всей площади города." // https://ru.wikipedia.org/wiki/%D0%A1%D0%B0%D0%BD%D0%BA%D1%82-%D0%9F%D0%B5%D1%82%D0%B5%D1%80%D0%B1%D1%83%D1%80%D0%B3
	elements := []*TreeElement{
		{Word: "водная", Sentences: []string{second_sentence}},
		{Word: "общая", Sentences: []string{second_sentence}},
		{Word: "просьба", Sentences: []string{first_sentence}},
		{Word: "составлять", Sentences: []string{second_sentence}},
		{Word: "часовня", Sentences: []string{first_sentence}},
	}

	tree := Tree{nil}
	tree.Create(elements, nil, true)

	assert.Equal(t, tree.Root.Left.Key.Word, "общая")
	assert.Equal(t, tree.Root.Key.Word, "просьба")
	assert.Equal(t, tree.Root.Right.Key.Word, "часовня")
	assert.Equal(t, tree.Root.Right.Left.Key.Word, "составлять")
	assert.Equal(t, tree.Root.Left.Left.Key.Word, "водная")

	assert.Equal(t, tree.Search("общая"), elements[1])
	assert.Equal(t, tree.Search("часовня"), elements[4])
	assert.Nil(t, tree.Search("слово"))
}

func TestCreateTreeElements(t *testing.T) {
	lemmatizedText := pb.LemmatizerResponse{
		Words: []*pb.LemmatizedWord{
			{Word: "привет", Sentence: "Привет, мир!", Id: 1},
			{Word: "привет", Sentence: "Привет, Владивосток!", Id: 2},
			{Word: "мир", Sentence: "Привет, мир!", Id: 3},
			{Word: "владивосток", Sentence: "Привет, Владивосток!", Id: 4},
		},
	}
	elements, charactersCount, wordsCount := CreateTreeElements(&lemmatizedText)

	assert.Equal(t, charactersCount, 58)
	assert.Equal(t, wordsCount, 4)
	assert.Equal(t, elements[0], &TreeElement{
		Word:      "владивосток",
		Sentences: []string{"Привет, Владивосток!"},
	})

	assert.Equal(t, elements[1], &TreeElement{
		Word:      "мир",
		Sentences: []string{"Привет, мир!"},
	})

	assert.Equal(t, elements[2], &TreeElement{
		Word:      "привет",
		Sentences: []string{"Привет, мир!", "Привет, Владивосток!"},
	})
}

func TestNode(t *testing.T) {
	element := TreeElement{
		Word:      "привет",
		Sentences: []string{"Предложение"},
	}

	node := NewNode(&element)
	assert.Equal(t, node, &Node{
		Parent: nil,
		Left:   nil,
		Right:  nil,
		Key:    &element,
	})
}
