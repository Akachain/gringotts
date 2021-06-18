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
	"strings"
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
	err = db.ProcessIndexesForChaincodeDeploy("indexPendingTx.json", "../../META-INF/statedb/couchdb/indexes/indexPendingTx.json")
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
	mintDto := token.MintToken{
		WalletId: suite.walletFromId,
		TokenId:  suite.STToken,
		Amount:   "678900",
	}
	mintByte, _ := json.Marshal(mintDto)
	mintRes := mock.MockInvokeTransaction(suite.T(), suite.stub, [][]byte{[]byte("Mint"), mintByte})
	assert.Empty(suite.T(), mintRes, "Mint invoke return err")

	// accounting balance
	suite.accountingBalance()
}

func (suite *IaoSCTestSuite) TestIAOSC_CreateAsset() {
	suite.createNewAsset()
}

func (suite *IaoSCTestSuite) TestIAO_CreateIAO() {
	// create new asset
	suite.createNewAsset()

	// issue AT
	issueDto := token.IssueToken{
		WalletId:        suite.walletFromId,
		FromTokenId:     suite.STToken,
		ToTokenId:       suite.ATToken,
		FromTokenAmount: "78900",
		ToTokenAmount:   "78900",
	}

	paramByte, _ := json.Marshal(issueDto)
	issueRes := mock.MockInvokeTransaction(suite.T(), suite.stub, [][]byte{[]byte("Issue"), paramByte})
	suite.T().Log(issueRes)
	assert.Emptyf(suite.T(), issueRes, "Issue AT token return error", issueRes)

	// accounting balance
	suite.accountingBalance()

	assetIaoDto := iao.AssetIao{
		AssetId:          suite.AssetId,
		AssetTokenAmount: "8900",
		StartDate:        string(time.Now().Unix()),
		EndDate:          string(time.Now().Unix()),
		Rate:             1.2,
	}
	paramByte, _ = json.Marshal(assetIaoDto)
	iaoRes := mock.MockInvokeTransaction(suite.T(), suite.stub, [][]byte{[]byte("CreateIao"), paramByte})
	suite.T().Log(iaoRes)
	assert.NotEmptyf(suite.T(), iaoRes, "Create wallet return error", iaoRes)

	// accounting balance
	suite.accountingBalance()

	// get check balance
	balanceOfFromWallet := suite.getBalance(suite.walletFromId, suite.ATToken)
	suite.T().Log(balanceOfFromWallet)
	assert.NotEmptyf(suite.T(), balanceOfFromWallet, "Get balance of From wallet return empty", balanceOfFromWallet)
	assert.Equal(suite.T(), "70000", balanceOfFromWallet, "Sub balance of From wallet failed")
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

func (suite *IaoSCTestSuite) accountingBalance() {
	lstTx := mock.MockInvokeTransaction(suite.T(), suite.stub, [][]byte{[]byte("GetAccountingTx")})
	suite.T().Log(lstTx)

	lstTx = strings.ReplaceAll(lstTx, "[", "")
	lstTx = strings.ReplaceAll(lstTx, "]", "")
	lstTx = strings.ReplaceAll(lstTx, "\"", "")
	suite.T().Log(lstTx)
	// accounting
	accountingDto := token.AccountingBalance{
		TxId: strings.Split(lstTx, ","),
	}
	paramByte, _ := json.Marshal(accountingDto)
	accountingRes := mock.MockInvokeTransaction(suite.T(), suite.stub, [][]byte{[]byte("CalculateBalance"), paramByte})
	assert.Empty(suite.T(), accountingRes, "CalculateBalance invoke return err")
}

func (suite *IaoSCTestSuite) getBalance(walletId, tokenId string) string {
	balanceDto := token.Balance{
		WalletId: walletId,
		TokenId:  tokenId,
	}
	paramByte, _ := json.Marshal(balanceDto)
	balanceOf := mock.MockInvokeTransaction(suite.T(), suite.stub, [][]byte{[]byte("GetBalance"), paramByte})
	suite.T().Log(balanceOf)

	return balanceOf
}
