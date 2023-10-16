package client

import (
	"testing"

	"github.com/cometbft/cometbft/rpc/client/http"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
)

func TestCometClient_Status(t *testing.T) {
	t.SkipNow()
	type fields struct {
		client *http.HTTP
	}
	cometClient, err := NewCometClient("http://localhost:26657")
	if err != nil {
		t.FailNow()
	}
	tests := []struct {
		name    string
		fields  fields
		want    *ctypes.ResultStatus
		wantErr bool
	}{
		{
			name: "get status",
			fields: fields{
				client: cometClient.client,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := CometClient{
				client: tt.fields.client,
			}
			got, err := s.Status()
			if err != nil {
				t.FailNow()
			}
			if got == nil {
				t.Errorf("Status() got = %v, want %v", got, tt.want)
			}
		})
	}
}
