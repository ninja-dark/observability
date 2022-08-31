package store

import (
	"context"
	"elastic/e"
	"elastic/l"
	"elastic/m"
	"errors"
	"fmt"

	
	"github.com/mitchellh/mapstructure"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"
)

var (
	ErrNotFound = errors.New("not found")
)

type ArticleStore struct {
	E e.E
	logger *zap.Logger
	tracer opentracing.Tracer
}

func NewArticleStore(logger *zap.Logger, tracer opentracing.Tracer) (ArticleStore, error) {
	e, err := e.NewE("articles", logger, tracer)
	if err != nil {
		return ArticleStore{}, err
	}
	return ArticleStore{E: e, logger: logger, tracer: tracer}, nil
}

func (s ArticleStore) Add(ctx context.Context, article m.Article) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ArticleStore.Add")
	defer span.Finish()
	
	err := s.E.Insert(ctx, article)
	if err != nil {
		s.logger.Error(fmt.Sprintf("cannot insert: %s", err))
		span.LogFields(log.Error(err))
		return err
	}
	span.LogFields(log.String("get query", fmt.Sprintf("%v", article)))
	return nil
}

func (s ArticleStore) Search(ctx context.Context, query string) ([]m.Article, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ArticleStore.Search")
	defer span.Finish()

	result, err := s.E.Search(ctx, query)
	if err != nil {
		s.logger.Error(fmt.Sprintf("cannot search article: %s", err))
		span.LogFields(log.Error(err))
		return nil, err
	}
	hits := result.Hits.Hits
	articles := []m.Article{}
	for _, hit := range hits {
		var article m.Article
		//map[string]interface{} -> struct
		err = mapstructure.Decode(hit.Source, &article)
		if err != nil {
			s.logger.Error(fmt.Sprintf("cannot decode: %s", err))
			span.LogFields(log.Error(err))
			return nil, err
		}
		article.Id = hit.ID
		articles = append(articles, article)
	}
	span.LogFields(log.String("get query", query))
	return articles, nil
}

func (s ArticleStore) Get(ctx context.Context, id string) (m.Article, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ArticleStore.Get")
	defer span.Finish()

	result, err := s.E.Get(ctx, id)
	if err != nil {
		s.logger.Error(fmt.Sprintf("cannot get article: %s", err))
		span.LogFields(log.Error(err))
		return m.Article{}, err
	}
	l.L(result)
	var article m.Article
	// err = mapstructure.Decode(result.Source, &article)
	// if err != nil {
	// 	return m.Article{}, err
	// }
	span.LogFields(log.String("get article", id))
	return article, nil
}
