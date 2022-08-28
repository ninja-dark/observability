package handler

import (
	"context"
	"elastic/m"
	"elastic/store"
	"encoding/json"
	"fmt"
	"net/http"

)

type ArticleHandler struct {
	S store.ArticleStore
}
// writeResponse - вспомогательная функция, которая записывет http статус-код и текстовое сообщение в ответ клиенту.
func writeResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte(message))
	_, _ = w.Write([]byte("\n"))
}
// writeJsonResponse - вспомогательная функция, которая запсывает http статус-код и сообщение в формате json в ответ клиенту.
func writeJsonResponse(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, fmt.Sprintf("can't marshal data: %s", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	writeResponse(w, status, string(response))
}



type ctxKey struct {
    name string
}

func NewArticleHandler(s store.ArticleStore) ArticleHandler {
	return ArticleHandler{S: s}
}
func (h ArticleHandler) Id(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
    articleID, ok := ctx.Value(ctxKey{"id"}).(int)
    if !ok {
        http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
        return
    }

    w.Write([]byte(fmt.Sprintf("article ID:%d", articleID)))
}

func (h ArticleHandler) Add(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var article m.Article
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		writeJsonResponse(w, http.StatusBadRequest, err)
        return
	}
	err = h.S.Add(ctx, article)
	if err != nil {
		writeJsonResponse(w, http.StatusBadRequest, err)
        return
	}
	writeJsonResponse(w, http.StatusOK, article)
}

type SearchRequest struct {
	Query string `json:"query"`
}

func (h ArticleHandler) Search(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var query SearchRequest
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		writeJsonResponse(w, http.StatusBadRequest, err)
        return
	}
	articles, err := h.S.Search(ctx, query.Query)
	if err != nil {
		writeJsonResponse(w, http.StatusBadRequest, err)
        return
	}
	writeJsonResponse(w, http.StatusOK, articles)

}