package server

import (
	"github.com/golang/mock/gomock"
	"github.com/izveigor/p2p-words/http/pkg/mock"
	"github.com/izveigor/p2p-words/http/pkg/server/pb"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"

	"net/url"
	"testing"
)

func TestSearchGet(t *testing.T) {
	wr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/search", nil)

	Search(&Client{Service: nil}, wr, req)
	assert.Equal(t, wr.Code, http.StatusOK)
}

func TestSearchPost(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer func() {
		ctrl.Finish()
	}()

	mockHTTPClient := mock.NewMockHTTPClient(ctrl)

	mockHTTPClient.EXPECT().SearchSentences(
		gomock.Any(),
		gomock.Any(),
	).Return(
		&pb.SearchSentencesResponse{
			Sentences: []string{
				"Предложение",
				"Другое Предложение",
			},
		}, nil,
	)

	wr := httptest.NewRecorder()

	postBody := url.Values{}
	postBody.Set("word", "Слово")

	req := httptest.NewRequest(http.MethodPost, "/search", nil)
	req.PostForm = postBody
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	Search(&Client{mockHTTPClient}, wr, req)
	assert.Equal(t, wr.Code, http.StatusOK)
	assert.Equal(t, wr.Body.String(), "{\"sentences\":[\"Предложение\",\"Другое Предложение\"]}\n")
}
