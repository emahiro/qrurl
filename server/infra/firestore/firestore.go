package firestore

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
)

var client *firestore.Client

func New(ctx context.Context) error {
	conf := &firebase.Config{ProjectID: os.Getenv("GCP_PROJECT_ID")}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return err
	}

	c, err := app.Firestore(ctx)
	if err != nil {
		return err
	}
	client = c
	return nil
}

func Close() error {
	return client.Close()
}

func Add[T any](ctx context.Context, collection string, data T) error {
	_, _, err := client.Collection(collection).Add(ctx, data)
	return err
}

func GetAll[T any](ctx context.Context, collection string) ([]T, error) {
	results := make([]T, 0)
	iter := client.Collection(collection).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var result T
		if err := doc.DataTo(&result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func GetLatestOne[T any](ctx context.Context, collection string) (*T, error) {
	query := client.Collection(collection).OrderBy("CreatedAt", firestore.Desc).Limit(1)
	iter := query.Documents(ctx)
	results := make([]T, 1)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var result T
		if err := doc.DataTo(&result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return &results[0], nil
}

func Query[T any](ctx context.Context, collenction, key, op, value string) ([]T, error) {
	query := client.Collection(collenction).Where(key, op, value)
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
