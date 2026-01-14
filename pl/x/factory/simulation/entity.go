package simulation

import (
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"pl/x/factory/keeper"
	"pl/x/factory/types"
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

func SimulateMsgCreateEntity(
	ak types.AuthKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
	txGen client.TxConfig,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		msg := &types.MsgCreateEntity{
			Creator: simAccount.Address.String()}

		clid := generateCollisionFreeClid()

		found, err := k.Entity.Has(ctx, clid)
		if err == nil && found {
			return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "Entity already exist"), nil, nil
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           txGen,
			Cdc:             nil,
			Msg:             msg,
			Context:         ctx,
			SimAccount:      simAccount,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(),
			AccountKeeper:   ak,
			Bankkeeper:      bk,
		}
		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}
