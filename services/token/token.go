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
	"github.com/Akachain/gringotts/entity"
	"github.com/Akachain/gringotts/errorcode"
	"github.com/Akachain/gringotts/glossary"
	"github.com/Akachain/gringotts/glossary/doc"
	"github.com/Akachain/gringotts/glossary/transaction"
	"github.com/Akachain/gringotts/helper"
	"github.com/Akachain/gringotts/helper/glogger"
	"github.com/Akachain/gringotts/pkg/unit"
	"github.com/Akachain/gringotts/services/base"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"strings"
)

type tokenService struct {
	*base.Base
}

func NewTokenService() *tokenService {
	return &tokenService{
		base.NewBase(),
	}
}

func (t *tokenService) TransferWithNote(ctx contractapi.TransactionContextInterface, fromWalletId, toWalletId string,
	amount float64, note string, tokenId string) (string, error) {
	glogger.GetInstance().Info(ctx, "-----------Token Service - TransferWithNote-----------")
	return t.transferToken(ctx, fromWalletId, toWalletId, amount, note, tokenId)
}

func (t *tokenService) Transfer(ctx contractapi.TransactionContextInterface, fromWalletId, toWalletId string, amount float64, tokenId string) (string, error) {
	glogger.GetInstance().Info(ctx, "-----------Token Service - Transfer-----------")
	return t.transferToken(ctx, fromWalletId, toWalletId, amount, "", tokenId)
}

