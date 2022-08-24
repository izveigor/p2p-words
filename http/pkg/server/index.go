package server

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/izveigor/p2p-words/http/pkg/server/pb"
)

func Index(HTTPServiceClient *Client, w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(2 << 32)
	switch r.Method {
	case "GET":
		response, err := HTTPServiceClient.Service.GetBooksInformation(
			context.Background(),
			&pb.GetBooksInformationRequest{},
		)
		if err != nil {
			panic(err)
		}

		tmpl, err := template.New("").ParseFiles(PATH_TO_STATIC+"base.html", PATH_TO_STATIC+"index.html")
		if err != nil {
			InternalServerError(w, err)
		}

		if err = tmpl.ExecuteTemplate(w, "base", response); err != nil {
			InternalServerError(w, err)
		}

	case "POST":
		file, header, err := r.FormFile("file")

		if err != nil || filepath.Ext(header.Filename) != ".txt" {
			http.Error(w, "Wrong file", http.StatusBadRequest)
		}
		defer file.Close()

		var buffer bytes.Buffer
		io.Copy(&buffer, file)

		response, err := HTTPServiceClient.Service.CreateBook(
			context.Background(),
			&pb.CreateBookRequest{
				Text: buffer.Bytes(),
				Name: header.Filename,
			},
		)

		if err == nil && response.Ok == true {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			InternalServerError(w, err)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
