package main

import (
	"context"

	"github.com/machinebox/graphql"
)

var clientInstance *Client

type Client struct {
	client *graphql.Client
}

func GetClientInstance() *Client {
	if clientInstance == nil {
		clientInstance = &Client{client: graphql.NewClient("http://localhost:8080/query")}
	}

	return clientInstance
}

func (c *Client) Request(req *graphql.Request, response *interface{}) error {
	if err := c.client.Run(context.Background(), req, response); err != nil {
		return err
	}
	return nil
}
