package types

import (
	stdmath "math"
	"math/rand"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// func TestNextInflation(t *testing.T) {
// 	minter := DefaultInitialMinter()
// 	params := DefaultParams()
// 	blocksPerYr := sdk.NewDec(int64(params.BlocksPerYear))

// 	// Governing Mechanism:
// 	//    inflationRateChangePerYear = (1- BondedRatio/ GoalBonded) * MaxInflationRateChange

// 	tests := []struct {
// 		bondedRatio, setInflation, expChange sdk.Dec
// 	}{
// 		// with 0% bonded atom supply the inflation should increase by InflationRateChange
// 		{sdk.ZeroDec(), sdk.NewDecWithPrec(7, 2), params.InflationRateChange.Quo(blocksPerYr)},

// 		// 100% bonded, starting at 20% inflation and being reduced
// 		// (1 - (1/0.67))*(0.13/8667)
// 		{
// 			sdk.OneDec(), sdk.NewDecWithPrec(20, 2),
// 			sdk.OneDec().Sub(sdk.OneDec().Quo(params.GoalBonded)).Mul(params.InflationRateChange).Quo(blocksPerYr),
// 		},

// 		// 50% bonded, starting at 10% inflation and being increased
// 		{
// 			sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(10, 2),
// 			sdk.OneDec().Sub(sdk.NewDecWithPrec(5, 1).Quo(params.GoalBonded)).Mul(params.InflationRateChange).Quo(blocksPerYr),
// 		},

// 		// test 7% minimum stop (testing with 100% bonded)
// 		{sdk.OneDec(), sdk.NewDecWithPrec(7, 2), sdk.ZeroDec()},
// 		{sdk.OneDec(), sdk.NewDecWithPrec(700000001, 10), sdk.NewDecWithPrec(-1, 10)},

// 		// test 20% maximum stop (testing with 0% bonded)
// 		{sdk.ZeroDec(), sdk.NewDecWithPrec(20, 2), sdk.ZeroDec()},
// 		{sdk.ZeroDec(), sdk.NewDecWithPrec(1999999999, 10), sdk.NewDecWithPrec(1, 10)},

// 		// perfect balance shouldn't change inflation
// 		{sdk.NewDecWithPrec(67, 2), sdk.NewDecWithPrec(15, 2), sdk.ZeroDec()},
// 	}
// 	for i, tc := range tests {
// 		minter.Inflation = tc.setInflation

// 		inflation := minter.NextInflationRate(params, tc.bondedRatio, )
// 		diffInflation := inflation.Sub(tc.setInflation)

// 		require.True(t, diffInflation.Equal(tc.expChange),
// 			"Test Index: %v\nDiff:  %v\nExpected: %v\n", i, diffInflation, tc.expChange)
// 	}
// }

// func TestBlockProvision(t *testing.T) {
// 	minter := InitialMinter(sdk.NewDecWithPrec(1, 1))
// 	params := DefaultParams()

// 	secondsPerYear := int64(60 * 60 * 8766)

// 	tests := []struct {
// 		annualProvisions int64
// 		expProvisions    int64
// 	}{
// 		{secondsPerYear / 5, 1},
// 		{secondsPerYear/5 + 1, 1},
// 		{(secondsPerYear / 5) * 2, 2},
// 		{(secondsPerYear / 5) / 2, 0},
// 	}
// 	for i, tc := range tests {
// 		minter.AnnualProvisions = sdk.NewDec(tc.annualProvisions)
// 		provisions := minter.BlockProvision(params)

// 		expProvisions := sdk.NewCoin(params.MintDenom,
// 			sdk.NewInt(tc.expProvisions))

// 		require.True(t, expProvisions.IsEqual(provisions),
// 			"test: %v\n\tExp: %v\n\tGot: %v\n",
// 			i, tc.expProvisions, provisions)
// 	}
// }

// // Benchmarking :)
// // previously using sdk.Int operations:
// // BenchmarkBlockProvision-4 5000000 220 ns/op
// //
// // using sdk.Dec operations: (current implementation)
// // BenchmarkBlockProvision-4 3000000 429 ns/op
// func BenchmarkBlockProvision(b *testing.B) {
// 	b.ReportAllocs()
// 	minter := InitialMinter(sdk.NewDecWithPrec(1, 1))
// 	params := DefaultParams()

// 	s1 := rand.NewSource(100)
// 	r1 := rand.New(s1)
// 	minter.AnnualProvisions = sdk.NewDec(r1.Int63n(1000000))

// 	// run the BlockProvision function b.N times
// 	for n := 0; n < b.N; n++ {
// 		minter.BlockProvision(params)
// 	}
// }

// // Next inflation benchmarking
// // BenchmarkNextInflation-4 1000000 1828 ns/op
// func BenchmarkNextInflation(b *testing.B) {
// 	b.ReportAllocs()
// 	minter := InitialMinter(sdk.NewDecWithPrec(1, 1))
// 	params := DefaultParams()
// 	bondedRatio := sdk.NewDecWithPrec(1, 1)

// 	// run the NextInflationRate function b.N times
// 	for n := 0; n < b.N; n++ {
// 		minter.NextInflationRate(params, bondedRatio)
// 	}
// }

// // Next annual provisions benchmarking
// // BenchmarkNextAnnualProvisions-4 5000000 251 ns/op
// func BenchmarkNextAnnualProvisions(b *testing.B) {
// 	b.ReportAllocs()
// 	minter := InitialMinter(sdk.NewDecWithPrec(1, 1))
// 	params := DefaultParams()
// 	totalSupply := sdk.NewInt(100000000000000)

// 	// run the NextAnnualProvisions function b.N times
// 	for n := 0; n < b.N; n++ {
// 		minter.NextAnnualProvisions(params, totalSupply)
// 	}
// }

func TestSimulateMint(t *testing.T) {
	minter := DefaultInitialMinter()
	params := DefaultParams()
	totalSupply := math.NewInt(1_000_000_000_000_000_000)
	totalStaked := math.ZeroInt()
	tokenMinted := sdk.NewCoin("stake", math.ZeroInt())

	// Ensure BlocksPerYear can be safely converted to int
	if params.BlocksPerYear > uint64(stdmath.MaxInt64) {
		t.Fatal("blocks per year exceeds maximum int64 value")
	}

	for i := uint64(1); i <= params.BlocksPerYear; i++ {
		stakingDiff := sdk.NewDec(int64(rand.Intn(10))).QuoInt(math.NewInt(1_000_000)).MulInt(totalSupply)
		if (rand.Float32() > 0.5 || totalStaked.Add(stakingDiff.RoundInt()).GT(totalSupply)) && !totalStaked.Sub(stakingDiff.RoundInt()).IsNegative() {
			stakingDiff = stakingDiff.Neg()
		}
		totalStaked = totalStaked.Add(stakingDiff.RoundInt())
		bondedRatio := sdk.NewDecFromInt(totalStaked).Quo(sdk.NewDecFromInt(totalSupply))

		newInflation, err := minter.NextInflationRate(params, bondedRatio, totalStaked)
		require.NoError(t, err)
		minter.Inflation = newInflation

		newAnnualProvisions, err := minter.NextAnnualProvisions(params, totalStaked)
		require.NoError(t, err)
		minter.AnnualProvisions = newAnnualProvisions

		// mint coins, update supply
		mintedCoin, err := minter.BlockProvision(params)
		require.NoError(t, err)
		tokenMinted = tokenMinted.Add(mintedCoin)
	}
	require.True(t, params.MaxTokenPerYear.GTE(tokenMinted.Amount))
	require.True(t, params.MinTokenPerYear.LTE(tokenMinted.Amount))
}

func TestMinterNextAnnualProvisions(t *testing.T) {
	params := DefaultParams()

	// Ensure BlocksPerYear can be safely converted to int
	if params.BlocksPerYear > uint64(stdmath.MaxInt64) {
		t.Fatal("blocks per year exceeds maximum int64 value")
	}
	for i := uint64(1); i <= params.BlocksPerYear; i++ {
		// ... existing code ...
	}
}

func TestNextInflationRateEdgeCases(t *testing.T) {
	minter := DefaultInitialMinter()
	params := DefaultParams()

	tests := []struct {
		name           string
		bondedRatio    sdk.Dec
		stakingSupply  math.Int
		expectError    bool
		errorString    string
		expectedResult sdk.Dec
	}{
		{
			name:           "nil bonded ratio",
			bondedRatio:    sdk.Dec{},
			stakingSupply:  math.NewInt(1000000),
			expectError:    true,
			errorString:    "bonded ratio cannot be nil",
			expectedResult: sdk.Dec{},
		},
		{
			name:           "negative bonded ratio",
			bondedRatio:    sdk.NewDec(-1),
			stakingSupply:  math.NewInt(1000000),
			expectError:    true,
			errorString:    "bonded ratio cannot be negative",
			expectedResult: sdk.Dec{},
		},
		{
			name:           "bonded ratio > 1",
			bondedRatio:    sdk.NewDec(2),
			stakingSupply:  math.NewInt(1000000),
			expectError:    true,
			errorString:    "bonded ratio cannot be greater than 1",
			expectedResult: sdk.Dec{},
		},
		{
			name:           "zero staking supply",
			bondedRatio:    sdk.NewDecWithPrec(5, 1),
			stakingSupply:  math.ZeroInt(),
			expectError:    false,
			expectedResult: sdk.ZeroDec(),
		},
		{
			name:           "zero bonded ratio",
			bondedRatio:    sdk.ZeroDec(),
			stakingSupply:  math.NewInt(1000000),
			expectError:    false,
			expectedResult: sdk.ZeroDec(),
		},
		{
			name:           "valid case",
			bondedRatio:    sdk.NewDecWithPrec(5, 1),
			stakingSupply:  math.NewInt(1000000),
			expectError:    false,
			expectedResult: sdk.NewDecWithPrec(7, 2), // 0.07 inflation rate
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := minter.NextInflationRate(params, tc.bondedRatio, tc.stakingSupply)
			if tc.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorString)
			} else {
				require.NoError(t, err)
				require.True(t, tc.expectedResult.Equal(result),
					"expected %s, got %s", tc.expectedResult, result)
			}
		})
	}
}

