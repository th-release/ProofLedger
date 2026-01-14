package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:    DefaultParams(),
		EntityMap: []Entity{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	entityIndexMap := make(map[string]struct{})

	for _, elem := range gs.EntityMap {
		index := fmt.Sprint(elem.Clid)
		if _, ok := entityIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for entity")
		}
		entityIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}
