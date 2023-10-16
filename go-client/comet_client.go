package client

import (
	"context"

	"github.com/cometbft/cometbft/rpc/client/http"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
)

type CometClient struct {
	client *http.HTTP
}

func NewCometClient(url string) (*CometClient, error) {
	hc, err := http.New(url, "/websocket")
	if err != nil {
		return nil, err
	}
	return &CometClient{
		client: hc,
	}, nil
}

func (s CometClient) Status() (*ctypes.ResultStatus, error) {
	return s.client.Status(context.Background())
}

func (s CometClient) Health() (*ctypes.ResultHealth, error) {
	return s.client.Health(context.Background())
}
