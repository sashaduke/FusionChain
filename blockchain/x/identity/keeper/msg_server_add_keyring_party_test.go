package keeper_test

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/qredo/fusionchain/testutil/keeper"
	"github.com/qredo/fusionchain/x/identity/keeper"
	"github.com/qredo/fusionchain/x/identity/types"
)

func Test_msgServer_AddKeyringParty(t *testing.T) {

	type args struct {
		msg        *types.MsgAddKeyringParty
		msgKeyring *types.MsgNewKeyring
	}
	tests := []struct {
		name        string
		args        args
		want        *types.MsgAddKeyringPartyResponse
		wantKeyring *types.Keyring
		wantErr     bool
	}{
		{
			name: "add a party to a keyring",
			args: args{
				msgKeyring: types.NewMsgNewKeyring("testCreator", "testDescription"),
				msg:        types.NewMsgAddKeyringParty("testCreator", 1, "testParty"),
			},
			want: &types.MsgAddKeyringPartyResponse{},
			wantKeyring: &types.Keyring{
				Id:          1,
				Creator:     "testCreator",
				Description: "testDescription",
				Admins:      []string{"testCreator"},
				Parties:     []string{"testParty"},
			},
			wantErr: false,
		},
		{
			name: "keyring not found",
			args: args{
				msgKeyring: types.NewMsgNewKeyring("testCreator", "testDescription"),
				msg:        types.NewMsgAddKeyringParty("testCreator", 2, "testParty"),
			},
			want:        &types.MsgAddKeyringPartyResponse{},
			wantKeyring: nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ik, ctx := keepertest.IdentityKeeper(t)
			goCtx := sdk.WrapSDKContext(ctx)
			msgSer := keeper.NewMsgServerImpl(*ik)
			keyringRes, err := msgSer.NewKeyring(goCtx, tt.args.msgKeyring)
			if err != nil {
				t.Errorf("NewKeyring() error = %v", err)
				return
			}
			got, err := msgSer.AddKeyringParty(goCtx, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddKeyringParty() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("AddKeyringParty() got = %v, want %v", got, tt.want)
				}
				gotKeyring, f := ik.KeyringsRepo().Get(ctx, keyringRes.Id)
				if !f {
					t.Errorf("NewKeyring() keyring not found")
					return
				}
				if !reflect.DeepEqual(gotKeyring, tt.wantKeyring) {
					t.Errorf("NewKeyring() got = %v, want %v", gotKeyring, tt.wantKeyring)
					return
				}
			}
		})
	}
}
