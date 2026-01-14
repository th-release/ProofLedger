package keeper_test

import (
	"strconv"
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
		expected := &types.MsgCreateEntity{Creator: creator,
			Clid: strconv.Itoa(i),
		}
		_, err := srv.CreateEntity(f.ctx, expected)
		require.NoError(t, err)
		rst, err := f.keeper.Entity.Get(f.ctx, expected.Clid)
		require.NoError(t, err)
		require.Equal(t, expected.Creator, rst.Creator)
	}
}
