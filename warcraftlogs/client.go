package warcraftlogs

import (
	"context"

	"github.com/hasura/go-graphql-client"
	"golang.org/x/oauth2/clientcredentials"
)

type Client struct {
	Client *graphql.Client
}

// New return a new graphql client for calling warcraft logs
func New() *Client {
	conf := clientcredentials.Config{
		ClientID:     "9761af97-7dc0-4f3e-b1d4-f9744af50f2e",
		ClientSecret: "",
		TokenURL:     "",
	}

	httpClient := conf.Client(context.Background())

	return &Client{
		graphql.NewClient("https://www.warcraftlogs.com/api/v2/client", httpClient),
	}
}
