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
	"github.com/Akachain/gringotts/dto"
	"github.com/Akachain/gringotts/helper"
	"github.com/Akachain/gringotts/smartcontract"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"gotest.tools/assert"
	"testing"
	"time"
)

func setupMock() (*contractapi.TransactionContext, smartcontract.BasicToken) {
	ts := new(timestamp.Timestamp)
	ts.Seconds = time.Now().Unix()
	sc := NewBaseToken()
	chaincode, _ := contractapi.NewChaincode(sc)
	stub := shimtest.NewMockStub("TMP", chaincode)
	stub.TxTimestamp = ts
	stub.TxID = helper.GenerateID("Test", "asd")

	txContext := new(contractapi.TransactionContext)
	txContext.SetStub(stub)

	return txContext, sc
}

func TestTokenBaseSC_Transfer(t *testing.T) {
	ctx, sc := setupMock()

	// create from wallet
	walletFrom := dto.CreateWallet{
		TokenId: "We123ea",
		Status:  "A",
	}
	walletFromId, err := sc.CreateWallet(ctx, walletFrom)
	assert.NilError(t, err, "Fail to create from wallet")

	walletTo := dto.CreateWallet{
		TokenId: "We123ea",
		Status:  "A",
	}
	walletToId, err := sc.CreateWallet(ctx, walletTo)
	assert.NilError(t, err, "Fail to create to wallet")

	transferDto := dto.TransferToken{
		FromWalletId: walletFromId,
		ToWalletId:   walletToId,
		Amount:       100,
	}
	err = sc.Transfer(ctx, transferDto)
	assert.NilError(t, err, "Fail to create to transfer token")
}

func TestTokenBaseSC_CreateWallet(t *testing.T) {
	ctx, sc := setupMock()

	walletFrom := dto.CreateWallet{
		TokenId: "We123ea",
		Status:  "A",
	}
	walletId, err := sc.CreateWallet(ctx, walletFrom)
	assert.NilError(t, err, "Fail to create wallet")

	balance, err := sc.GetBalance(ctx, dto.Balance{WalletId: walletId})
	assert.NilError(t, err, "Fail to get balance of wallet")
	assert.Equal(t, balance, "0")
}

func TestTokenBaseSC_UpdateWallet(t *testing.T) {
	ctx, sc := setupMock()

	walletFrom := dto.CreateWallet{
		TokenId: "We123ea",
		Status:  "A",
	}
	walletId, err := sc.CreateWallet(ctx, walletFrom)
	assert.NilError(t, err, "Fail to create wallet")

	balance, err := sc.GetBalance(ctx, dto.Balance{WalletId: walletId})
	assert.NilError(t, err, "Fail to get balance of wallet")
	assert.Equal(t, balance, "0")

	err = sc.UpdateWallet(ctx, dto.UpdateWallet{
		WalletId: walletId,
		Status:   "I",
	})
	assert.NilError(t, err, "Fail to update status of wallet")
}

func TestTokenBaseSC_Mint(t *testing.T) {
	ctx, sc := setupMock()

	walletFrom := dto.CreateWallet{
		TokenId: "We123ea",
		Status:  "A",
	}
	walletId, err := sc.CreateWallet(ctx, walletFrom)
	assert.NilError(t, err, "Fail to create wallet")

	balance, err := sc.GetBalance(ctx, dto.Balance{WalletId: walletId})
	assert.NilError(t, err, "Fail to get balance of wallet")
	assert.Equal(t, balance, "0")

	err = sc.Mint(ctx, dto.MintToken{
		WalletId: walletId,
		Amount:   321,
	})
	assert.NilError(t, err, "Fail to mint token for wallet")
}

func TestTokenBaseSC_Burn(t *testing.T) {
	ctx, sc := setupMock()

	// create wallet
	walletFrom := dto.CreateWallet{
		TokenId: "We123ea",
		Status:  "A",
	}
	walletId, err := sc.CreateWallet(ctx, walletFrom)
	assert.NilError(t, err, "Fail to create wallet")

	balance, err := sc.GetBalance(ctx, dto.Balance{WalletId: walletId})
	assert.NilError(t, err, "Fail to get balance of wallet")
	assert.Equal(t, balance, "0")

	// burn token
	err = sc.Burn(ctx, dto.BurnToken{
		WalletId: walletId,
		Amount:   100,
	})
	assert.NilError(t, err, "Fail to mint token for wallet")
}

func TestTokenBaseSC_CreateTokenType(t *testing.T) {
	ctx, sc := setupMock()

	tokenId, err := sc.CreateTokenType(ctx, dto.CreateTokenType{
		Name: "Stable",
		Rate: 0.123,
	})
	assert.NilError(t, err, "Fail to create token type")
	t.Log(tokenId)
	assert.Check(t, tokenId != "", "Token Id return empty")
}
