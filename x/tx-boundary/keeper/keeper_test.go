package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/notional-labs/composable/v6/app"
	"github.com/notional-labs/composable/v6/app/helpers"
	"github.com/notional-labs/composable/v6/x/tx-boundary/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context
	// querier sdk.Querier
	app *app.ComposableApp
}

func (s *KeeperTestSuite) SetupTest() {
	s.app = helpers.SetupComposableAppWithValSet(s.T())
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "centauri-1", Time: time.Now().UTC()})
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

var (
	defaultBoundary = types.Boundary{
		TxLimit:             5,
		BlocksPerGeneration: 5,
	}
	newBoundary = types.Boundary{
		TxLimit:             10,
		BlocksPerGeneration: 5,
	}
	failBoundary = types.Boundary{
		TxLimit:             10,
		BlocksPerGeneration: 0,
	}
)

func (s *KeeperTestSuite) TestSetDelegateBoundary() {
	for _, tc := range []struct {
		desc             string
		expectedBoundary types.Boundary
		malleate         func() error
		shouldErr        bool
		expectedErr      string
	}{
		{
			desc:             "Case success",
			expectedBoundary: newBoundary,
			malleate: func() error {
				return s.app.TxBoundaryKeepper.SetDelegateBoundary(s.ctx, newBoundary)
			},
			shouldErr: false,
		},
		{
			desc:             "Case fail",
			expectedBoundary: failBoundary,
			malleate: func() error {
				return s.app.TxBoundaryKeepper.SetDelegateBoundary(s.ctx, failBoundary)
			},
			shouldErr:   true,
			expectedErr: "BlocksPerGeneration must not be zero",
		},
		{
			desc:             "Do no thing",
			expectedBoundary: defaultBoundary,
			malleate: func() error {
				return nil
			},
			shouldErr: false,
		},
	} {
		tc := tc
		s.Run(tc.desc, func() {
			s.SetupTest()
			err := tc.malleate()
			if !tc.shouldErr {
				res := s.app.TxBoundaryKeepper.GetDelegateBoundary(s.ctx)
				s.Equal(res, tc.expectedBoundary)
			} else {
				s.Equal(err.Error(), tc.expectedErr)
			}
		})
	}
}

func (s *KeeperTestSuite) TestSetRedelegateBoundary() {
	s.app.TxBoundaryKeepper.SetRedelegateBoundary(s.ctx, types.Boundary{ //nolint:errcheck
		TxLimit:             10,
		BlocksPerGeneration: 5,
	})

	for _, tc := range []struct {
		desc             string
		expectedBoundary types.Boundary
		malleate         func() error
		shouldErr        bool
		expectedErr      string
	}{
		{
			desc:             "Case success",
			expectedBoundary: newBoundary,
			malleate: func() error {
				return s.app.TxBoundaryKeepper.SetRedelegateBoundary(s.ctx, newBoundary)
			},
			shouldErr: false,
		},
		{
			desc:             "Success",
			expectedBoundary: failBoundary,
			malleate: func() error {
				return s.app.TxBoundaryKeepper.SetRedelegateBoundary(s.ctx, failBoundary)
			},
			shouldErr:   true,
			expectedErr: "BlocksPerGeneration must not be zero",
		},
		{
			desc:             "Do no thing",
			expectedBoundary: defaultBoundary,
			malleate: func() error {
				return nil
			},
			shouldErr: false,
		},
	} {
		tc := tc
		s.Run(tc.desc, func() {
			s.SetupTest()
			err := tc.malleate()
			if !tc.shouldErr {
				res := s.app.TxBoundaryKeepper.GetRedelegateBoundary(s.ctx)
				s.Equal(res, tc.expectedBoundary)
			} else {
				s.Equal(err.Error(), tc.expectedErr)
			}
		})
	}
}