func TestBlockProvisionEdgeCases(t *testing.T) {
	minter := DefaultInitialMinter()
	params := DefaultParams()

	tests := []struct {
		name          string
		mintDenom     string
		blocksPerYear uint64
		expectError   bool
		errorString   string
	}{
		{
			name:          "empty mint denom",
			mintDenom:     "",
			blocksPerYear: params.BlocksPerYear,
			expectError:   true,
			errorString:   "mint denom cannot be empty",
		},
		{
			name:          "invalid mint denom",
			mintDenom:     "invalid!denom",
			blocksPerYear: params.BlocksPerYear,
			expectError:   true,
			errorString:   "invalid mint denom",
		},
		{
			name:          "zero blocks per year",
			mintDenom:     params.MintDenom,
			blocksPerYear: 0,
			expectError:   true,
			errorString:   "blocks per year must be positive",
		},
		{
			name:          "valid case",
			mintDenom:     params.MintDenom,
			blocksPerYear: params.BlocksPerYear,
			expectError:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testParams := params
			testParams.MintDenom = tc.mintDenom
			testParams.BlocksPerYear = tc.blocksPerYear

			_, err := minter.BlockProvision(testParams)
			if tc.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorString)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNextAnnualProvisionsEdgeCases(t *testing.T) {
	minter := DefaultInitialMinter()
	params := DefaultParams()

	tests := []struct {
		name        string
		totalSupply math.Int
		expectError bool
		errorString string
	}{
		{
			name:        "zero total supply",
			totalSupply: sdk.ZeroInt(),
			expectError: false,
		},
		{
			name:        "negative total supply",
			totalSupply: sdk.NewInt(-1),
			expectError: true,
			errorString: "total supply cannot be negative",
		},
		{
			name:        "maximum total supply",
			totalSupply: sdk.NewIntFromUint64(^uint64(0)),
			expectError: true,
			errorString: "total supply too large",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := minter.NextAnnualProvisions(params, tc.totalSupply)
			if tc.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorString)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}
		})
	}
}

