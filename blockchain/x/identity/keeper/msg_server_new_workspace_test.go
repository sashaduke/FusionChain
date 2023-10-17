package keeper_test

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/qredo/fusionchain/testutil/keeper"
	"github.com/qredo/fusionchain/x/identity/keeper"
	"github.com/qredo/fusionchain/x/identity/types"
)

func Test_msgServer_NewWorkspace(t *testing.T) {
	type args struct {
		msg *types.MsgNewWorkspace
	}
	tests := []struct {
		name    string
		args    args
		want    *types.MsgNewWorkspaceResponse
		wantErr bool
	}{
		{
			name: "create a workspace",
			args: args{
				msg: types.NewMsgNewWorkspace("test", 0, 0),
			},
			want: &types.MsgNewWorkspaceResponse{
				Address: "qredoworkspace14a2hpadpsy9h5m6us54",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ik, ctx := keepertest.IdentityKeeper(t)
			goCtx := sdk.WrapSDKContext(ctx)

			msgSer := keeper.NewMsgServerImpl(*ik)
			got, err := msgSer.NewWorkspace(goCtx, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWorkspace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWorkspace() got = %v, want %v", got, tt.want)
			}
			createdWrk := ik.GetWorkspace(ctx, got.Address)
			if createdWrk.Creator != tt.args.msg.Creator || createdWrk.AdminPolicyId != tt.args.msg.AdminPolicyId || createdWrk.SignPolicyId != tt.args.msg.SignPolicyId {
				t.Errorf("NewWorkspace() created workspace = %v, want %v", createdWrk, tt.args.msg)
			}
		})
	}
}
