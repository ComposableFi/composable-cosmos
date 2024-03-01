package v6_4_6_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	apptesting "github.com/notional-labs/composable/v6/app"
	"github.com/notional-labs/composable/v6/app/upgrades/v6_4_6"
	"github.com/notional-labs/composable/v6/bech32-migration/utils"
	"github.com/stretchr/testify/suite"
	alliancetypes "github.com/terra-money/alliance/x/alliance/types"
)

const (
	COIN_DENOM   = "upica"
	CONNECTION_0 = "connection-0"
	PORT_0       = "port-0"
)

type UpgradeTestSuite struct {
	apptesting.KeeperTestHelper
}

func TestUpgradeTestSuite(t *testing.T) {
	suite.Run(t, new(UpgradeTestSuite))
}

// Ensures the test does not error out.
func (s *UpgradeTestSuite) TestForMigratingNewPrefix() {
	// DEFAULT PREFIX: centauri
	sdk.SetAddrCacheEnabled(false)

	sdk.GetConfig().SetBech32PrefixForAccount(utils.OldBech32PrefixAccAddr, utils.OldBech32PrefixAccPub)
	sdk.GetConfig().SetBech32PrefixForValidator(utils.OldBech32PrefixValAddr, utils.OldBech32PrefixValPub)
	sdk.GetConfig().SetBech32PrefixForConsensusNode(utils.OldBech32PrefixConsAddr, utils.OldBech32PrefixConsPub)

	s.Setup(s.T())

	acc1, proposal := prepareForTestingGovModule(s)

	oldConsAddress := prepareForTestingSlashingModule(s)

	oldValAddress, oldValAddress2, acc3, afterOneDay := prepareForTestingStakingModule(s)

	baseAccount, stakingModuleAccount, baseVestingAccount, continuousVestingAccount, delayedVestingAccount, periodicVestingAccount, permanentLockedAccount := prepareForTestingAuthModule(s)

	prepareForTestingAllianceModule(s)

	prepareForTestingICAHostModule(s)

	/* == UPGRADE == */
	upgradeHeight := int64(5)
	s.ConfirmUpgradeSucceeded(v6_4_6.UpgradeName, upgradeHeight)

	/* == CHECK AFTER UPGRADE == */
	checkUpgradeGovModule(s, acc1, proposal)
	checkUpgradeSlashingModule(s, oldConsAddress)
	checkUpgradeStakingModule(s, oldValAddress, oldValAddress2, acc3, afterOneDay)
	checkUpgradeAuthModule(s, baseAccount, stakingModuleAccount, baseVestingAccount, continuousVestingAccount, delayedVestingAccount, periodicVestingAccount, permanentLockedAccount)
	checkUpgradeAllianceModule(s)
	checkUpgradeICAHostModule(s)
}

func prepareForTestingGovModule(s *UpgradeTestSuite) (sdk.AccAddress, govtypes.Proposal) {
	/* PREPARE FOR TESTING GOV MODULE */
	acc1 := s.TestAccs[0]

	// MINT NEW TOKEN FOR BALANCE CHECKING
	s.App.BankKeeper.MintCoins(s.Ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(COIN_DENOM, math.NewIntFromUint64(100000000))))

	// VOTE AND DEPOSIT
	proposal, err := s.App.GovKeeper.SubmitProposal(s.Ctx, []sdk.Msg{}, "", "test", "description", acc1)
	s.Suite.Equal(err, nil)

	s.App.GovKeeper.SetVote(s.Ctx, govtypes.Vote{
		ProposalId: proposal.Id,
		Voter:      acc1.String(),
		Options:    nil,
		Metadata:   "",
	})

	s.App.GovKeeper.SetDeposit(s.Ctx, govtypes.Deposit{
		ProposalId: proposal.Id,
		Depositor:  acc1.String(),
		Amount:     sdk.NewCoins(sdk.NewCoin(COIN_DENOM, math.NewIntFromUint64(1))),
	})

	return acc1, proposal
}

