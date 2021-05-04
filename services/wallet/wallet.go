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
	"github.com/Akachain/gringotts/entity"
	"github.com/Akachain/gringotts/errorcode"
	"github.com/Akachain/gringotts/glossary"
	"github.com/Akachain/gringotts/glossary/doc"
	"github.com/Akachain/gringotts/helper"
	"github.com/Akachain/gringotts/helper/glogger"
	"github.com/Akachain/gringotts/services/base"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"strings"
)

type walletService struct {
	*base.Base
}

func NewWalletService() *walletService {
	return &walletService{base.NewBase()}
}

func (w *walletService) Create(ctx contractapi.TransactionContextInterface, tokenId string, status glossary.Status) (string, error) {
	glogger.GetInstance().Info(ctx, "-----------Wallet Service - Create-----------")

	// validate with token type id
	token, err := w.GetTokenType(ctx, tokenId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Create - Get token type failed with error (%v)", err)
		return "", err
	}
	if token.Status != glossary.Active {
		glogger.GetInstance().Errorf(ctx, "Create - Token has status inactive", err)
		return "", helper.RespError(errorcode.InvalidParam)
	}

	walletEntity := entity.NewWallet(ctx)
	walletEntity.Status = status
	walletEntity.TokenId = tokenId
	if err := w.Repo.Create(ctx, walletEntity, doc.Wallets, helper.WalletKey(walletEntity.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Create - Create wallet failed with error (%v)", err)
		return "", helper.RespError(errorcode.BizUnableCreateWallet)
	}
	glogger.GetInstance().Infof(ctx, "-----------Wallet Service - Create Succeed: id (%s)-----------", walletEntity.Id)

	return walletEntity.Id, nil
}

func (w *walletService) Update(ctx contractapi.TransactionContextInterface, walletId string, status glossary.Status) error {
	glogger.GetInstance().Info(ctx, "-----------Wallet Service - Update-----------")
	txTime, _ := ctx.GetStub().GetTxTimestamp()

	wallet, err := w.GetWallet(ctx, walletId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Update - Get wallet failed with error (%v)", err)
		return helper.RespError(errorcode.BizUnableGetWallet)
	}
	if wallet.Status != glossary.Active {
		glogger.GetInstance().Errorf(ctx, "Update - Wallet has status inactive", err)
		return helper.RespError(errorcode.InvalidParam)
	}

	wallet.Status = status
	wallet.UpdatedAt = helper.TimestampISO(txTime.Seconds)

	if err := w.Repo.Update(ctx, wallet, doc.Wallets, helper.WalletKey(wallet.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Update - Update wallet failed with error (%v)", err)
		return helper.RespError(errorcode.BizUnableUpdateWallet)
	}
	glogger.GetInstance().Info(ctx, "-----------Wallet Service - Update Succeed-----------")

	return nil
}

func (w *walletService) BalanceOf(ctx contractapi.TransactionContextInterface, walletId string) (string, error) {
	glogger.GetInstance().Info(ctx, "-----------Wallet Service - BalanceOf-----------")

	wallet, err := w.GetWallet(ctx, walletId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "BalanceOf - Get wallet failed with error (%v)", err)
		return "-1", helper.RespError(errorcode.BizUnableGetWallet)
	}
	glogger.GetInstance().Infof(ctx, "-----------Wallet Service - BalanceOf wallet: (%s)-----------", wallet.Balances)

	return wallet.Balances, nil
}

func (w *walletService) EnrollToken(ctx contractapi.TransactionContextInterface, tokenId string, fromWalletId []string, toWalletId []string) error {
	glogger.GetInstance().Info(ctx, "-----------Wallet Service - EnrollToken-----------")
	isExisted, err := w.Repo.IsExist(ctx, doc.Enrollments, helper.EnrollmentKey(tokenId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Wallet Service - Check exit enrollment failed with err (%v)", err)
		return helper.RespError(errorcode.BizUnableCreateEnrollment)
	}

	if isExisted {
		enrollment, err := w.GetEnrollment(ctx, tokenId)
		if err != nil {
			glogger.GetInstance().Errorf(ctx, "Wallet Service - Get enrollment failed with err (%v)", err)
			return err
		}
		if len(fromWalletId) > 0 && !helper.ArrayContains(fromWalletId, "") {
			enrollment.FromWalletId = enrollment.FromWalletId + "," + strings.Join(fromWalletId, ",")
		}
		if len(toWalletId) > 0 && !helper.ArrayContains(toWalletId, "") {
			enrollment.ToWalletId = enrollment.ToWalletId + "," + strings.Join(toWalletId, ",")
		}
		if err := w.Repo.Update(ctx, enrollment, doc.Enrollments, helper.EnrollmentKey(tokenId)); err != nil {
			glogger.GetInstance().Errorf(ctx, "Wallet Service - Update enrollment failed with err (%v)", err)
			return helper.RespError(errorcode.BizUnableCreateEnrollment)
		}
	} else {
		enrollmentEntity := entity.NewEnrollment(ctx)
		if len(fromWalletId) > 0 && !helper.ArrayContains(fromWalletId, "") {
			enrollmentEntity.FromWalletId = strings.Join(fromWalletId, ",")
		}
		if len(toWalletId) > 0 && !helper.ArrayContains(toWalletId, "") {
			enrollmentEntity.ToWalletId = strings.Join(toWalletId, ",")
		}
		if err := w.Repo.Create(ctx, enrollmentEntity, doc.Enrollments, helper.EnrollmentKey(tokenId)); err != nil {
			glogger.GetInstance().Errorf(ctx, "Wallet Service - Create enrollment failed with err (%v)", err)
			return helper.RespError(errorcode.BizUnableCreateEnrollment)
		}
	}
	return nil
}
