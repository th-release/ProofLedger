package factory

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"pl/x/factory/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "ListEntity",
					Use:       "list-entity",
					Short:     "List all entity",
				},
				{
					RpcMethod:      "GetEntity",
					Use:            "get-entity [id]",
					Short:          "Gets a entity",
					Alias:          []string{"show-entity"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "clid"}},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateEntity",
					Use:            "create-entity [clid] [hash] [event-time]",
					Short:          "Create a new entity",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "clid"}, {ProtoField: "hash"}, {ProtoField: "event_time"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
