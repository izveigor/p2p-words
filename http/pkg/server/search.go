package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"text/template"

	"github.com/izveigor/p2p-words/http/pkg/server/pb"
)

type searchResult struct {
	Sentences []string `json:"sentences"`
}

type searchInput struct {
	Word string
}

func Search(HTTPServiceClient *Client, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl, err := template.New("").ParseFiles(PATH_TO_STATIC+"base.html", PATH_TO_STATIC+"search.html")
		if err != nil {
			InternalServerError(w, err)
		}

		if err = tmpl.ExecuteTemplate(w, "base", nil); err != nil {
			InternalServerError(w, err)
		}
	case "POST":
		input := searchInput{}
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&input)
		var word string = strings.ToLower(input.Word)

		response, err := HTTPServiceClient.Service.SearchSentences(
			context.Background(),
			&pb.SearchSentencesRequest{
				Word: word,
			},
		)
		if err != nil {
			http.Error(w, "Wrong word", http.StatusBadRequest)
		}

		data := searchResult{
			Sentences: response.Sentences,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
