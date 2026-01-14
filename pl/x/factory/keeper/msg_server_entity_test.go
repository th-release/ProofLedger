package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"pl/x/factory/keeper"
	"pl/x/factory/types"
)

func TestEntityMsgServerCreate(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator, err := f.addressCodec.BytesToString([]byte("signerAddr__________________"))
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateEntity{Creator: creator}
		_, err := srv.CreateEntity(f.ctx, expected)
		require.NoError(t, err)
	}
}
