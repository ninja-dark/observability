package handler

import (
	"context"
	"elastic/m"
	"elastic/store"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ArticleHandler struct {
	S      store.ArticleStore
	logger *zap.Logger
	tracer opentracing.Tracer
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

func NewArticleHandler(s store.ArticleStore, logger *zap.Logger, tracer opentracing.Tracer) ArticleHandler {
	return ArticleHandler{
		S:      s,
		logger: logger,
		tracer: tracer,
	}
}


func (h ArticleHandler) Id(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(context.Background(), h.tracer, "article.handler.Id")
	defer span.Finish()
	h.logger.Info("ArticleHandler id called", zap.Field{Key: "method", String: r.Method, Type:zapcore.StringType})
	
	id := chi.URLParam(r, "id")
	

	article, err := h.S.Get(ctx, id)
	if err != nil {
		h.logger.Error(fmt.Sprintf("cannot get article id: %s", err))
		span.LogFields(log.Error(err),
	)
		span.LogFields(
			log.Error(err),
			log.String("articleId", article.Id),
		)

	}
	span.LogFields(log.String("get article by id", fmt.Sprintf("%v", article)))
	
	writeJsonResponse(w, http.StatusOK, article)
}

func (h ArticleHandler) Add(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(context.Background(), h.tracer, "article.handler.Add")
	defer span.Finish()

	h.logger.Info("ArticleHandler called add", zap.Field{Key: "method", String: r.Method, Type:zapcore.StringType})
	var article m.Article
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		h.logger.Error(fmt.Sprintf("cannot decode JSON: %s", err))
		span.LogFields(log.Error(err))
		
		writeJsonResponse(w, http.StatusBadRequest, err)
		return
	}
	err = h.S.Add(ctx, article)
	if err != nil {
		h.logger.Error(fmt.Sprintf("cannot add article: %s", err))
		span.LogFields(log.Error(err))
		writeJsonResponse(w, http.StatusBadRequest, err)
		return
	}
	span.LogFields(log.String("Add article", fmt.Sprintf("%v", article)))
	writeJsonResponse(w, http.StatusOK, article)
}

type SearchRequest struct {
	Query string `json:"query"`
}

func (h ArticleHandler) Search(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(context.Background(), h.tracer, "article.handler.search")
	defer span.Finish()
	
	h.logger.Info("ArticleHandler called search", zap.Field{Key: "method", String: r.Method, Type:zapcore.StringType})
	var query SearchRequest
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		h.logger.Error(fmt.Sprintf("cannot decode JSON: %s", err))
		span.LogFields(log.Error(err))
		writeJsonResponse(w, http.StatusBadRequest, err)
		return
	}
	articles, err := h.S.Search(ctx, query.Query)
	if err != nil {
		h.logger.Error(fmt.Sprintf("cannot search article: %s", err))
		span.LogFields(log.Error(err))
		writeJsonResponse(w, http.StatusBadRequest, err)
		return
	}
	span.LogFields(log.String(" searcha articles", fmt.Sprintf("%v", articles)))
	writeJsonResponse(w, http.StatusOK, articles)

}