func (t *tokenService) Mint(ctx contractapi.TransactionContextInterface, walletId string, amount float64, tokenId string) error {
	glogger.GetInstance().Info(ctx, "-----------Token Service - Mint-----------")

	// validate wallet exited
	if _, err := t.GetActiveWallet(ctx, walletId); err != nil {
		glogger.GetInstance().Errorf(ctx, "Mint - Get wallet mint failed with error (%s)", err.Error())
		return err
	}

	// convert balance to akc base
	amountUnit := unit.NewBalanceUnitFromFloat(amount)

	// create tx mint token
	txMint := entity.NewTransaction(ctx)
	txMint.FromWallet = glossary.SystemWallet
	txMint.ToWallet = walletId
	txMint.FromTokenId = tokenId
	txMint.ToTokenId = tokenId
	txMint.Amount = amountUnit.String()
	txMint.TxType = transaction.Mint

	if err := t.Repo.Create(ctx, txMint, doc.Transactions, helper.TransactionKey(txMint.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Mint - Create mint transaction failed with error (%s)", err.Error())
		return helper.RespError(errorcode.BizUnableCreateTX)
	}
	glogger.GetInstance().Infof(ctx, "-----------Token Service - Mint succeed (%s)-----------", txMint.Id)

	return nil
}

func (t *tokenService) Burn(ctx contractapi.TransactionContextInterface, walletId string, amount float64, tokenId string) error {
	glogger.GetInstance().Info(ctx, "-----------Token Service - Burn-----------")

	// validate burn wallet exist
	wallet, err := t.GetActiveWallet(ctx, walletId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Burn - Get wallet mint failed with error (%v)", err)
		return err
	}

	// get balance of token
	balanceToken, err := t.GetBalanceOfToken(ctx, wallet.Id, tokenId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Burn - Get balance of token failed with error (%v)", err)
		return err
	}

	// check balance enough to burn
	if helper.CompareFloatBalance(balanceToken.Balances, amount) < 0 {
		glogger.GetInstance().Error(ctx, "Burn - Wallet balance is insufficient", err)
		return helper.RespError(errorcode.BizBalanceNotEnough)
	}

	// convert balance to akc base
	amountUnit := unit.NewBalanceUnitFromFloat(amount)

	// create tx burn token
	txBurn := entity.NewTransaction(ctx)
	txBurn.SpenderWallet = walletId
	txBurn.FromWallet = walletId
	txBurn.ToWallet = glossary.SystemWallet
	txBurn.FromTokenId = tokenId
	txBurn.ToTokenId = tokenId
	txBurn.Amount = amountUnit.String()
	txBurn.TxType = transaction.Burn

	if err := t.Repo.Create(ctx, txBurn, doc.Transactions, helper.TransactionKey(txBurn.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Burn - Create burn transaction failed with error (%v)", err)
		return helper.RespError(errorcode.BizUnableCreateTX)
	}
	glogger.GetInstance().Infof(ctx, "-----------Token Service - Burn succeed (%s)-----------", txBurn.Id)

	return nil
}

func (t *tokenService) CreateType(ctx contractapi.TransactionContextInterface, name string, tickerToken string, rate float64, maxSupply string) (string, error) {
	glogger.GetInstance().Info(ctx, "-----------Token Service - CreateType-----------")

	tokenEntity := entity.NewToken(ctx)
	tokenEntity.Name = name
	tokenEntity.TickerToken = tickerToken
	tokenEntity.Rate = rate
	tokenEntity.MaxSupply = maxSupply

	if err := t.Repo.Create(ctx, tokenEntity, doc.Tokens, helper.TokenKey(tokenEntity.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "CreateType - Create token type failed with error (%s)", err.Error())
		return "", helper.RespError(errorcode.BizUnableCreateToken)
	}
	glogger.GetInstance().Infof(ctx, "-----------Token Service - CreateType succeed (%s)-----------", tokenEntity.Id)

	return tokenEntity.Id, nil
}

func (t *tokenService) Exchange(ctx contractapi.TransactionContextInterface, fromWalletId, toWalletId string,
	fromTokenId string, toTokenId string, amount float64) error {
	glogger.GetInstance().Info(ctx, "-----------Token Service - Exchange-----------")

	// validate from wallet and to wallet have active or not
	if _, _, err := t.ValidatePairWallet(ctx, fromWalletId, toWalletId); err != nil {
		glogger.GetInstance().Errorf(ctx, "Exchange - Validation swap failed with error (%s)", err.Error())
		return err
	}

	// TODO: validate token

	// convert balance to akc base
	amountUnit := unit.NewBalanceUnitFromFloat(amount)

	// create new swap transaction
	txEntity := entity.NewTransaction(ctx)
	txEntity.SpenderWallet = fromWalletId
	txEntity.FromWallet = fromWalletId
	txEntity.ToWallet = toWalletId
	txEntity.FromTokenId = fromTokenId
	txEntity.ToTokenId = toTokenId
	txEntity.TxType = transaction.Exchange
	txEntity.Amount = amountUnit.String()

	if err := t.Repo.Create(ctx, txEntity, doc.Transactions, helper.TransactionKey(txEntity.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Exchange - Create exchange transaction failed with error (%v)", err)
		return helper.RespError(errorcode.BizUnableCreateTX)
	}
	glogger.GetInstance().Infof(ctx, "-----------Token Service - Exchange succeed (%s)-----------", txEntity.Id)

	return nil
}

func (t *tokenService) Issue(ctx contractapi.TransactionContextInterface, fromTokenId, toTokenId, walletId string, amount float64) error {
	glogger.GetInstance().Info(ctx, "-----------Token Service - Issue-----------")

	// validate wallet active or not
	wallet, err := t.GetWallet(ctx, walletId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Issue - Get wallet failed with err (%s)", err.Error())
		return err
	}

	// check permission to issue new token
	enrollment, err := t.GetEnrollment(ctx, toTokenId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Issue - Get enrollment failed with error (%v)", err)
		return err
	}
	if enrollment.FromWalletId != "" {
		if !strings.Contains(enrollment.FromWalletId, wallet.Id) {
			glogger.GetInstance().Errorf(ctx, "Issue - From wallet do not have permission issue token (%s)", toTokenId)
			return helper.RespError(errorcode.BizIssueNotPermission)
		}
	}
	if enrollment.ToWalletId != "" {
		if !strings.Contains(enrollment.ToWalletId, wallet.Id) {
			glogger.GetInstance().Errorf(ctx, "Issue - To wallet do not have permission issue token (%s)", toTokenId)
			return helper.RespError(errorcode.BizIssueNotPermission)
		}
	}

	// calculate amount sub from wallet and add to wallet
	amountUpdate, err := t.calculateAmountSwap(ctx, wallet, fromTokenId, amount)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Issue - Calculate balance issue failed with error (%v)", err)
		return err
	}

	// create new swap transaction
	txEntity := entity.NewTransaction(ctx)
	txEntity.SpenderWallet = wallet.Id
	txEntity.FromWallet = wallet.Id
	txEntity.ToWallet = wallet.Id
	txEntity.FromTokenId = fromTokenId
	txEntity.ToTokenId = toTokenId
	txEntity.TxType = transaction.Issue
	txEntity.Amount = amountUpdate

	if err := t.Repo.Create(ctx, txEntity, doc.Transactions, helper.TransactionKey(txEntity.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Issue - Create swap transaction failed with error (%v)", err)
		return helper.RespError(errorcode.BizUnableCreateTX)
	}
	glogger.GetInstance().Infof(ctx, "-----------Token Service - Issue succeed (%s)-----------", txEntity.Id)

	return nil
}

func (t *tokenService) validateTransfer(ctx contractapi.TransactionContextInterface, fromWalletId, toWalletId string, tokenId string, amount float64) error {
	walletFrom, _, err := t.ValidatePairWallet(ctx, fromWalletId, toWalletId)
	if err != nil {
		return err
	}

	balanceToken, err := t.GetBalanceOfToken(ctx, walletFrom.Id, tokenId)
	if err != nil {
		return err
	}

	// check balance enough to transfer
	if helper.CompareFloatBalance(balanceToken.Balances, amount) < 0 {
		glogger.GetInstance().Error(ctx, "ValidateTransfer - Balance of from wallet is insufficient")
		return helper.RespError(errorcode.BizBalanceNotEnough)
	}
	return nil
}

func (t *tokenService) calculateAmountSwap(ctx contractapi.TransactionContextInterface, wallet *entity.Wallet, tokenId string, amount float64) (string, error) {
	tokenFrom, err := t.GetTokenType(ctx, tokenId)
	if err != nil {
		glogger.GetInstance().Error(ctx, "CalculateAmountSwap - Get token type failed with err (%v)", err)
	}

	amountUpdate := amount * tokenFrom.Rate
	amountUpdateUnit := unit.NewBalanceUnitFromFloat(amountUpdate)

	balanceToken, err := t.GetBalanceOfToken(ctx, wallet.Id, tokenId)
	if err != nil {
		return "-1", err
	}

	// check balance enough to transfer
	if helper.CompareStringBalance(balanceToken.Balances, amountUpdateUnit.String()) < 0 {
		glogger.GetInstance().Error(ctx, "CalculateAmountSwap - Balance of from wallet is insufficient")
		return "-1", helper.RespError(errorcode.BizBalanceNotEnough)
	}
	return amountUpdateUnit.String(), nil
}

func (t *tokenService) transferToken(ctx contractapi.TransactionContextInterface, fromWalletId, toWalletId string,
	amount float64, note string, tokenId string) (string, error) {
	if err := t.validateTransfer(ctx, fromWalletId, toWalletId, tokenId, amount); err != nil {
		glogger.GetInstance().Errorf(ctx, "Transfer - Validation transfer failed with error (%v)", err)
		return "", err
	}

	// convert amount to base unit
	amountUnit := unit.NewBalanceUnitFromFloat(amount)

	// create new transfer transaction
	txEntity := entity.NewTransaction(ctx)
	txEntity.SpenderWallet = fromWalletId
	txEntity.FromWallet = fromWalletId
	txEntity.ToWallet = toWalletId
	txEntity.FromTokenId = tokenId
	txEntity.ToTokenId = tokenId
	txEntity.TxType = transaction.Transfer
	txEntity.Amount = amountUnit.String()
	txEntity.Note = note

	if err := t.Repo.Create(ctx, txEntity, doc.Transactions, helper.TransactionKey(txEntity.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Transfer - Create transfer transaction failed with error (%v)", err)
		return "", helper.RespError(errorcode.BizUnableCreateTX)
	}
	glogger.GetInstance().Infof(ctx, "-----------Token Service - Transfer succeed (%s)-----------", txEntity.Id)

	return txEntity.Id, nil
}
