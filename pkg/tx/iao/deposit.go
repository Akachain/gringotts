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
	"github.com/Akachain/gringotts/entity"
	"github.com/Akachain/gringotts/glossary/doc"
	"github.com/Akachain/gringotts/glossary/iao"
	"github.com/Akachain/gringotts/glossary/transaction"
	"github.com/Akachain/gringotts/helper"
	"github.com/Akachain/gringotts/helper/glogger"
	"github.com/Akachain/gringotts/pkg/tx/base"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
)

type txDeposit struct {
	*base.TxBase
}

func NewTxDeposit() *txDeposit {
	return &txDeposit{
		base.NewTxBase(),
	}
}

func (t *txDeposit) AccountingTx(ctx contractapi.TransactionContextInterface, tx *entity.Transaction, mapBalanceToken map[string]*entity.BalanceCache) (*entity.Transaction, error) {
	txTime, _ := ctx.GetStub().GetTxTimestamp()
	// case tx is iao transfer, note is iao id
	iaoEntity, err := t.GetIao(ctx, tx.Note)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Iao TxTransfer - Get Iao failed with err (%s)", err.Error())
		tx.Status = transaction.Rejected
		return tx, err
	}

	if err := t.SubAmount(ctx, mapBalanceToken, tx.FromWallet, tx.FromTokenId, tx.FromTokenAmount); err != nil {
		glogger.GetInstance().Errorf(ctx, "Iao TxTransfer - Transaction (%s): Unable to sub temp amount of From wallet", tx.Id)
		tx.Status = transaction.Rejected
		return tx, errors.WithMessage(err, "Sub balance of from wallet failed")
	}

	assetTokenUpdate, err := helper.AddBalance(iaoEntity.AssetTokenAmount, tx.ToTokenAmount)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Iao TxTransfer - Transaction (%s): Unable to add balance of iao wallet", tx.Id, err.Error())
		tx.Status = transaction.Rejected
		return tx, errors.WithMessage(err, "Add balance of iao wallet failed")
	}

	iaoEntity.AssetTokenAmount = assetTokenUpdate
	iaoEntity.Status = iao.Open
	iaoEntity.UpdatedAt = helper.TimestampISO(txTime.Seconds)

	if err := t.Repo.Update(ctx, iaoEntity, doc.Iao, helper.IaoKey(iaoEntity.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Iao TxTransfer - Update IAO failed with err(%s)", err.Error())
		tx.Status = transaction.Rejected
		if err := t.RollbackTxHandler(ctx, tx, mapBalanceToken, transaction.SubFromWallet); err != nil {
			glogger.GetInstance().Errorf(ctx, "Iao TxTransfer - Rollback handle transaction (%s) failed with error (%v)", tx.Id, err)
		}
		return tx, errors.WithMessage(err, "Add balance of iao wallet failed")
	}

	// update remaining asset token
	assetEntity, err := t.GetAsset(ctx, iaoEntity.AssetId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Iao TxTransfer - Transaction (%s): Unable to get asset (%s)", tx.Id, err.Error())
		tx.Status = transaction.Rejected
		return tx, errors.WithMessage(err, "Get asset failed")
	}
	remainToken, err := helper.SubBalance(assetEntity.RemainingToken, tx.ToTokenAmount)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Iao Service - Sub amount token of asset failed with err (%s)", err.Error())
		tx.Status = transaction.Rejected
		return tx, errors.WithMessage(err, "Sub remaining asset token of asset failed")
	}
	assetEntity.RemainingToken = remainToken
	if err := t.Repo.Update(ctx, assetEntity, doc.Asset, helper.AssetKey(assetEntity.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Iao Service - Update asset failed with error (%s)", err.Error())
		tx.Status = transaction.Rejected
		return tx, errors.WithMessage(err, "Update remaining asset token of asset failed")
	}

	tx.Status = transaction.Confirmed
	return tx, nil
}