type BlockProvisionTestCase struct {
	name          string
	minter        Minter
	params        Params
	totalSupply   math.Int
	expProvisions sdk.Dec
	expError      bool
}

type AnnualProvisionsTestCase struct {
	name          string
	minter        Minter
	params        Params
	totalSupply   math.Int
	expProvisions sdk.Dec
	expError      bool
}

func TestBlockProvision(t *testing.T) {
	tests := []BlockProvisionTestCase{
		{
			name:          "default params",
			minter:        DefaultInitialMinter(),
			params:        DefaultParams(),
			totalSupply:   math.NewInt(1000000),
			expProvisions: sdk.NewDec(1),
			expError:      false,
		},
		{
			name:   "zero total supply",
			minter: InitialMinter(sdk.NewDecWithPrec(10, 2)), // 10% inflation instead of 13%
			params: Params{
				MintDenom:           "stake",
				InflationRateChange: sdk.NewDecWithPrec(13, 2),
				GoalBonded:          sdk.NewDecWithPrec(67, 2),
				BlocksPerYear:       uint64(60 * 60 * 8766 / 6), // Different from special case
				MaxTokenPerYear:     sdk.NewIntFromUint64(1000000000000000),
				MinTokenPerYear:     sdk.NewIntFromUint64(800000000000000),
			},
			totalSupply:   math.ZeroInt(),
			expProvisions: sdk.ZeroDec(),
			expError:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.minter.BlockProvision(tc.params)
			if tc.expError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.True(t, tc.expProvisions.Equal(sdk.NewDecFromInt(result.Amount)))
			}
		})
	}
}

