package types_test

import (
	"testing"

	"github.com/lavanet/lava/x/pairing/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{

				UniquePaymentStorageClientProviderList: []types.UniquePaymentStorageClientProvider{
					{
						Index: "0",
					},
					{
						Index: "1",
					},
				},
				ClientPaymentStorageList: []types.ClientPaymentStorage{
					{
						Index: "0",
					},
					{
						Index: "1",
					},
				},
				EpochPaymentsList: []types.EpochPayments{
					{
						Index: "0",
					},
					{
						Index: "1",
					},
				},
				FixatedServicersToPairList: []types.FixatedServicersToPair{
					{
						Index: "0",
					},
					{
						Index: "1",
					},
				},
				FixatedStakeToMaxCuList: []types.FixatedStakeToMaxCu{
					{
						Index: "0",
					},
					{
						Index: "1",
					},
				},
				FixatedEpochBlocksOverlapList: []types.FixatedEpochBlocksOverlap{
					{
						Index: "0",
					},
					{
						Index: "1",
					},
				},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated uniquePaymentStorageClientProvider",
			genState: &types.GenesisState{
				UniquePaymentStorageClientProviderList: []types.UniquePaymentStorageClientProvider{
					{
						Index: "0",
					},
					{
						Index: "0",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated clientPaymentStorage",
			genState: &types.GenesisState{
				ClientPaymentStorageList: []types.ClientPaymentStorage{
					{
						Index: "0",
					},
					{
						Index: "0",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated epochPayments",
			genState: &types.GenesisState{
				EpochPaymentsList: []types.EpochPayments{
					{
						Index: "0",
					},
					{
						Index: "0",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated fixatedServicersToPair",
			genState: &types.GenesisState{
				FixatedServicersToPairList: []types.FixatedServicersToPair{
					{
						Index: "0",
					},
					{
						Index: "0",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated fixatedStakeToMaxCu",
			genState: &types.GenesisState{
				FixatedStakeToMaxCuList: []types.FixatedStakeToMaxCu{
					{
						Index: "0",
					},
					{
						Index: "0",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated fixatedEpochBlocksOverlap",
			genState: &types.GenesisState{
				FixatedEpochBlocksOverlapList: []types.FixatedEpochBlocksOverlap{
					{
						Index: "0",
					},
					{
						Index: "0",
					},
				},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
