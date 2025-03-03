package main

import (
	"context"

	"github.com/machinebox/graphql"
)

var clientInstance GraphQLClientInterface

type Client struct {
	client *graphql.Client
}

type GraphQLClientInterface interface {
	Request(req *graphql.Request, resp interface{}) error
}

func GetClientInstance() GraphQLClientInterface {
	if clientInstance == nil {
		clientInstance = &Client{client: graphql.NewClient("http://localhost:8080/query")}
	}

	return clientInstance
}

func (c *Client) Request(req *graphql.Request, response interface{}) error {
	if err := c.client.Run(context.Background(), req, response); err != nil {
		return err
	}
	return nil
}

func SetClientInstance(client GraphQLClientInterface) {
	clientInstance = client
}
