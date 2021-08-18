// Copyright (c) 2021 akachain
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package iao

import (
	"encoding/json"
	"fmt"
	"github.com/Akachain/akc-go-sdk-v2/mock"
	"github.com/Akachain/gringotts/dto/iao"
	"github.com/Akachain/gringotts/dto/token"
	"github.com/Akachain/gringotts/entity"
	"github.com/Akachain/gringotts/glossary"
	"github.com/Akachain/gringotts/glossary/doc"
	"github.com/Akachain/gringotts/helper"
	"github.com/Akachain/gringotts/smartcontract"
	"github.com/Akachain/gringotts/smartcontract/basic"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type IaoScTest struct {
	smartcontract.BasicToken
	smartcontract.Iao
}

func NewIaoSCTest() *IaoScTest {
	return &IaoScTest{basic.NewBaseToken(), NewIaoSc()}
}

type ResultAsset struct {
	AssetId string `json:"assetId"`
	TokenId string `json:"tokenId"`
}

func setupMock() (*mock.MockStubExtend, error) {
	// Initialize MockStubExtend
	chaincodeName := "IaoSC"
	sc := NewIaoSCTest()
	chaincode, _ := contractapi.NewChaincode(sc)
	stub := mock.NewMockStubExtend(shimtest.NewMockStub(chaincodeName, chaincode), chaincode, ".")

	// Create a new database, Drop old database
	db, err := mock.NewCouchDBHandler(true, chaincodeName)
	if err != nil {
		return nil, err
	}
	stub.SetCouchDBConfiguration(db)

	// Process indexes
	err = db.ProcessIndexesForChaincodeDeploy("../../META-INF/statedb/couchdb/indexes/indexPendingTx.json")
	if err != nil {
		return nil, err

	}
	return stub, nil
}

type IaoSCTestSuite struct {
	suite.Suite
	walletFromId string
	walletToId   string
	STToken      string
	ATToken      string
	AssetId      string
	IaoId        string
	stub         *mock.MockStubExtend
}

func (suite *IaoSCTestSuite) SetupTest() {
	stub, err := setupMock()
	assert.Nilf(suite.T(), err, "Setup Mock return error not nil")
	suite.stub = stub

	// create ST token type
	stableToken := token.CreateTokenType{
		Name:        "Stable Token",
		TickerToken: "ST",
		MaxSupply:   "12345678900",
	}
	paramByte, _ := json.Marshal(stableToken)
	suite.STToken = mock.MockInvokeTransaction(suite.T(), suite.stub, [][]byte{[]byte("CreateTokenType"), paramByte})
	suite.T().Log(suite.STToken)
	assert.NotEmpty(suite.T(), suite.STToken, "Create Token Type return empty")

	// create AT token type
	atToken := token.CreateTokenType{
		Name:        "Asset Token",
		TickerToken: "AT",
		MaxSupply:   "78900",
	}
	paramByte, _ = json.Marshal(atToken)
	suite.ATToken = mock.MockInvokeTransaction(suite.T(), suite.stub, [][]byte{[]byte("CreateTokenType"), paramByte})
	suite.T().Log(suite.ATToken)
	assert.NotEmpty(suite.T(), suite.ATToken, "Create Token Type return empty")

	// create wallet
	wallet := token.CreateWallet{
		TokenId: suite.STToken,
		Status:  "A",
	}
	walletByte, _ := json.Marshal(wallet)

	suite.walletFromId = mock.MockInvokeTransaction(suite.T(), suite.stub, [][]byte{[]byte("CreateWallet"), walletByte})
	assert.NotEmpty(suite.T(), suite.walletFromId, "Create from wallet return empty")

	suite.walletToId = mock.MockInvokeTransaction(suite.T(), suite.stub, [][]byte{[]byte("CreateWallet"), walletByte})
	assert.NotEmpty(suite.T(), suite.walletToId, "Create to wallet return empty")

	// mint balance for From wallet
	mintDto := new(token.MintToken)
	mintDto.WalletId = suite.walletFromId
	mintDto.TokenId = suite.STToken
	mintDto.Amount = "678900"
	mintDto.Metadata = "Req11"
	mintByte, _ := json.Marshal(mintDto)
	mintRes := mock.MockInvokeTransaction(suite.T(), suite.stub, [][]byte{[]byte("Mint"), mintByte})
	assert.Empty(suite.T(), mintRes, "Mint invoke return err")
}

func (suite *IaoSCTestSuite) TestIAO_CreateAsset() {
	suite.createNewAsset()
}

func (suite *IaoSCTestSuite) TestIAO_CreateIAO() {
	suite.createIao()
}

func (suite *IaoSCTestSuite) createIao() {
	// create new asset
	suite.createNewAsset()

	assetIaoDto := iao.AssetIao{
		AssetId:          suite.AssetId,
		AssetTokenAmount: "8900",
		StartDate:        fmt.Sprint(time.Now().Unix()),
		EndDate:          fmt.Sprint(time.Now().Unix()),
	}
	paramByte, _ := json.Marshal(assetIaoDto)
	iaoRes := mock.MockInvokeTransaction(suite.T(), suite.stub, [][]byte{[]byte("CreateIao"), paramByte})
	suite.T().Log(iaoRes)
	assert.NotEmptyf(suite.T(), iaoRes, "Create wallet return error", iaoRes)
	suite.IaoId = iaoRes
}

func (suite *IaoSCTestSuite) createNewAsset() {
	assetDto := iao.CreateAsset{
		Code:        "VIN123654",
		Name:        "A Smart city",
		OwnerWallet: suite.walletFromId,
		TokenName:   "HP Smart",
		TickerToken: "HPS",
		MaxSupply:   "78900",
		TotalValue:  "456321789",
		DocumentUrl: "http://xxx?yyy_zzz",
	}
	assetByte, _ := json.Marshal(assetDto)

	res := mock.MockInvokeTransaction(suite.T(), suite.stub, [][]byte{[]byte("CreateAsset"), assetByte})
	suite.T().Log(res)
	assert.NotEmpty(suite.T(), res, "Create wallet return empty")
	data := new(ResultAsset)
	err := json.Unmarshal([]byte(res), &data)
	assert.Nilf(suite.T(), err, "Parse data failed", err)
	suite.ATToken = data.TokenId
	suite.AssetId = data.AssetId
	suite.T().Log("ATToken", suite.ATToken)

	// Check if the created data exist in the ledger
	compositeKey, _ := suite.stub.CreateCompositeKey(doc.Tokens, helper.TokenKey(suite.ATToken))
	state, err := suite.stub.GetState(compositeKey)
	assert.Nilf(suite.T(), err, "Get token failed", err)

	tokenEntity := new(entity.Token)
	err = json.Unmarshal(state, &tokenEntity)
	assert.Nilf(suite.T(), err, "Parse token failed", err)
	assert.Equal(suite.T(), suite.ATToken, tokenEntity.Id)
	assert.Equal(suite.T(), glossary.Active, tokenEntity.Status)
}

func TestBaseSCTestSuite(t *testing.T) {
	suite.Run(t, new(IaoSCTestSuite))
}
