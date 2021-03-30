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

package token

import (
	"github.com/Akachain/gringotts/dto"
	"github.com/Akachain/gringotts/errorcode"
	"github.com/Akachain/gringotts/glossary"
	"github.com/Akachain/gringotts/glossary/doc"
	"github.com/Akachain/gringotts/glossary/transaction"
	"github.com/Akachain/gringotts/helper"
	"github.com/Akachain/gringotts/helper/glogger"
	"github.com/Akachain/gringotts/internal/entity"
	"github.com/Akachain/gringotts/pkg/unit"
	"github.com/Akachain/gringotts/services/base"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type tokenService struct {
	*base.Base
}

func NewTokenService() *tokenService {
	return &tokenService{
		base.NewBase(),
	}
}

func (t *tokenService) Transfer(ctx contractapi.TransactionContextInterface, transferDto dto.TransferToken) error {
	glogger.GetInstance().Info(ctx, "-----------Token Service - Transfer-----------")

	if err := t.validateTransfer(ctx, transferDto.FromWalletId, transferDto.ToWalletId, transferDto.Amount); err != nil {
		glogger.GetInstance().Errorf(ctx, "Transfer - Validation transfer failed with error (%v)", err)
		return err
	}

	// create new transfer transaction
	txEntity := transferDto.ToEntity(ctx)
	if err := t.Repo.Create(ctx, txEntity, doc.Transactions, helper.TransactionKey(txEntity.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Transfer - Create transfer transaction failed with error (%v)", err)
		return helper.RespError(errorcode.BizUnableCreateTX)
	}
	glogger.GetInstance().Infof(ctx, "-----------Token Service - Transfer succeed (%s)-----------", txEntity.Id)

	return nil
}

func (t *tokenService) Mint(ctx contractapi.TransactionContextInterface, mintDto dto.MintToken) error {
	glogger.GetInstance().Info(ctx, "-----------Token Service - Mint-----------")

	// validate wallet exited
	wallet, err := t.GetWallet(ctx, mintDto.WalletId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Mint - Get wallet mint failed with error (%v)", err)
		return err
	}
	if wallet.Status != glossary.Active {
		glogger.GetInstance().Errorf(ctx, "Burn - Wallet has status inactive", err)
		return helper.RespError(errorcode.InvalidParam)
	}

	// convert balance to akc base
	amountUnit := unit.NewBalanceUnitFromFloat(mintDto.Amount)

	// create tx mint token
	txMint := entity.NewTransaction(ctx)
	txMint.FromWallet = glossary.SystemWallet
	txMint.ToWallet = mintDto.WalletId
	txMint.Amount = amountUnit.String()
	txMint.TxType = transaction.Mint
	txMint.Status = transaction.Pending

	if err := t.Repo.Create(ctx, txMint, doc.Transactions, helper.TransactionKey(txMint.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Mint - Create mint transaction failed with error (%v)", err)
		return helper.RespError(errorcode.BizUnableCreateTX)
	}
	glogger.GetInstance().Infof(ctx, "-----------Token Service - Mint succeed (%s)-----------", txMint.Id)

	return nil
}

func (t *tokenService) Burn(ctx contractapi.TransactionContextInterface, burnDto dto.BurnToken) error {
	glogger.GetInstance().Info(ctx, "-----------Token Service - Burn-----------")

	// validate burn wallet exist
	wallet, err := t.GetWallet(ctx, burnDto.WalletId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Burn - Get wallet mint failed with error (%v)", err)
		return err
	}
	if wallet.Status != glossary.Active {
		glogger.GetInstance().Errorf(ctx, "Burn - Wallet has status inactive", err)
		return helper.RespError(errorcode.InvalidParam)
	}

	// check balance enough to burn
	if helper.CompareFloatBalance(wallet.Balances, burnDto.Amount) < 1 {
		glogger.GetInstance().Error(ctx, "Burn - Wallet balance is insufficient", err)
		return helper.RespError(errorcode.BizBalanceNotEnough)
	}

	// convert balance to akc base
	amountUnit := unit.NewBalanceUnitFromFloat(burnDto.Amount)

	// create tx burn token
	txBurn := entity.NewTransaction(ctx)
	txBurn.FromWallet = burnDto.WalletId
	txBurn.ToWallet = glossary.SystemWallet
	txBurn.Amount = amountUnit.String()
	txBurn.TxType = transaction.Burn
	txBurn.Status = transaction.Pending

	if err := t.Repo.Create(ctx, txBurn, doc.Transactions, helper.TransactionKey(txBurn.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Burn - Create burn transaction failed with error (%v)", err)
		return helper.RespError(errorcode.BizUnableCreateTX)
	}
	glogger.GetInstance().Infof(ctx, "-----------Token Service - Burn succeed (%s)-----------", txBurn.Id)

	return nil
}

func (t *tokenService) CreateType(ctx contractapi.TransactionContextInterface, tokenType dto.CreateTokenType) (string, error) {
	glogger.GetInstance().Info(ctx, "-----------Token Service - CreateType-----------")

	tokenEntity := tokenType.ToEntity(ctx)
	if err := t.Repo.Create(ctx, tokenEntity, doc.Tokens, helper.TokenKey(tokenEntity.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "CreateType - Create token type failed with error (%v)", err)
		return "", helper.RespError(errorcode.BizUnableCreateToken)
	}
	glogger.GetInstance().Infof(ctx, "-----------Token Service - CreateType succeed (%s)-----------", tokenEntity.Id)

	return tokenEntity.Id, nil
}

func (t *tokenService) validateTransfer(ctx contractapi.TransactionContextInterface, fromWalletId, toWalletId string, amount float64) error {
	// validate to wallet exist
	toWallet, err := t.GetWallet(ctx, toWalletId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Transfer - Get wallet to failed with error (%v)", err)
		return err
	}
	if toWallet.Status != glossary.Active {
		glogger.GetInstance().Errorf(ctx, "Transfer - To Wallet has status inactive", err)
		return helper.RespError(errorcode.InvalidParam)
	}

	// validate from wallet exist
	walletFrom, err := t.GetWallet(ctx, fromWalletId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Transfer - Get wallet from failed with error (%v)", err)
		return err
	}
	if walletFrom.Status != glossary.Active {
		glogger.GetInstance().Errorf(ctx, "Transfer - From Wallet has status inactive", err)
		return helper.RespError(errorcode.InvalidParam)
	}

	// check balance enough to transfer
	if helper.CompareFloatBalance(walletFrom.Balances, amount) < 1 {
		glogger.GetInstance().Error(ctx, "Transfer - Balance of from wallet is insufficient", err)
		return helper.RespError(errorcode.BizBalanceNotEnough)
	}
	return nil
}
