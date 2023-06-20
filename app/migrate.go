package app

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MigrateFn func() error

func CreateStoreLoaderWithMigrateFn(migrateHeight int64, migrateFn MigrateFn) baseapp.StoreLoader {
	return func(ms sdk.CommitMultiStore) error {
		if migrateHeight == ms.LastCommitID().Version+1 {
			err := migrateFn()
			if err != nil {
				panic(err)
			}
		}

		return baseapp.DefaultStoreLoader(ms)
	}
}
