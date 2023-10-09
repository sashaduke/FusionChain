package keeper

import (
	"encoding/json"

	"github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	qassets "github.com/qredo/fusionchain/x/qassets/keeper"
)

// NewDefaultMessageHandler constructor
func NewDefaultMessageHandler(
	router keeper.MessageRouter,
	ics4Wrapper types.ICS4Wrapper,
	channelKeeper types.ChannelKeeper,
	capabilityKeeper types.CapabilityKeeper,
	bankKeeper types.Burner,
	qassetsKeeper qassets.Keeper,
	unpacker codectypes.AnyUnpacker,
	portSource types.ICS20TransferPortSource,
	customEncoders ...*keeper.MessageEncoders,
) keeper.Messenger {
	encoders := keeper.DefaultEncoders(unpacker, portSource)
	for _, e := range customEncoders {
		encoders = encoders.Merge(e)
	}
	return keeper.NewMessageHandlerChain(
		keeper.NewSDKMessageHandler(router, encoders),
		keeper.NewIBCRawPacketHandler(ics4Wrapper, channelKeeper, capabilityKeeper),
		keeper.NewBurnCoinMessageHandler(bankKeeper),
		NewQAssetMintMessageHandler(qassetsKeeper),
		NewQAssetBurnMessageHandler(qassetsKeeper),
	)
}

type MsgMint struct {
	Creator               string `json:"creator"`
	FromWalletId          uint64 `json:"from_wallet_id"`
	ToWorkspaceWalletAddr string `json:"to_workspace_wallet_addr"`
	IsToken               bool   `json:"is_token"`
	TokenName             string `json:"token_name"`
	TokenContractAddr     string `json:"token_contract_addr"`
	Amount                uint64 `json:"amount"`
}
type MsgBurn struct {
	Creator                 string `json:"creator"`
	FromWorkspaceWalletAddr string `json:"from_workspace_wallet_addr"`
	ToWalletId              uint64 `json:"to_wallet_id"`
	IsToken                 bool   `json:"is_token"`
	TokenName               string `json:"token_name"`
	TokenContractAddr       string `json:"token_contract_addr"`
	Amount                  uint64 `json:"amount"`
}

func NewQAssetMintMessageHandler(k qassets.Keeper) keeper.MessageHandlerFunc {
	return func(ctx sdk.Context, contractAddr sdk.AccAddress, _ string, m wasmvmtypes.CosmosMsg) (events []sdk.Event, data [][]byte, err error) {
		var msg MsgMint
		if err := json.Unmarshal(m.Custom, &msg); err != nil {
			return nil, nil, InvalidRequest{Kind: "could not deserialise QAssetMsg"}
		}
		k.Mint(ctx, msg.Creator, msg.FromWalletId, msg.ToWorkspaceWalletAddr, msg.IsToken, msg.TokenName, msg.TokenContractAddr, msg.Amount)
		return nil, nil, nil
	}
}

func NewQAssetBurnMessageHandler(k qassets.Keeper) keeper.MessageHandlerFunc {
	return func(ctx sdk.Context, contractAddr sdk.AccAddress, _ string, m wasmvmtypes.CosmosMsg) (events []sdk.Event, data [][]byte, err error) {
		var msg MsgBurn
		if err := json.Unmarshal(m.Custom, &msg); err != nil {
			return nil, nil, InvalidRequest{Kind: "could not deserialise QAssetMsg"}
		}
		k.Burn(ctx, msg.Creator, msg.FromWorkspaceWalletAddr, msg.ToWalletId, msg.IsToken, msg.TokenName, msg.TokenContractAddr, msg.Amount)
		return nil, nil, nil
	}
}
