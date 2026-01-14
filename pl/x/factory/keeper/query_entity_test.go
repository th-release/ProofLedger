package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"pl/x/factory/keeper"
	"pl/x/factory/types"
)

func createNEntity(keeper keeper.Keeper, ctx context.Context, n int) []types.Entity {
	items := make([]types.Entity, n)
	for i := range items {
		items[i].Clid = strconv.Itoa(i)
		items[i].Hash = strconv.Itoa(i)
		items[i].EventTime = int64(i)
		_ = keeper.Entity.Set(ctx, items[i].Clid, items[i])
	}
	return items
}

func TestEntityQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNEntity(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetEntityRequest
		response *types.QueryGetEntityResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetEntityRequest{
				Clid: msgs[0].Clid,
			},
			response: &types.QueryGetEntityResponse{Entity: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetEntityRequest{
				Clid: msgs[1].Clid,
			},
			response: &types.QueryGetEntityResponse{Entity: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetEntityRequest{
				Clid: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.GetEntity(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestEntityQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNEntity(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllEntityRequest {
		return &types.QueryAllEntityRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListEntity(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Entity), step)
			require.Subset(t, msgs, resp.Entity)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListEntity(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Entity), step)
			require.Subset(t, msgs, resp.Entity)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListEntity(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.Entity)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListEntity(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
