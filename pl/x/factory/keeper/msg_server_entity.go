package keeper

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"os"
	"strconv"
	"time"

	"pl/x/factory/types"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func generateCollisionFreeClid() string {
	now := time.Now()
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		strconv.FormatInt(now.UnixNano(), 36),     // 나노초까지 (충돌 극히 어려움)
		strconv.FormatInt(int64(os.Getpid()), 36), // 프로세스 ID
		strconv.Itoa(os.Getppid()),                // 부모 프로세스 ID
		hostnameHash(),                            // 서버/컨테이너 고유값
		randomHex(12),                             // 48bit 랜덤 (충돌 확률 10^-14 이하)
	)
}

func hostnameHash() string {
	h, _ := os.Hostname()
	sum := crc32.ChecksumIEEE([]byte(h))
	return strconv.FormatUint(uint64(sum), 36)
}

func randomHex(n int) string {
	b := make([]byte, n/2+1)
	rand.Read(b)
	return hex.EncodeToString(b)[:n]
}

func (k msgServer) CreateEntity(ctx context.Context, msg *types.MsgCreateEntity) (*types.MsgCreateEntityResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	clid := generateCollisionFreeClid()

	// Check if the value already exists
	ok, err := k.Entity.Has(ctx, clid)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	} else if ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var entity = types.Entity{
		Creator:   msg.Creator,
		Clid:      clid,
		Hash:      msg.Hash,
		EventTime: time.Now().UnixMilli(),
	}

	if err := k.Entity.Set(ctx, entity.Clid, entity); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgCreateEntityResponse{}, nil
}
