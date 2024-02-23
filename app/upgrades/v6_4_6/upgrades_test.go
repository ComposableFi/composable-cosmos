package v6_4_6_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	apptesting "github.com/notional-labs/composable/v6/app"
	"github.com/notional-labs/composable/v6/app/upgrades/v6_4_6"
	"github.com/notional-labs/composable/v6/bech32-migration/utils"
	"github.com/stretchr/testify/suite"
)

const (
	COIN_DENOM = "upica"
)

type UpgradeTestSuite struct {
	apptesting.KeeperTestHelper
}

func TestUpgradeTestSuite(t *testing.T) {
	suite.Run(t, new(UpgradeTestSuite))
}

// Ensures the test does not error out.
func (s *UpgradeTestSuite) TestForMigratingNewPrefix() {
	s.Setup(s.T())

	acc1, proposal, oldAccBal := prepareForTestingGovModule(s)

	oldConsAddress := prepareForTestingSlashingModule(s)

	oldValAddress, oldValAddress2, acc3, afterOneDay := prepareForTestingStakingModule(s)

	/* == UPGRADE == */
	upgradeHeight := int64(5)
	s.ConfirmUpgradeSucceeded(v6_4_6.UpgradeName, upgradeHeight)

	/* == CHECK AFTER UPGRADE == */
	checkUpgradeGovModule(s, acc1, proposal, oldAccBal)
	checkUpgradeSlashingModule(s, oldConsAddress)
	checkUpgradeStakingModule(s, oldValAddress, oldValAddress2, acc3, afterOneDay)
	// migration of auth has been well tested on above cases
}

func prepareForTestingGovModule(s *UpgradeTestSuite) (sdk.AccAddress, govtypes.Proposal, sdk.Coin) {
	/* PREPARE FOR TESTING GOV MODULE */
	acc1 := s.TestAccs[0]

	// MINT NEW TOKEN FOR BALANCE CHECKING
	s.App.BankKeeper.MintCoins(s.Ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(COIN_DENOM, math.NewIntFromUint64(100000000))))
	s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, minttypes.ModuleName, acc1, sdk.NewCoins(sdk.NewCoin(COIN_DENOM, math.NewIntFromUint64(1))))

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

	// CHECK BALANCE OF A ACCOUNT
	oldAccBal := s.App.BankKeeper.GetBalance(s.Ctx, acc1, COIN_DENOM)
	return acc1, proposal, oldAccBal
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

func checkUpgradeGovModule(s *UpgradeTestSuite, acc1 sdk.AccAddress, proposal govtypes.Proposal, oldAccBal sdk.Coin) {
	// CONVERT ACC TO NEW PREFIX
	_, bz, _ := bech32.DecodeAndConvert(acc1.String())
	newBech32Addr, _ := bech32.ConvertAndEncode(utils.NewBech32PrefixAccAddr, bz)
	newAddr, err := utils.AccAddressFromOldBech32(newBech32Addr, utils.NewBech32PrefixAccAddr)
	s.Suite.Equal(err, nil)

	// CHECK PROPOSAL
	proposal, found := s.App.GovKeeper.GetProposal(s.Ctx, proposal.Id)
	s.Suite.Equal(found, true)
	s.Suite.Equal(proposal.Proposer, newBech32Addr)

	// CHECK BALANCE OF NEW ADDRESS
	newAccBal := s.App.BankKeeper.GetBalance(s.Ctx, newAddr, COIN_DENOM)
	s.Suite.Equal(oldAccBal, newAccBal)

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
	// parsedNewPrefixValAddr, _ := utils.ValAddressFromOldBech32(newBech32Addr, utils.NewBech32PrefixValAddr)
	newValAddr, err := utils.ValAddressFromOldBech32(newBech32Addr, utils.NewBech32PrefixValAddr)
	s.Suite.Equal(err, nil)

	_, bzVal2, _ := bech32.DecodeAndConvert(oldValAddress2.String())
	newBech32AddrVal2, _ := bech32.ConvertAndEncode(utils.NewBech32PrefixValAddr, bzVal2)
	// parsedNewPrefixVal2Addr, _ := utils.ValAddressFromOldBech32(newBech32AddrVal2, utils.NewBech32PrefixValAddr)
	newValAddr2, err := utils.ValAddressFromOldBech32(newBech32AddrVal2, utils.NewBech32PrefixValAddr)
	s.Suite.Equal(err, nil)

	_, bz1, _ := bech32.DecodeAndConvert(acc1.String())
	newBech32DelAddr, _ := bech32.ConvertAndEncode(utils.NewBech32PrefixAccAddr, bz1)
	// parsedNewPrefixAccAddr, _ := utils.AccAddressFromOldBech32(newBech32DelAddr, utils.NewBech32PrefixAccAddr)
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
