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
		ClientID:     "",
		ClientSecret: "",
		TokenURL:     "https://www.warcraftlogs.com/oauth/token",
	}

	httpClient := conf.Client(context.Background())

	return &Client{
		graphql.NewClient("https://www.warcraftlogs.com/api/v2/client", httpClient),
	}
}