func TestNextInflationRate(t *testing.T) {
	tests := []struct {
		name          string
		minter        Minter
		params        Params
		bondedRatio   sdk.Dec
		stakingSupply math.Int
		expRate       sdk.Dec
	}{
		{
			name:          "inflation at goal",
			minter:        DefaultInitialMinter(),
			params:        DefaultParams(),
			bondedRatio:   sdk.NewDecWithPrec(67, 2),
			stakingSupply: math.NewInt(1000000),
			expRate:       sdk.NewDecWithPrec(13, 2),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.minter.NextInflationRate(tc.params, tc.bondedRatio, tc.stakingSupply)
			require.NoError(t, err)
			require.True(t, tc.expRate.Equal(result))
		})
	}
}

func TestNextAnnualProvisions(t *testing.T) {
	tests := []struct {
		name          string
		minter        Minter
		params        Params
		stakingSupply math.Int
		expProvisions sdk.Dec
	}{
		{
			name:          "zero staking supply",
			minter:        Minter{},
			params:        Params{},
			stakingSupply: math.ZeroInt(),
			expProvisions: sdk.ZeroDec(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.minter.NextAnnualProvisions(tc.params, tc.stakingSupply)
			require.NoError(t, err)
			require.True(t, tc.expProvisions.Equal(result))
		})
	}
}

func TestAnnualProvisions(t *testing.T) {
	tests := []AnnualProvisionsTestCase{
		{
			name:          "default params",
			minter:        DefaultInitialMinter(),
			params:        DefaultParams(),
			totalSupply:   math.NewInt(1000000),
			expProvisions: sdk.NewDec(1),
			expError:      false,
		},
		{
			name:          "zero total supply",
			minter:        DefaultInitialMinter(),
			params:        DefaultParams(),
			totalSupply:   math.ZeroInt(),
			expProvisions: sdk.ZeroDec(),
			expError:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.minter.NextAnnualProvisions(tc.params, tc.totalSupply)
			if tc.expError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.True(t, tc.expProvisions.Equal(result))
			}
		})
	}
}
