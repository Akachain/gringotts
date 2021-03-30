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

package wallet

import (
	"github.com/Akachain/gringotts/dto"
	"github.com/Akachain/gringotts/errorcode"
	"github.com/Akachain/gringotts/glossary"
	"github.com/Akachain/gringotts/glossary/doc"
	"github.com/Akachain/gringotts/helper"
	"github.com/Akachain/gringotts/helper/glogger"
	"github.com/Akachain/gringotts/services/base"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type walletService struct {
	*base.Base
}

func NewWalletService() *walletService {
	return &walletService{base.NewBase()}
}

func (w *walletService) Create(ctx contractapi.TransactionContextInterface, createWalletDto dto.CreateWallet) (string, error) {
	glogger.GetInstance().Info(ctx, "-----------Wallet Service - Create-----------")

	// validate with token type id
	token, err := w.GetTokenType(ctx, createWalletDto.TokenId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Create - Get token type failed with error (%v)", err)
		return "", err
	}
	if token.Status != glossary.Active {
		glogger.GetInstance().Errorf(ctx, "Create - Token has status inactive", err)
		return "", helper.RespError(errorcode.InvalidParam)
	}

	walletEntity := createWalletDto.ToEntity(ctx)
	if err := w.Repo.Create(ctx, walletEntity, doc.Wallets, helper.WalletKey(walletEntity.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Create - Create wallet failed with error (%v)", err)
		return "", helper.RespError(errorcode.BizUnableCreateWallet)
	}
	glogger.GetInstance().Infof(ctx, "-----------Wallet Service - Create Succeed: id (%s)-----------", walletEntity.Id)

	return walletEntity.Id, nil
}

func (w *walletService) Update(ctx contractapi.TransactionContextInterface, updateWalletDto dto.UpdateWallet) error {
	glogger.GetInstance().Info(ctx, "-----------Wallet Service - Update-----------")
	txTime, _ := ctx.GetStub().GetTxTimestamp()

	wallet, err := w.GetWallet(ctx, updateWalletDto.WalletId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Update - Get wallet failed with error (%v)", err)
		return helper.RespError(errorcode.BizUnableGetWallet)
	}
	if wallet.Status != glossary.Active {
		glogger.GetInstance().Errorf(ctx, "Update - Wallet has status inactive", err)
		return helper.RespError(errorcode.InvalidParam)
	}

	wallet.Status = updateWalletDto.Status
	wallet.UpdatedAt = helper.TimestampISO(txTime.Seconds)

	if err := w.Repo.Update(ctx, wallet, doc.Wallets, helper.WalletKey(wallet.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Update - Update wallet failed with error (%v)", err)
		return helper.RespError(errorcode.BizUnableUpdateWallet)
	}
	glogger.GetInstance().Info(ctx, "-----------Wallet Service - Update Succeed-----------")

	return nil
}

func (w *walletService) BalanceOf(ctx contractapi.TransactionContextInterface, balanceDto dto.Balance) (string, error) {
	glogger.GetInstance().Info(ctx, "-----------Wallet Service - BalanceOf-----------")

	wallet, err := w.GetWallet(ctx, balanceDto.WalletId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "BalanceOf - Get wallet failed with error (%v)", err)
		return "-1", helper.RespError(errorcode.BizUnableGetWallet)
	}
	glogger.GetInstance().Infof(ctx, "-----------Wallet Service - BalanceOf wallet: (%s)-----------", wallet.Balances)

	return wallet.Balances, nil
}
