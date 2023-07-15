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

func Add(ctx context.Context, collection string, data any) error {
	_, _, err := client.Collection(collection).Add(ctx, data)
	return err
}

func GetAll(ctx context.Context, collection string) ([]map[string]any, error) {
	result := make([]map[string]any, 0)
	iter := client.Collection(collection).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		result = append(result, doc.Data())
	}
	return result, nil
}

func Query(ctx context.Context, collenction, key, op, value string) ([]map[string]any, error) {
	query := client.Collection(collenction).Where(key, op, value)
	iter := query.Documents(ctx)
	results := make([]map[string]any, 0)
	for {
		ss, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		results = append(results, ss.Data())
	}
	return results, nil
}
