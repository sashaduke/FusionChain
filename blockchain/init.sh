#!/bin/bash

Default_ACCOUNT1="shulgin"
Default_ACCOUNT2="pootis"
CHAINID="fusion_420-1"
MONIKER="fusiontestnet"
KEYRING="test"
#KEYALGO="eth_secp256k1"
LOGLEVEL="info"
# trace evm
TRACE="--trace"
# TRACE=""
FCHAIN_NODE="tcp://localhost:27657"
FCHAIN_HOME="~/.fusiond/"
FCHAIN_GAS_PRICES="1000000000nQRDO"


install_if_not_present() {
    local command_name="$1"
    if ! command -v "$command_name" > /dev/null 2>&1; then
        echo "$command_name is not installed. We will need root rights to install it"
        if [[ "$OSTYPE" == "linux-gnu"* ]]; then
            # Assume Ubuntu/Debian for simplicity. Adjust for other distros if needed.
            sudo apt update && sudo apt install -y "$command_name"
        elif [[ "$OSTYPE" == "darwin"* ]]; then
            # macOS //not tested
            brew install "$command_name"
        else
            echo "Unsupported OS for auto-install. Install $command_name manually. More info: https://example.com/$command_name/download/"
            exit 1
        fi
    fi
}

function ssed {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        gsed "$@"
    else
        sed "$@"
    fi
}

# Check and install jq
install_if_not_present "jq"

# Check and install gsed on macOS // to be tested
if [[ "$OSTYPE" == "darwin"* ]]; then
    install_if_not_present "gsed"
fi

# validate dependencies are installed
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }
if [[ "$OSTYPE" == "darwin"* ]]; then
    command -v gsed > /dev/null 2>&1 || { echo >&2 "gnu-sed not installed. Run 'brew install gnu-sed' to fix this."; exit 1; }
fi

# remove existing daemon and client
rm -rf ~/.fusiond*

set -e

make install

fusiond config keyring-backend $KEYRING
fusiond config chain-id $CHAINID

#fusiond keys add $Default_ACCOUNT1 --keyring-backend $KEYRING --algo $KEYALGO
echo "exclude try nephew main caught favorite tone degree lottery device tissue tent ugly mouse pelican gasp lava flush pen river noise remind balcony emerge" | fusiond keys add $Default_ACCOUNT1 --recover
echo "maximum fold demand spend gauge describe expect end grain entry glow purse enlist chronic robust shy panic arrange eye retreat video chat sense rare" | fusiond keys add $Default_ACCOUNT2 --recover
# echo "wool bind jeans erase promote ten session sleep logic brick drift moral twist assume people action donor guard govern three disagree share clinic oppose" | fusiond keys add $K3 --recover

# Set moniker and chain-id
fusiond init $MONIKER --chain-id $CHAINID

# Change parameter token denominations to nQRDO
cat $HOME/.fusiond/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="nQRDO"' > $HOME/.fusiond/config/tmp_genesis.json && mv $HOME/.fusiond/config/tmp_genesis.json $HOME/.fusiond/config/genesis.json
cat $HOME/.fusiond/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="nQRDO"' > $HOME/.fusiond/config/tmp_genesis.json && mv $HOME/.fusiond/config/tmp_genesis.json $HOME/.fusiond/config/genesis.json
cat $HOME/.fusiond/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="nQRDO"' > $HOME/.fusiond/config/tmp_genesis.json && mv $HOME/.fusiond/config/tmp_genesis.json $HOME/.fusiond/config/genesis.json
cat $HOME/.fusiond/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="nQRDO"' > $HOME/.fusiond/config/tmp_genesis.json && mv $HOME/.fusiond/config/tmp_genesis.json $HOME/.fusiond/config/genesis.json

# Initialize state for testing
# - First keyring, owned by shulgin
cat $HOME/.fusiond/config/genesis.json | jq '.app_state["identity"]["keyrings"]=[{"id": 0, "creator": "qredo1d652c9nngq5cneak2whyaqa4g9ehr8psyl0t7j", "description": "ACME corp", admins: ["qredo1d652c9nngq5cneak2whyaqa4g9ehr8psyl0t7j"], parties: ["qredo1d652c9nngq5cneak2whyaqa4g9ehr8psyl0t7j"]}]' > $HOME/.fusiond/config/tmp_genesis.json && mv $HOME/.fusiond/config/tmp_genesis.json $HOME/.fusiond/config/genesis.json
cat $HOME/.fusiond/config/genesis.json | jq '.app_state["identity"]["workspaces"]=[{"address": "qredoworkspace14a2hpadpsy9h5m6us54", "creator": "qredo1d652c9nngq5cneak2whyaqa4g9ehr8psyl0t7j", "owners": ["qredo1d652c9nngq5cneak2whyaqa4g9ehr8psyl0t7j"], address: "qredoworkspace14a2hpadpsy9h5m6us54"}]' > $HOME/.fusiond/config/tmp_genesis.json && mv $HOME/.fusiond/config/tmp_genesis.json $HOME/.fusiond/config/genesis.json