func prepareForTestingSlashingModule(s *UpgradeTestSuite) sdk.ConsAddress {
	/* PREPARE FOR TESTING SLASHING MODULE */
	acc2 := s.TestAccs[1]

	oldConsAddress, err := utils.ConsAddressFromOldBech32(acc2.String(), utils.OldBech32PrefixAccAddr)
	s.Suite.Equal(err, nil)

	// CHECK ValidatorSigningInfo
	s.App.SlashingKeeper.SetValidatorSigningInfo(s.Ctx, oldConsAddress, slashingtypes.ValidatorSigningInfo{
		Address: oldConsAddress.String(),
	})
	return oldConsAddress
}

func prepareForTestingStakingModule(s *UpgradeTestSuite) (sdk.ValAddress, sdk.ValAddress, sdk.AccAddress, time.Time) {
	/* PREPARE FOR TESTING SLASHING MODULE */
	acc3 := s.TestAccs[2]

	// MINT NEW TOKEN FOR BALANCE CHECKING
	s.App.BankKeeper.MintCoins(s.Ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(COIN_DENOM, math.NewIntFromUint64(100000000))))
	s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, minttypes.ModuleName, acc3, sdk.NewCoins(sdk.NewCoin(COIN_DENOM, math.NewIntFromUint64(100000000))))

	// validator.OperatorAddress
	oldValAddress := s.SetupValidator(stakingtypes.Bonded)
	oldValAddress2 := s.SetupValidator(stakingtypes.Bonded)

	// delegation.DelegatorAddress & delegation.ValidatorAddress
	s.StakingHelper.Delegate(acc3, oldValAddress, sdk.NewInt(300))

	// redelegation.DelegatorAddress & redelegation.ValidatorSrcAddress & redelegation.ValidatorDstAddress
	completionTime, err := s.App.StakingKeeper.BeginRedelegation(s.Ctx, acc3, oldValAddress, oldValAddress2, sdk.NewDec(100))
	afterOneDay := completionTime.AddDate(0, 0, 1)
	s.Require().NoError(err)

	// Undelegate part of the tokens from val2 (test instant unbonding on undelegation started before upgrade)
	s.StakingHelper.Undelegate(acc3, oldValAddress, sdk.NewInt(10), true)
	return oldValAddress, oldValAddress2, acc3, afterOneDay
}

func prepareForTestingAuthModule(s *UpgradeTestSuite) (sdk.AccAddress, sdk.AccAddress, sdk.AccAddress, sdk.AccAddress, sdk.AccAddress, sdk.AccAddress, sdk.AccAddress) {
	addr0 := s.TestAccs[0]
	baseAccount := authtypes.NewBaseAccount(addr0, nil, 0, 0)
	s.App.AccountKeeper.SetAccount(s.Ctx, baseAccount)

	addr6 := s.TestAccs[6]
	baseAccount6 := authtypes.NewBaseAccount(addr6, nil, 0, 0)
	stakingModuleAccount := authtypes.NewModuleAccount(baseAccount6, "name", "name")
	s.App.AccountKeeper.SetAccount(s.Ctx, stakingModuleAccount)

	addr2 := s.TestAccs[2]
	baseAccount2 := authtypes.NewBaseAccount(addr2, nil, 0, 0)
	baseVestingAccount := authvesting.NewBaseVestingAccount(baseAccount2, sdk.NewCoins(sdk.NewCoin(COIN_DENOM, math.NewIntFromUint64(1))), 60)
	s.App.AccountKeeper.SetAccount(s.Ctx, baseVestingAccount)

	continuousVestingAccount := CreateVestingAccount(s)

	addr3 := s.TestAccs[3]
	baseAccount3 := authtypes.NewBaseAccount(addr3, nil, 0, 0)
	baseVestingAccount2 := authvesting.NewBaseVestingAccount(baseAccount3, sdk.NewCoins(sdk.NewCoin(COIN_DENOM, math.NewIntFromUint64(1))), 60)
	delayedVestingAccount := authvesting.NewDelayedVestingAccountRaw(baseVestingAccount2)
	s.App.AccountKeeper.SetAccount(s.Ctx, delayedVestingAccount)

	addr4 := s.TestAccs[4]
	baseAccount4 := authtypes.NewBaseAccount(addr4, nil, 0, 0)
	baseVestingAccount3 := authvesting.NewBaseVestingAccount(baseAccount4, sdk.NewCoins(sdk.NewCoin(COIN_DENOM, math.NewIntFromUint64(1))), 60)
	periodicVestingAccount := authvesting.NewPeriodicVestingAccountRaw(baseVestingAccount3, 0, vestingtypes.Periods{})
	s.App.AccountKeeper.SetAccount(s.Ctx, periodicVestingAccount)

	addr5 := s.TestAccs[5]
	baseAccount5 := authtypes.NewBaseAccount(addr5, nil, 0, 0)
	permanentLockedAccount := authvesting.NewPermanentLockedAccount(baseAccount5, sdk.NewCoins(sdk.NewCoin(COIN_DENOM, math.NewIntFromUint64(1))))
	s.App.AccountKeeper.SetAccount(s.Ctx, permanentLockedAccount)

	return baseAccount.GetAddress(), stakingModuleAccount.GetAddress(), baseVestingAccount.GetAddress(), continuousVestingAccount.GetAddress(), delayedVestingAccount.GetAddress(), periodicVestingAccount.GetAddress(), permanentLockedAccount.GetAddress()
}

