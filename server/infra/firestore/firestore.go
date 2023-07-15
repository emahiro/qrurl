package firestore

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
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
