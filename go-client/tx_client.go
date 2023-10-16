package client

import (
	"google.golang.org/grpc"
)

type TxClient struct {
	*RawTxClient
	*TreasuryTxClient
	*CometClient
}

func NewTxClient(id Identity, chainID string, c *grpc.ClientConn, accountFetcher AccountFetcher, cometClient *CometClient) *TxClient {
	raw := NewRawTxClient(id, chainID, c, accountFetcher)
	return &TxClient{
		RawTxClient:      raw,
		TreasuryTxClient: NewTreasuryTxClient(raw),
		CometClient:      cometClient,
	}
}