func prepareForTestingICAHostModule(s *UpgradeTestSuite) {
	acc1 := s.TestAccs[0]
	s.App.ICAHostKeeper.SetInterchainAccountAddress(s.Ctx, CONNECTION_0, PORT_0, acc1.String())
}

func prepareForTestingAllianceModule(s *UpgradeTestSuite) {
	oldValAddress := s.SetupValidator(stakingtypes.Bonded)
	_, bz, _ := bech32.DecodeAndConvert(oldValAddress.String())
	oldBech32Addr, _ := bech32.ConvertAndEncode(utils.OldBech32PrefixValAddr, bz)

	s.App.AllianceKeeper.InitGenesis(s.Ctx, &alliancetypes.GenesisState{
		ValidatorInfos: []alliancetypes.ValidatorInfoState{{
			ValidatorAddress: oldBech32Addr,
			Validator:        alliancetypes.NewAllianceValidatorInfo(),
		}},
	})
}

func checkUpgradeGovModule(s *UpgradeTestSuite, acc1 sdk.AccAddress, proposal govtypes.Proposal) {
	// CONVERT ACC TO NEW PREFIX
	_, bz, _ := bech32.DecodeAndConvert(acc1.String())
	newBech32Addr, _ := bech32.ConvertAndEncode(utils.NewBech32PrefixAccAddr, bz)
	newAddr, err := utils.AccAddressFromOldBech32(newBech32Addr, utils.NewBech32PrefixAccAddr)
	s.Suite.Equal(err, nil)

	// CHECK PROPOSAL
	proposal, found := s.App.GovKeeper.GetProposal(s.Ctx, proposal.Id)
	s.Suite.Equal(found, true)
	s.Suite.Equal(proposal.Proposer, newBech32Addr)

	// CHECK VOTER AND DEPOSITER OF NEW ADDRESS
	existed_proposal, _ := s.App.GovKeeper.GetProposal(s.Ctx, proposal.Id)
	s.Suite.Equal(existed_proposal.Proposer, newBech32Addr)

	vote, found := s.App.GovKeeper.GetVote(s.Ctx, proposal.Id, newAddr)
	s.Suite.Equal(found, true)
	s.Suite.Equal(vote.Voter, newBech32Addr)

	deposit, found := s.App.GovKeeper.GetDeposit(s.Ctx, proposal.Id, newAddr)
	s.Suite.Equal(found, true)
	s.Suite.Equal(deposit.Depositor, newBech32Addr)
}

