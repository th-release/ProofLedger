package factory

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"pl/testutil/sample"
	factorysimulation "pl/x/factory/simulation"
	"pl/x/factory/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	factoryGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		EntityMap: []types.Entity{{Creator: sample.AccAddress(),
			Clid: "0",
		}, {Creator: sample.AccAddress(),
			Clid: "1",
		}}}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&factoryGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateEntity          = "op_weight_msg_factory"
		defaultWeightMsgCreateEntity int = 100
	)

	var weightMsgCreateEntity int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateEntity, &weightMsgCreateEntity, nil,
		func(_ *rand.Rand) {
			weightMsgCreateEntity = defaultWeightMsgCreateEntity
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateEntity,
		factorysimulation.SimulateMsgCreateEntity(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
