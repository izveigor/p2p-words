package p2p

import (
	"sort"

	"github.com/izveigor/p2p-words/network/pkg/p2p/pb"
)

type TreeElement struct {
	Word      string
	Sentences []string
}

type Node struct {
	Parent *Node
	Left   *Node
	Right  *Node
	Key    *TreeElement
}

type Tree struct {
	Root *Node
}

func NewNode(element *TreeElement) (node *Node) {
	node = &Node{
		Parent: nil,
		Left:   nil,
		Right:  nil,
		Key:    element,
	}

	return
}

func (t *Tree) Create(elements []*TreeElement, parent *Node, is_left bool) {
	var length int = len(elements)
	if length == 0 {
		return
	}

	var mid int = length / 2
	node := NewNode(elements[mid])

	if t.Root == nil {
		t.Root = node
	}

	if parent != nil {
		node.Parent = parent
		if is_left {
			parent.Left = node
		} else {
			parent.Right = node
		}
	}

	t.Create(elements[:mid], node, true)
	t.Create(elements[mid+1:], node, false)
}

func (t *Tree) Search(word string) *TreeElement {
	var x *Node = t.Root
	for x != nil && word != x.Key.Word {
		if word < x.Key.Word {
			x = x.Left
		} else {
			x = x.Right
		}
	}

	if x != nil {
		return x.Key
	}
	return nil
}

func CreateTreeElementsFunction(lemmatizedText *pb.LemmatizerResponse) (elements []*TreeElement, charactersCount int, wordsCount int) {
	counter := map[string][]string{}
	charactersCount = 0
	wordsCount = 0
	setSentences := make(map[string]bool)

	for _, body := range lemmatizedText.Words {
		var word, sentence string = body.Word, body.Sentence
		counter[word] = append(counter[word], sentence)
		if _, ok := setSentences[sentence]; !ok {
			charactersCount += len(sentence)
		}
		wordsCount += 1
		setSentences[sentence] = true
	}

	for word, sentences := range counter {
		elements = append(elements, &TreeElement{
			Word:      word,
			Sentences: sentences,
		})
	}

	sort.Slice(
		elements, func(i, j int) bool {
			return elements[i].Word < elements[j].Word
		},
	)

	return
}

var CreateTreeElements = CreateTreeElementsFunction
