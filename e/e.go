package e

import (
	"bytes"
	"context"
	"elastic/l"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
)

type E struct {
	C         *elasticsearch.Client
	IndexName string
}

type I interface{}
type M map[string]I

func NewE(indexName string) (E, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://elasticsearch:9200"},
	})
	if err != nil {
		return E{}, err
	}
	_, err = client.Ping()
	if err != nil {
		return E{}, err
	}

	return E{
		C:         client,
		IndexName: indexName,
	}, nil
}

func (e E) Info() (M, error) {
	res, err := e.C.Info()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var r M
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	return r, nil
}

func (e E) Insert(ctx context.Context, i I) error {
	data, err := json.Marshal(i)
	if err != nil {
		return err
	}
	id := e.GetId(i)
	req := esapi.IndexRequest{
		Index:      e.IndexName,
		DocumentID: id,
		Refresh:    "true",
		Body:       bytes.NewBuffer(data),
	}
	res, err := req.Do(ctx, e.C)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	l.L(res.Status())
	return nil
}
func (e E) Search(ctx context.Context, q string) (SearchResponse, error) {
	var r SearchResponse
	var buf bytes.Buffer
	query := M{
		"query": M{
			"match": M{
				"title": q,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return r, err
	}

	// Perform the search request.
	res, err := e.C.Search(
		e.C.Search.WithContext(context.Background()),
		e.C.Search.WithIndex(e.IndexName),
		e.C.Search.WithBody(&buf),
		e.C.Search.WithTrackTotalHits(true),
		e.C.Search.WithPretty(),
	)
	if err != nil {
		return r, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return r, err
	}
	return r, nil
}
func (e E) Get(ctx context.Context, id string) (M, error) {
	req := esapi.GetRequest{
		Index:      e.IndexName,
		DocumentID: id,
	}
	res, err := req.Do(ctx, e.C)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	l.L(res.Status())
	var r M
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	return r, nil
}
func (e E) GetId(i I) string {
	typeOf := reflect.TypeOf(i)
	if typeOf.Kind() == reflect.Map {
		m, ok := i.(M)
		if ok {
			id, ok := m["id"]
			if ok {
				return rawId(id)
			}
		}
		m1, ok := i.(map[string]interface{})
		if ok {
			id, ok := m1["id"]
			if ok {
				return rawId(id)
			}
		}
	}
	if typeOf.Kind() == reflect.Struct {
		v := reflect.ValueOf(i)
		for i := 0; i < v.NumField(); i++ {
			field := typeOf.Field(i)
			tag, ok := field.Tag.Lookup("json")
			if ok {
				if tag == "id" {
					idField := v.Field(i)
					idType := idField.Type()
					if idType.Kind() == reflect.String {
						return rawId(idField.String())
					}
					if idType.Kind() == reflect.Int {
						return rawId(idField.Int())
					}
					return uuid.New().String()
				}
			}
		}
	}
	return uuid.New().String()
}
func rawId(i interface{}) string {
	id, ok := i.(string)
	if ok {
		if id != "" {
			return id
		}
	}
	idInt, ok := i.(int)
	if ok {
		if idInt > 0 {
			return fmt.Sprintf("%v", idInt)
		}
	}
	return uuid.New().String()
}

type SearchResponse struct {
	Shards   Shards `json:"_shards"`
	Hits     Hits   `json:"hits"`
	TimedOut bool   `json:"timed_out"`
	Took     int64  `json:"took"`
}

type Hits struct {
	Hits     []Hit   `json:"hits"`
	MaxScore float64 `json:"max_score"`
	Total    Total   `json:"total"`
}

type Hit struct {
	ID     string  `json:"_id"`
	Index  string  `json:"_index"`
	Score  float64 `json:"_score"`
	Source M       `json:"_source"`
}

type Source struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Total struct {
	Relation string `json:"relation"`
	Value    int64  `json:"value"`
}

type Shards struct {
	Failed     int64 `json:"failed"`
	Skipped    int64 `json:"skipped"`
	Successful int64 `json:"successful"`
	Total      int64 `json:"total"`
}
