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

package sidechain_transfer

import (
	"github.com/Akachain/gringotts/entity"
	"github.com/Akachain/gringotts/glossary"
	"github.com/Akachain/gringotts/glossary/doc"
	"github.com/Akachain/gringotts/glossary/sidechain"
	"github.com/Akachain/gringotts/glossary/transaction"
	"github.com/Akachain/gringotts/helper"
	"github.com/Akachain/gringotts/helper/glogger"
	"github.com/Akachain/gringotts/pkg/tx/base"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"strings"
)

type txSideChainTransfer struct {
	*base.TxBase
}

func NewTxSideChainTransfer() *txSideChainTransfer {
	return &txSideChainTransfer{
		base.NewTxBase(),
	}
}

func (t *txSideChainTransfer) AccountingTx(ctx contractapi.TransactionContextInterface, tx *entity.Transaction, mapBalanceToken map[string]*entity.BalanceCache) (*entity.Transaction, error) {
	if tx.FromWallet == glossary.SystemWallet || tx.ToWallet == glossary.SystemWallet {
		glogger.GetInstance().Errorf(ctx, "TxSideChainTransfer - Transaction (%s) has from/to wallet Id is system type", tx.Id)
		tx.Status = transaction.Rejected
		return tx, errors.New("From/To wallet id invalidate")
	}

	lstNameChain := strings.Split(tx.Note, "_")
	if err := t.SubAmount(ctx, mapBalanceToken, tx.FromWallet, tx.FromTokenId, tx.FromTokenAmount); err != nil {
		glogger.GetInstance().Errorf(ctx, "TxSideChainTransfer- Transaction (%s): Unable to sub temp amount of From wallet", tx.Id)
		tx.Status = transaction.Rejected
		return tx, errors.WithMessage(err, "Sub balance of from wallet failed")
	}

	isExisted, balanceData, err := t.Repo.GetAndCheckExist(ctx, doc.IaoBalances, helper.BalanceKey(tx.ToWallet, tx.ToTokenId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "TxSideChainTransfer - Get balance of token failed with error (%s)", err.Error())
		if err := t.RollbackTxHandler(ctx, tx, mapBalanceToken, transaction.SubFromWallet); err != nil {
			glogger.GetInstance().Errorf(ctx, "TxSideChainTransfer - Rollback handle transaction (%s) failed with error (%s)", tx.Id, err.Error())
		}
		tx.Status = transaction.Rejected
		return tx, errors.WithMessage(err, "Get balance of to wallet failed")
	}

	if isExisted {
		balanceToWallet := new(entity.Balance)
		if err = mapstructure.Decode(balanceData, &balanceToWallet); err != nil {
			glogger.GetInstance().Errorf(ctx, "TxSideChainTransfer - Decode balance of token failed with error (%s)", err.Error())
			if err := t.RollbackTxHandler(ctx, tx, mapBalanceToken, transaction.SubFromWallet); err != nil {
				glogger.GetInstance().Errorf(ctx, "TxSideChainTransfer - Rollback handle transaction (%s) failed with error (%s)", tx.Id, err.Error())
			}
			tx.Status = transaction.Rejected
			return tx, errors.WithMessage(err, "Get balance of to wallet failed")
		}
		updateBalance, err := helper.AddBalance(balanceToWallet.Balances, tx.ToTokenAmount)
		if err != nil {
			glogger.GetInstance().Errorf(ctx, "TxSideChainTransfer - Get balance of token failed with error (%s)", err.Error())
			if err := t.RollbackTxHandler(ctx, tx, mapBalanceToken, transaction.SubFromWallet); err != nil {
				glogger.GetInstance().Errorf(ctx, "TxSideChainTransfer - Rollback handle transaction (%s) failed with error (%s)", tx.Id, err.Error())
			}
			tx.Status = transaction.Rejected
			return tx, errors.WithMessage(err, "Add balance of to wallet failed")
		}
		balanceToWallet.Balances = updateBalance
		if err := t.Repo.Update(ctx, balanceToWallet, doc.IaoBalances, helper.BalanceKey(tx.ToWallet, tx.ToTokenId)); err != nil {
			glogger.GetInstance().Errorf(ctx, "TxSideChainTransfer - Update balance failed with err (%s)", err.Error())
			if err := t.RollbackTxHandler(ctx, tx, mapBalanceToken, transaction.SubFromWallet); err != nil {
				glogger.GetInstance().Errorf(ctx, "TxSideChainTransfer - Rollback handle transaction (%s) failed with error (%s)", tx.Id, err.Error())
			}
			tx.Status = transaction.Rejected
			return tx, err
		}
	}

	balanceEntity := entity.NewBalance(sidechain.SideName(lstNameChain[1]), ctx)
	balanceEntity.WalletId = tx.ToWallet
	balanceEntity.TokenId = tx.ToTokenId
	balanceEntity.Balances = tx.ToTokenAmount

	if err := t.Repo.Create(ctx, balanceEntity, doc.IaoBalances, helper.BalanceKey(tx.ToWallet, tx.ToTokenId)); err != nil {
		glogger.GetInstance().Errorf(ctx, "TxSideChainTransfer - Create balance failed with err (%s)", err.Error())
		if err := t.RollbackTxHandler(ctx, tx, mapBalanceToken, transaction.SubFromWallet); err != nil {
			glogger.GetInstance().Errorf(ctx, "TxSideChainTransfer - Rollback handle transaction (%s) failed with error (%s)", tx.Id, err.Error())
		}
		tx.Status = transaction.Rejected
		return tx, err
	}
	tx.Status = transaction.Confirmed

	return tx, nil
}