# Set gas limit and base fee in genesis
cat $HOME/.fusiond/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="20000000"' > $HOME/.fusiond/config/tmp_genesis.json && mv $HOME/.fusiond/config/tmp_genesis.json $HOME/.fusiond/config/genesis.json
cat $HOME/.fusiond/config/genesis.json | jq '.app_state["feemarket"]["params"]["base_fee"]="50"' > $HOME/.fusiond/config/tmp_genesis.json && mv $HOME/.fusiond/config/tmp_genesis.json $HOME/.fusiond/config/genesis.json



# disable produce empty block and enable prometheus metrics
ssed -i 's/create_empty_blocks = true/create_empty_blocks = false/g' $HOME/.fusiond/config/config.toml
ssed -i 's/prometheus = false/prometheus = true/' $HOME/.fusiond/config/config.toml
ssed -i 's/prometheus-retention-time  = "0"/prometheus-retention-time  = "1000000000000"/g' $HOME/.fusiond/config/app.toml
ssed -i 's/enabled = false/enabled = true/g' $HOME/.fusiond/config/app.toml

# enable Tendermint RPC cors from localhost
ssed -i 's/cors_allowed_origins = \[\]/cors_allowed_origins = ["http:\/\/localhost:5173"]/g' $HOME/.fusiond/config/config.toml

if [[ $1 == "pending" ]]; then
  echo "pending mode is on, please wait for the first block committed."
  ssed -i 's/create_empty_blocks_interval = "0s"/create_empty_blocks_interval = "30s"/g' $HOME/.fusiond/config/config.toml
  ssed -i 's/timeout_propose = "3s"/timeout_propose = "30s"/g' $HOME/.fusiond/config/config.toml
  ssed -i 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "5s"/g' $HOME/.fusiond/config/config.toml
  ssed -i 's/timeout_prevote = "1s"/timeout_prevote = "10s"/g' $HOME/.fusiond/config/config.toml
  ssed -i 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "5s"/g' $HOME/.fusiond/config/config.toml
  ssed -i 's/timeout_precommit = "1s"/timeout_precommit = "10s"/g' $HOME/.fusiond/config/config.toml
  ssed -i 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "5s"/g' $HOME/.fusiond/config/config.toml
  ssed -i 's/timeout_commit = "5s"/timeout_commit = "150s"/g' $HOME/.fusiond/config/config.toml
  ssed -i 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "150s"/g' $HOME/.fusiond/config/config.toml
fi

# Allocate genesis accounts (cosmos formatted addresses)
fusiond add-genesis-account $Default_ACCOUNT1 100000000000000000000000000nQRDO --keyring-backend $KEYRING
fusiond add-genesis-account $Default_ACCOUNT2 100000000000000000000000000nQRDO --keyring-backend $KEYRING

# Sign genesis transaction
fusiond gentx $Default_ACCOUNT1 1000000000000000000000nQRDO --keyring-backend $KEYRING --chain-id $CHAINID #--output-document $HOME/.fusiond/config/gentx-sk1.json
# fusiond gentx $Default_ACCOUNT2 1000000000000000000000nQRDO --keyring-backend $KEYRING --chain-id $CHAINID --output-document $HOME/.fusiond/config/gentx/gentx-sk2.json

# Collect genesis tx
fusiond collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
fusiond validate-genesis

if [[ $1 != "--defaultports" ]]; then
# Change ports to avoid conflicts with other chains running locally
  ssed -i 's/26656/27656/g' $HOME/.fusiond/config/config.toml
  ssed -i 's/26657/27657/g' $HOME/.fusiond/config/config.toml
  ssed -i 's/6060/6760/g' $HOME/.fusiond/config/config.toml
  ssed -i 's/1317/1717/g' $HOME/.fusiond/config/app.toml
  ssed -i 's/9090/9790/g' $HOME/.fusiond/config/app.toml
  ssed -i 's/9091/9791/g' $HOME/.fusiond/config/app.toml
fi



# Creating of the fchain alias
if [[ "$SHELL" == *"/bash" ]]; then
    # Bash
    echo "alias fchain=\"fusiond --node $FCHAIN_NODE --home $FCHAIN_HOME --from $Default_ACCOUNT1 --gas-prices $FCHAIN_GAS_PRICES\"" >> ~/.bashrc
    source  ~/.bashrc 
    echo "fchain alias is added"
elif [[ "$SHELL" == *"/zsh" ]]; then
    # Zsh // not tested
    config_file="~/.zshrc"
    echo "alias fchain="fusiond --node $FCHAIN_NODE --home $FCHAIN_HOME --from $Default_ACCOUNT1 --gas-prices $FCHAIN_GAS_PRICES"'" >> ~/.zshrc
    source ~/.zshrc
    echo "fchain alias is added"
else
    echo "The shell is not compatible with this script. fchain alias to be added manually."
fi



# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
fusiond start --pruning=nothing --evm.tracer=json $TRACE --log_level $LOGLEVEL --minimum-gas-prices=0.0001nQRDO --json-rpc.api eth,txpool,personal,net,debug,web3,miner --api.enable  --api.enabled-unsafe-cors
