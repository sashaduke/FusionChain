package client

import (
	"google.golang.org/grpc"
	insecurecreds "google.golang.org/grpc/credentials/insecure"
)

type QueryClient struct {
	*AuthQueryClient
	*TreasuryQueryClient
	*CometClient
}

func NewQueryClient(url string, insecure bool) (*QueryClient, error) {
	var opts []grpc.DialOption
	if insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecurecreds.NewCredentials()))
	}
	grpcConn, err := grpc.Dial(url, opts...)
	if err != nil {
		return nil, err
	}

	qc, err := NewQueryClientWithConn(grpcConn)
	if err != nil {
		return nil, err
	}

	return qc, nil
}

func NewQueryClientWithConn(c *grpc.ClientConn) (*QueryClient, error) {
	comet, err := NewCometClient(c.Target())
	if err != nil {
		return nil, err
	}
	return &QueryClient{
		AuthQueryClient:     NewAuthQueryClient(c),
		TreasuryQueryClient: NewTreasuryQueryClient(c),
		CometClient:         comet,
	}, nil
}
