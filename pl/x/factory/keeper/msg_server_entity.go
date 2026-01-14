package keeper

import (
	"context"
	"fmt"

	"pl/x/factory/types"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateEntity(ctx context.Context, msg *types.MsgCreateEntity) (*types.MsgCreateEntityResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Check if the value already exists
	ok, err := k.Entity.Has(ctx, msg.Clid)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	} else if ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var entity = types.Entity{
		Creator:   msg.Creator,
		Clid:      msg.Clid,
		Hash:      msg.Hash,
		EventTime: msg.EventTime,
	}

	if err := k.Entity.Set(ctx, entity.Clid, entity); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgCreateEntityResponse{}, nil
}
