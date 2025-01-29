package rpcclient

import (
	"net/rpc"
)

type Args struct {
	A, B int
}

type Reply struct {
	Result int
}

type Client struct {
	*rpc.Client
}

func NewClient(serverAddress string) (*Client, error) {
	client, err := rpc.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}
	return &Client{client}, nil
}
