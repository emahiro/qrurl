package firestore

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/cockroachdb/errors"
	"google.golang.org/api/iterator"
)

var client *firestore.Client

func New(ctx context.Context) error {
	conf := &firebase.Config{ProjectID: os.Getenv("GCP_PROJECT_ID")}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return errors.WithHint(err, "failed to initialize firebase app")
	}

	c, err := app.Firestore(ctx)
	if err != nil {
		return errors.WithHint(err, "failed to initialize firestore client")
	}
	client = c
	return nil
}

func Client() *firestore.Client {
	return client
}

func Close() error {
	return client.Close()
}

func Add[T any](ctx context.Context, collection string, data T) error {
	_, _, err := client.Collection(collection).Add(ctx, data)
	return err
}

type QueryOption struct {
	Key     string
	Op      string
	Value   any
	OrderBy string
	Desc    bool
	Limit   int
}

func Query[T any](ctx context.Context, collection string, opt QueryOption) ([]T, error) {
	limit := opt.Limit
	if limit == 0 {
		limit = 1000
	}

	query := client.Collection(collection).Limit(limit)

	if opt.Key != "" && opt.Op != "" && opt.Value != "" {
		query = query.Where(opt.Key, opt.Op, opt.Value)
	}

	if opt.OrderBy != "" {
		if opt.Desc {
			query = query.OrderBy(opt.OrderBy, firestore.Desc)
		} else {
			query = query.OrderBy(opt.OrderBy, firestore.Asc)
		}
	}

	iter := query.Documents(ctx)
	results := make([]T, 0)
	for {
		ss, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var result T
		if err := ss.DataTo(&result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}