func checkUpgradeSlashingModule(s *UpgradeTestSuite, oldConsAddress sdk.ConsAddress) {
	// CONVERT TO ACC TO NEW PREFIX
	_, bz, _ := bech32.DecodeAndConvert(oldConsAddress.String())
	newBech32Addr, _ := bech32.ConvertAndEncode(utils.NewBech32PrefixConsAddr, bz)
	newAddr, err := utils.ConsAddressFromOldBech32(newBech32Addr, utils.NewBech32PrefixConsAddr)
	s.Suite.Equal(err, nil)

	valSigningInfo, found := s.App.SlashingKeeper.GetValidatorSigningInfo(s.Ctx, newAddr)
	s.Suite.Equal(found, true)
	s.Suite.Equal(valSigningInfo.Address, newBech32Addr)
}

func checkUpgradeStakingModule(s *UpgradeTestSuite, oldValAddress sdk.ValAddress, oldValAddress2 sdk.ValAddress, acc1 sdk.AccAddress, afterOneDay time.Time) {
	// CONVERT TO ACC TO NEW PREFIX
	_, bz, _ := bech32.DecodeAndConvert(oldValAddress.String())
	newBech32Addr, _ := bech32.ConvertAndEncode(utils.NewBech32PrefixValAddr, bz)
	newValAddr, err := utils.ValAddressFromOldBech32(newBech32Addr, utils.NewBech32PrefixValAddr)
	s.Suite.Equal(err, nil)

	_, bzVal2, _ := bech32.DecodeAndConvert(oldValAddress2.String())
	newBech32AddrVal2, _ := bech32.ConvertAndEncode(utils.NewBech32PrefixValAddr, bzVal2)
	newValAddr2, err := utils.ValAddressFromOldBech32(newBech32AddrVal2, utils.NewBech32PrefixValAddr)
	s.Suite.Equal(err, nil)

	_, bz1, _ := bech32.DecodeAndConvert(acc1.String())
	newBech32DelAddr, _ := bech32.ConvertAndEncode(utils.NewBech32PrefixAccAddr, bz1)
	newAccAddr, err := utils.AccAddressFromOldBech32(newBech32DelAddr, utils.NewBech32PrefixAccAddr)
	s.Suite.Equal(err, nil)

	val, found := s.App.StakingKeeper.GetValidator(s.Ctx, newValAddr)
	s.Suite.Equal(found, true)
	s.Suite.Equal(val.OperatorAddress, newBech32Addr)

	delegation, found := s.App.StakingKeeper.GetDelegation(s.Ctx, newAccAddr, newValAddr)
	s.Suite.Equal(found, true)
	s.Suite.Equal(delegation.DelegatorAddress, newBech32DelAddr)
	s.Suite.Equal(delegation.ValidatorAddress, newBech32Addr)

	unbonding, found := s.App.StakingKeeper.GetUnbondingDelegation(s.Ctx, newAccAddr, newValAddr)
	s.Suite.Equal(found, true)
	s.Suite.Equal(unbonding.DelegatorAddress, newBech32DelAddr)
	s.Suite.Equal(unbonding.ValidatorAddress, newBech32Addr)

	s.Ctx = s.Ctx.WithBlockTime(afterOneDay)

	redelegation, found := s.App.StakingKeeper.GetRedelegation(s.Ctx, newAccAddr, newValAddr, newValAddr2)
	s.Suite.Equal(found, true)
	s.Suite.Equal(redelegation.DelegatorAddress, newBech32DelAddr)
	s.Suite.Equal(redelegation.ValidatorSrcAddress, newBech32Addr)
	s.Suite.Equal(redelegation.ValidatorDstAddress, newBech32AddrVal2)
}

