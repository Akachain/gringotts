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

package basic

import (
	"encoding/json"
	"github.com/Akachain/akc-go-sdk-v2/mock"
	"github.com/Akachain/gringotts/dto/token"
	"github.com/Akachain/gringotts/entity"
	"github.com/Akachain/gringotts/glossary"
	"github.com/Akachain/gringotts/glossary/doc"
	"github.com/Akachain/gringotts/helper"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

func setupMock() *mock.MockStubExtend {
	// Initialize MockStubExtend
	chaincodeName := "TMP"
	sc := NewBaseToken()
	chaincode, _ := contractapi.NewChaincode(sc)
	stub := mock.NewMockStubExtend(shimtest.NewMockStub(chaincodeName, chaincode), chaincode, ".")

	// Create a new database, Drop old database
	db, _ := mock.NewCouchDBHandler(false, chaincodeName)
	stub.SetCouchDBConfiguration(db)
	return stub
}

type BaseSCTestSuite struct {
	suite.Suite
	walletFromId string
	walletToId   string
	tokenId      string
	stub         *mock.MockStubExtend
}

func (suite *BaseSCTestSuite) SetupTest() {
	suite.stub = setupMock()

	// create token type
	stableToken := token.CreateTokenType{
		Name:        "Stable",
		TickerToken: "ST",
	}
	paramByte, _ := json.Marshal(stableToken)
	suite.tokenId = mock.MockInvokeTransaction(suite.T(), suite.stub, [][]byte{[]byte("CreateTokenType"), paramByte})
	suite.T().Log(suite.tokenId)
	assert.NotEmpty(suite.T(), suite.tokenId, "Create Token Type return empty")

	// create wallet
	wallet := token.CreateWallet{
		TokenId: suite.tokenId,
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
		TokenId:  suite.tokenId,
		Amount:   "100000",
	}
	mintByte, _ := json.Marshal(mintDto)
	mintRes := mock.MockInvokeTransaction(suite.T(), suite.stub, [][]byte{[]byte("Mint"), mintByte})
	assert.Empty(suite.T(), mintRes, "Mint invoke return err")

	// accounting balance
	suite.accountingBalance()
}

func (suite *BaseSCTestSuite) TestTokenBaseSC_CreateWallet() {
	walletFrom := token.CreateWallet{
		TokenId: suite.tokenId,
		Status:  "A",
	}
	walletByte, _ := json.Marshal(walletFrom)

	walletId := mock.MockInvokeTransaction(suite.T(), suite.stub, [][]byte{[]byte("CreateWallet"), walletByte})
	suite.T().Log(walletId)
	assert.NotEmpty(suite.T(), walletId, "Create wallet return empty")

	// Check if the created data exist in the ledger
	compositeKey, _ := suite.stub.CreateCompositeKey(doc.Wallets, helper.WalletKey(walletId))
	state, err := suite.stub.GetState(compositeKey)
	assert.Nilf(suite.T(), err, "Get wallet failed", err)

	walletEntity := new(entity.Wallet)
	err = json.Unmarshal(state, &walletEntity)
	assert.Nilf(suite.T(), err, "Parse wallet failed", err)
	assert.Equal(suite.T(), walletId, walletEntity.Id)
	assert.Equal(suite.T(), glossary.Active, walletEntity.Status)
}

func (suite *BaseSCTestSuite) TestTokenBaseSC_Transfer() {
	transferDto := token.TransferToken{
		FromWalletId: suite.walletFromId,
		ToWalletId:   suite.walletToId,
		TokenId:      suite.tokenId,
		Amount:       "100",
	}
	paramByte, _ := json.Marshal(transferDto)
	transferRes := mock.MockInvokeTransaction(suite.T(), suite.stub, [][]byte{[]byte("Transfer"), paramByte})
	suite.T().Log(transferRes)
	assert.Emptyf(suite.T(), transferRes, "Create wallet return error", transferRes)

	// accounting balance
	suite.accountingBalance()

	// get check balance
	balanceOfFromWallet := suite.getBalance(suite.walletFromId)
	suite.T().Log(balanceOfFromWallet)
	assert.NotEmptyf(suite.T(), balanceOfFromWallet, "Get balance of From wallet return empty", balanceOfFromWallet)
	assert.Equal(suite.T(), "9990000000000", balanceOfFromWallet, "Sub balance of From wallet failed")

	balanceOfToWallet := suite.getBalance(suite.walletToId)
	suite.T().Log(balanceOfToWallet)
	assert.NotEmptyf(suite.T(), balanceOfToWallet, "Get balance of To wallet return error", balanceOfToWallet)
	assert.Equal(suite.T(), "10000000000", balanceOfToWallet, "Sub balance of To wallet failed")
}

func (suite *BaseSCTestSuite) TestTokenBaseSC_UpdateWallet() {
	updateWalletDto := token.UpdateWallet{
		WalletId: suite.walletFromId,
		Status:   glossary.InActive,
	}
	paramByte, _ := json.Marshal(updateWalletDto)
	updateRes := mock.MockInvokeTransaction(suite.T(), suite.stub, [][]byte{[]byte("UpdateWallet"), paramByte})
	suite.T().Log(updateRes)
	assert.Emptyf(suite.T(), updateRes, "Update wallet return error", updateRes)
}

func (suite *BaseSCTestSuite) TestTokenBaseSC_Mint() {
	mintDto := token.MintToken{
		WalletId: suite.walletToId,
		Amount:   "200",
	}
	paramByte, _ := json.Marshal(mintDto)
	mintRes := mock.MockInvokeTransaction(suite.T(), suite.stub, [][]byte{[]byte("Mint"), paramByte})
	suite.T().Log(mintRes)
	assert.Emptyf(suite.T(), mintRes, "Mint token return error", mintRes)

	// accounting balance
	suite.accountingBalance()

	// checking balance
	balanceRes := suite.getBalance(suite.walletToId)
	assert.NotEmpty(suite.T(), balanceRes, "Get balance wallet return empty")
	assert.Equal(suite.T(), "20000000000", balanceRes, "Balance mint not enough")
}

func (suite *BaseSCTestSuite) TestTokenBaseSC_Burn() {
	burnDto := token.BurnToken{
		WalletId: suite.walletFromId,
		Amount:   "1000",
	}
	paramByte, _ := json.Marshal(burnDto)
	burnRes := mock.MockInvokeTransaction(suite.T(), suite.stub, [][]byte{[]byte("Burn"), paramByte})
	suite.T().Log(burnRes)
	assert.Emptyf(suite.T(), burnRes, "Burn token return error", burnRes)

	// accounting balance
	suite.accountingBalance()

	// checking balance
	balanceRes := suite.getBalance(suite.walletFromId)
	assert.NotEmpty(suite.T(), balanceRes, "Get balance wallet return empty")
	assert.Equal(suite.T(), "9900000000000", balanceRes, "Balance mint not enough")
}

func (suite *BaseSCTestSuite) TestTokenBaseSC_CreateTokenType() {
	tokenTypeDto := token.CreateTokenType{
		Name:        "Test Token",
		TickerToken: "TS",
	}
	paramByte, _ := json.Marshal(tokenTypeDto)
	createTokenRes := mock.MockInvokeTransaction(suite.T(), suite.stub, [][]byte{[]byte("CreateTokenType"), paramByte})
	suite.T().Log(createTokenRes)
	assert.NotEmptyf(suite.T(), createTokenRes, "Burn token return error", createTokenRes)
}

func TestBaseSCTestSuite(t *testing.T) {
	suite.Run(t, new(BaseSCTestSuite))
}

func (suite *BaseSCTestSuite) accountingBalance() {
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

func (suite *BaseSCTestSuite) getBalance(walletId string) string {
	balanceDto := token.Balance{
		WalletId: walletId,
	}
	paramByte, _ := json.Marshal(balanceDto)
	balanceOf := mock.MockInvokeTransaction(suite.T(), suite.stub, [][]byte{[]byte("GetBalance"), paramByte})
	suite.T().Log(balanceOf)

	return balanceOf
}
