package keeper

import (
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/lavanet/lava/x/epochstorage/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace

		bankKeeper    types.BankKeeper
		accountKeeper types.AccountKeeper
		specKeeper    types.SpecKeeper

		fixationRegistries map[string]func(sdk.Context) any
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,

	bankKeeper types.BankKeeper, accountKeeper types.AccountKeeper, specKeeper types.SpecKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	keeper := &Keeper{

		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
		bankKeeper: bankKeeper, accountKeeper: accountKeeper, specKeeper: specKeeper,

		fixationRegistries: make(map[string]func(sdk.Context) any),
	}

	keeper.AddFixationRegistry(string(types.KeyEpochBlocks), func(ctx sdk.Context) any { return keeper.EpochBlocksRaw(ctx) })
	keeper.AddFixationRegistry(string(types.KeyEpochsToSave), func(ctx sdk.Context) any { return keeper.EpochsToSaveRaw(ctx) })
	keeper.AddFixationRegistry(string(types.KeyUnstakeHoldBlocks), func(ctx sdk.Context) any { return keeper.UnstakeHoldBlocksRaw(ctx) })

	return keeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) AddFixationRegistry(fixationKey string, getParamFunction func(sdk.Context) any) {

	if _, ok := k.fixationRegistries[fixationKey]; ok {
		panic(fmt.Sprintf("duplicate fixation registry %s", fixationKey))
	}
	k.fixationRegistries[fixationKey] = getParamFunction
}

func (k *Keeper) GetFixationRegistries() map[string]func(sdk.Context) any {
	return k.fixationRegistries
}