func checkUpgradeAuthModule(s *UpgradeTestSuite, baseAccount sdk.AccAddress, stakingModuleAccount sdk.AccAddress, baseVestingAccount sdk.AccAddress, continuousVestingAccount sdk.AccAddress, delayedVestingAccount sdk.AccAddress, periodicVestingAccount sdk.AccAddress, permanentLockedAccount sdk.AccAddress) {
	/* CHECK BASE ACCOUNT */
	_, bz, _ := bech32.DecodeAndConvert(baseAccount.String())
	newBech32AddrBaseAccount, _ := bech32.ConvertAndEncode(utils.NewBech32PrefixAccAddr, bz)
	var newPrefixAddrBA authtypes.AccountI
	newPrefixAddrBA = s.App.AccountKeeper.GetAccount(s.Ctx, baseAccount)
	switch acci := newPrefixAddrBA.(type) {
	case *authtypes.BaseAccount:
		acc := acci
		s.Suite.Equal(acc.Address, newBech32AddrBaseAccount)
	default:
		s.Suite.NotNil(nil)
	}

	/* CHECK MODULE ACCOUNT */
	_, bz, _ = bech32.DecodeAndConvert(stakingModuleAccount.String())
	newBech32AddrModuleAccount, _ := bech32.ConvertAndEncode(utils.NewBech32PrefixAccAddr, bz)
	var newPrefixAddrMA authtypes.AccountI
	newPrefixAddrMA = s.App.AccountKeeper.GetAccount(s.Ctx, stakingModuleAccount)
	switch acci := newPrefixAddrMA.(type) {
	case *authtypes.ModuleAccount:
		acc := acci
		s.Suite.Equal(acc.Address, newBech32AddrModuleAccount)
	default:
		s.Suite.NotNil(nil)
	}

	/* CHECK BASE VESTING ACCOUNT */
	_, bz, _ = bech32.DecodeAndConvert(baseVestingAccount.String())
	newBech32AddrBaseVestingAccount, _ := bech32.ConvertAndEncode(utils.NewBech32PrefixAccAddr, bz)
	var newPrefixAddrBVA authtypes.AccountI
	newPrefixAddrBVA = s.App.AccountKeeper.GetAccount(s.Ctx, baseVestingAccount)
	switch acci := newPrefixAddrBVA.(type) {
	case *vestingtypes.BaseVestingAccount:
		acc := acci
		s.Suite.Equal(acc.Address, newBech32AddrBaseVestingAccount)
	default:
		s.Suite.NotNil(nil)
	}

	// CHECK CONTINUOUS VESTING ACCOUNT AND MULTISIG
	_, bz, _ = bech32.DecodeAndConvert(continuousVestingAccount.String())
	newBech32AddrConVestingAccount, _ := bech32.ConvertAndEncode(utils.NewBech32PrefixAccAddr, bz)
	var newPrefixAddrCVA authtypes.AccountI
	newPrefixAddrCVA = s.App.AccountKeeper.GetAccount(s.Ctx, continuousVestingAccount)
	switch acci := newPrefixAddrCVA.(type) {
	case *vestingtypes.ContinuousVestingAccount:
		acc := acci
		s.Suite.Equal(acc.Address, newBech32AddrConVestingAccount)
	default:
		s.Suite.NotNil(nil)
	}

	// CHECK DELAYED VESTING ACCOUNT
	_, bz, _ = bech32.DecodeAndConvert(delayedVestingAccount.String())
	newBech32AddrDelayedVestingAccount, _ := bech32.ConvertAndEncode(utils.NewBech32PrefixAccAddr, bz)
	var newPrefixAddrDVA authtypes.AccountI
	newPrefixAddrDVA = s.App.AccountKeeper.GetAccount(s.Ctx, delayedVestingAccount)
	switch acci := newPrefixAddrDVA.(type) {
	case *vestingtypes.DelayedVestingAccount:
		acc := acci
		s.Suite.Equal(acc.Address, newBech32AddrDelayedVestingAccount)
	default:
		s.Suite.NotNil(nil)
	}

	// CHECK PERIODIC VESTING ACCOUNT
	_, bz, _ = bech32.DecodeAndConvert(periodicVestingAccount.String())
	newBech32AddrPeriodicVestingAccount, _ := bech32.ConvertAndEncode(utils.NewBech32PrefixAccAddr, bz)
	var newPrefixAddrPVA authtypes.AccountI
	newPrefixAddrPVA = s.App.AccountKeeper.GetAccount(s.Ctx, periodicVestingAccount)
	switch acci := newPrefixAddrPVA.(type) {
	case *vestingtypes.PeriodicVestingAccount:
		acc := acci
		s.Suite.Equal(acc.Address, newBech32AddrPeriodicVestingAccount)
	default:
		s.Suite.NotNil(nil)
	}

	// CHECK PERMANENT LOCKED ACCOUNT
	_, bz, _ = bech32.DecodeAndConvert(permanentLockedAccount.String())
	newBech32AddrPermanentVestingAccount, _ := bech32.ConvertAndEncode(utils.NewBech32PrefixAccAddr, bz)
	var newPrefixAddrPLA authtypes.AccountI
	newPrefixAddrPLA = s.App.AccountKeeper.GetAccount(s.Ctx, permanentLockedAccount)
	switch acci := newPrefixAddrPLA.(type) {
	case *vestingtypes.PermanentLockedAccount:
		acc := acci
		s.Suite.Equal(acc.Address, newBech32AddrPermanentVestingAccount)
	default:
		s.Suite.NotNil(nil)
	}
}

