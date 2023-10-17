package keeper_test

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/qredo/fusionchain/testutil/keeper"
	"github.com/qredo/fusionchain/x/identity/keeper"
	"github.com/qredo/fusionchain/x/identity/types"
)

func Test_msgServer_NewKeyring(t *testing.T) {
	type args struct {
		msg *types.MsgNewKeyring
	}
	tests := []struct {
		name        string
		args        args
		want        *types.MsgNewKeyringResponse
		wantCreated *types.Keyring
		wantErr     bool
	}{
		{
			name: "create a keyring",
			args: args{
				msg: types.NewMsgNewKeyring("testCreator", "testDescription"),
			},
			want: &types.MsgNewKeyringResponse{Id: 1},
			wantCreated: &types.Keyring{
				Id:          1,
				Creator:     "testCreator",
				Description: "testDescription",
				Admins:      []string{"testCreator"},
				Parties:     nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ik, ctx := keepertest.IdentityKeeper(t)
			goCtx := sdk.WrapSDKContext(ctx)
			msgSer := keeper.NewMsgServerImpl(*ik)
			got, err := msgSer.NewKeyring(goCtx, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewKeyring() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKeyring() got = %v, want %v", got, tt.want)
			}
			gotKeyring, f := ik.KeyringsRepo().Get(ctx, got.Id)
			if !f {
				t.Errorf("NewKeyring() keyring not found")
				return
			}

			if !reflect.DeepEqual(gotKeyring, tt.wantCreated) {
				t.Errorf("NewKeyring() got = %v, want %v", gotKeyring, tt.wantCreated)
				return
			}
		})
	}
}
