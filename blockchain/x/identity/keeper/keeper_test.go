package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/x/params/types"
)

type memStoreKey struct {
	kv map[string]string
}

func newMemoryStoreKey() *memStoreKey {
	kv := make(map[string]string)
	return &memStoreKey{kv: kv}
}

func (m memStoreKey) Name() string {
	return "storeKey"
}

func (m memStoreKey) String() string {
	return "storeKey"
}

func Test_NewKeeper(t *testing.T) {
	if k := NewKeeper(nil, newMemoryStoreKey(), newMemoryStoreKey(),
		types.NewSubspace(nil, nil, newMemoryStoreKey(), newMemoryStoreKey(), "subspace"),
		nil); k == nil {
		t.Fatal("expeced non-nil struct")
	}
}