func checkUpgradeAllianceModule(s *UpgradeTestSuite) {
	// the validator address in alliance genesis file is converted into accAddr type
	// and then used for key storage
	// so the migration do not affect this module
	genesis := s.App.AllianceKeeper.ExportGenesis(s.Ctx)
	s.Suite.Equal(strings.Contains(genesis.ValidatorInfos[0].ValidatorAddress, "pica"), true)
}

func checkUpgradeICAHostModule(s *UpgradeTestSuite) {
	acc1 := s.TestAccs[0]
	interchainAccount, _ := s.App.ICAHostKeeper.GetInterchainAccountAddress(s.Ctx, CONNECTION_0, PORT_0)
	s.Suite.Equal(acc1.String(), interchainAccount)
}

func CreateVestingAccount(s *UpgradeTestSuite,
) vestingtypes.ContinuousVestingAccount {
	str := `{"@type":"/cosmos.vesting.v1beta1.ContinuousVestingAccount","base_vesting_account":{"base_account":{"address":"centauri1alga5e8vr6ccr9yrg0kgxevpt5xgmgrvfkc5p8","pub_key":{"@type":"/cosmos.crypto.multisig.LegacyAminoPubKey","threshold":4,"public_keys":[{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"AlnzK22KrkylnvTCvZZc8eZnydtQuzCWLjJJSMFUvVHf"},{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"Aiw2Ftg+fnoHDU7M3b0VMRsI0qurXlerW0ahtfzSDZA4"},{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"AvEHv+MVYRVau8FbBcJyG0ql85Tbbn7yhSA0VGmAY4ku"},{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"Az5VHWqi3zMJu1rLGcu2EgNXLLN+al4Dy/lj6UZTzTCl"},{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"Ai4GlSH3uG+joMnAFbQC3jQeHl9FPvVTlRmwIFt7d7TI"},{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"A2kAzH2bZr530jmFq/bRFrT2q8SRqdnfIebba+YIBqI1"}]},"account_number":46,"sequence":27},"original_vesting":[{"denom":"upica","amount":"22165200000000"}],"delegated_free":[{"denom":"upica","amount":"443382497453"}],"delegated_vesting":[{"denom":"upica","amount":"22129422502547"}],"end_time":1770994800},"start_time":1676300400}`

	var acc vestingtypes.ContinuousVestingAccount
	if err := json.Unmarshal([]byte(str), &acc); err != nil {
		panic(err)
	}

	err := banktestutil.FundAccount(s.App.BankKeeper, s.Ctx, acc.BaseAccount.GetAddress(),
		acc.GetOriginalVesting())
	if err != nil {
		panic(err)
	}

	err = banktestutil.FundAccount(s.App.BankKeeper, s.Ctx, acc.BaseAccount.GetAddress(),
		sdk.NewCoins(sdk.NewCoin(COIN_DENOM, math.NewIntFromUint64(1))))
	if err != nil {
		panic(err)
	}

	s.App.AccountKeeper.SetAccount(s.Ctx, &acc)
	return acc
}
