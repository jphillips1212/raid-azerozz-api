package firestore

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type Client struct {
	Client *firestore.Client
	Ctx    *context.Context
}

// New return a new firestore client for accessing the firestore db
func New() *Client {
	// Use the application default credentials
	ctx := context.Background()
	opt := option.WithCredentialsFile("./ServiceAccountKey.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	return &Client{
		client,
		&ctx,
	}
}
