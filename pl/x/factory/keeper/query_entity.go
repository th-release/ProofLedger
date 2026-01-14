package keeper

import (
	"context"
	"errors"

	"pl/x/factory/types"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListEntity(ctx context.Context, req *types.QueryAllEntityRequest) (*types.QueryAllEntityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	entitys, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Entity,
		req.Pagination,
		func(_ string, value types.Entity) (types.Entity, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllEntityResponse{Entity: entitys, Pagination: pageRes}, nil
}

func (q queryServer) GetEntity(ctx context.Context, req *types.QueryGetEntityRequest) (*types.QueryGetEntityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.Entity.Get(ctx, req.Clid)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetEntityResponse{Entity: val}, nil
}
