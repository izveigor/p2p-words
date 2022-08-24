package server

import (
	"github.com/golang/mock/gomock"
	"github.com/izveigor/p2p-words/http/pkg/mock"
	"github.com/izveigor/p2p-words/http/pkg/server/pb"
	"github.com/stretchr/testify/assert"

	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIndexGet(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer func() {
		ctrl.Finish()
	}()

	mockHTTPClient := mock.NewMockHTTPClient(ctrl)

	mockHTTPClient.EXPECT().GetBooksInformation(
		gomock.Any(),
		gomock.Any(),
	).Return(
		&pb.GetBooksInformationResponse{
			Books: []*pb.Book{
				{
					Name:            "Название",
					CharactersCount: 50,
					WordsCount:      20,
				},
			},
			Number: 1,
		}, nil,
	)

	wr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	Index(&Client{mockHTTPClient}, wr, req)
	assert.Equal(t, wr.Code, http.StatusOK)

	assert.NotEqual(t, strings.Index(wr.Body.String(), "<h5 class=\"text-center\">Название: Название</h5>"), -1)
	assert.NotEqual(t, strings.Index(wr.Body.String(), "<span class=\"col-9\">50</span>"), -1)
	assert.NotEqual(t, strings.Index(wr.Body.String(), "<span class=\"col-9\">20</span>"), -1)
}

func TestIndexPost(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer func() {
		ctrl.Finish()
	}()

	mockHTTPClient := mock.NewMockHTTPClient(ctrl)

	mockHTTPClient.EXPECT().CreateBook(
		gomock.Any(),
		gomock.Any(),
	).Return(
		&pb.CreateBookResponse{
			Ok: true,
		}, nil,
	)

	wr := httptest.NewRecorder()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "file.txt")
	part.Write([]byte("Текст..."))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	Index(&Client{mockHTTPClient}, wr, req)

	assert.Equal(t, wr.Code, http.StatusSeeOther)
}
