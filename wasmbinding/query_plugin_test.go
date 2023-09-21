package wasmbinding_test

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/golang/protobuf/proto" //nolint:staticcheck
	"github.com/stretchr/testify/suite"

	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	"github.com/notional-labs/centauri/v5/app"
	helpers "github.com/notional-labs/centauri/v5/app/helpers"
	"github.com/notional-labs/centauri/v5/wasmbinding"
)

type StargateTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *app.CentauriApp
}

func (s *StargateTestSuite) SetupTest() {
	s.app = helpers.SetupCentauriAppWithValSet(s.T())
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "centauri-1", Time: time.Now().UTC()})
	s.app.IBCKeeper.ChannelKeeper.SetNextSequenceSend(s.ctx, "transfer", "channel-1", 102)
}

func TestStargateTestSuite(t *testing.T) {
	suite.Run(t, new(StargateTestSuite))
}

func (s *StargateTestSuite) TestStargateQuerier() {
	testCases := []struct {
		name                   string
		testSetup              func()
		path                   string
		requestData            func() []byte
		responseProtoStruct    interface{}
		expectedQuerierError   bool
		expectedUnMarshalError bool
		resendRequest          bool
	}{
		{
			name: "happy path",
			path: "/ibc.core.channel.v1.Query/NextSequenceSend",
			requestData: func() []byte {
				epochrequest := channeltypes.QueryNextSequenceSendRequest{
					PortId: "transfer",
					ChannelId: "channel-1",
				}
				bz, err := proto.Marshal(&epochrequest)
				s.Require().NoError(err)
				return bz
			},
			responseProtoStruct: &channeltypes.QueryNextSequenceSendResponse{},
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.SetupTest()
			if tc.testSetup != nil {
				tc.testSetup()
			}


			stargateQuerier := wasmbinding.StargateQuerier(*s.app.GRPCQueryRouter(), s.app.AppCodec())
			stargateRequest := &wasmvmtypes.StargateQuery{
				Path: tc.path,
				Data: tc.requestData(),
			}
			stargateResponse, err := stargateQuerier(s.ctx, stargateRequest)
			if tc.expectedQuerierError {
				s.Require().Error(err)
				return
			}

			s.Require().NoError(err)

			protoResponse, ok := tc.responseProtoStruct.(proto.Message)
			s.Require().True(ok)

			// test correctness by unmarshalling json response into proto struct
			err = s.app.AppCodec().UnmarshalJSON(stargateResponse, protoResponse)
			if tc.expectedUnMarshalError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(protoResponse)
			}

			if tc.resendRequest {
				stargateQuerier = wasmbinding.StargateQuerier(*s.app.GRPCQueryRouter(), s.app.AppCodec())
				stargateRequest = &wasmvmtypes.StargateQuery{
					Path: tc.path,
					Data: tc.requestData(),
				}
				resendResponse, err := stargateQuerier(s.ctx, stargateRequest)
				s.Require().NoError(err)
				s.Require().Equal(stargateResponse, resendResponse)
			}
		})
	}
}

func (s *StargateTestSuite) TestConvertProtoToJsonMarshal() {
	testCases := []struct {
		name                  string
		queryPath             string
		protoResponseStruct   codec.ProtoMarshaler
		originalResponse      string
		expectedProtoResponse codec.ProtoMarshaler
		expectedError         bool
	}{
		{
			name:                "successful conversion from proto response to json marshalled response",
			queryPath:           "/cosmos.bank.v1beta1.Query/AllBalances",
			originalResponse:    "0a090a036261721202333012050a03666f6f",
			protoResponseStruct: &banktypes.QueryAllBalancesResponse{},
			expectedProtoResponse: &banktypes.QueryAllBalancesResponse{
				Balances: sdk.NewCoins(sdk.NewCoin("bar", sdk.NewInt(30))),
				Pagination: &query.PageResponse{
					NextKey: []byte("foo"),
				},
			},
		},
		{
			name:                "invalid proto response struct",
			queryPath:           "/cosmos.bank.v1beta1.Query/AllBalances",
			originalResponse:    "0a090a036261721202333012050a03666f6f",
			protoResponseStruct: &channeltypes.QueryNextSequenceSendResponse{},
			expectedError:       true,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.SetupTest()

			originalVersionBz, err := hex.DecodeString(tc.originalResponse)
			s.Require().NoError(err)

			jsonMarshalledResponse, err := wasmbinding.ConvertProtoToJSONMarshal(tc.protoResponseStruct, originalVersionBz, s.app.AppCodec())
			if tc.expectedError {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			// check response by json marshalling proto response into json response manually
			jsonMarshalExpectedResponse, err := s.app.AppCodec().MarshalJSON(tc.expectedProtoResponse)
			s.Require().NoError(err)
			s.Require().Equal(jsonMarshalledResponse, jsonMarshalExpectedResponse)
		})
	}
}
