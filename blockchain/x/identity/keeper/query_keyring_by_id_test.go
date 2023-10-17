package keeper_test

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/qredo/fusionchain/testutil/keeper"
	"github.com/qredo/fusionchain/x/identity/keeper"
	"github.com/qredo/fusionchain/x/identity/types"
)

func TestKeeper_KeyringByID(t *testing.T) {

	type args struct {
		msg *types.MsgNewKeyring
		req *types.QueryKeyringByIdRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *types.QueryKeyringByIdResponse
		wantErr bool
	}{
		{
			name: "get a keyring by id",
			args: args{
				msg: types.NewMsgNewKeyring("testCreator", "testDescription"),
				req: &types.QueryKeyringByIdRequest{
					Id: 1,
				},
			},
			want: &types.QueryKeyringByIdResponse{Keyring: &types.Keyring{
				Id:          1,
				Creator:     "testCreator",
				Description: "testDescription",
				Admins:      []string{"testCreator"},
				Parties:     nil,
			}},
			wantErr: false,
		},
		{
			name: "keyring by id not found",
			args: args{
				msg: types.NewMsgNewKeyring("testCreator", "testDescription"),
				req: &types.QueryKeyringByIdRequest{
					Id: 5,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ik, ctx := keepertest.IdentityKeeper(t)
			goCtx := sdk.WrapSDKContext(ctx)
			msgSer := keeper.NewMsgServerImpl(*ik)
			_, err := msgSer.NewKeyring(goCtx, tt.args.msg)
			got, err := ik.KeyringByID(goCtx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyringByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyringByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
