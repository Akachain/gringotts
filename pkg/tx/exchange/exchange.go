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

package exchange

import (
	"github.com/Akachain/gringotts/entity"
	"github.com/Akachain/gringotts/glossary"
	"github.com/Akachain/gringotts/glossary/transaction"
	"github.com/Akachain/gringotts/helper/glogger"
	"github.com/Akachain/gringotts/pkg/tx"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
)

type txExchange struct {
	*tx.TxBase
}

func NewTxExchange() tx.Handler {
	return &txExchange{tx.NewTxBase()}
}

func (t *txExchange) AccountingTx(ctx contractapi.TransactionContextInterface, tx *entity.Transaction, mapBalanceToken map[string]string) (*entity.Transaction, error) {
	// TODO: handle rollback map current balance when sub/add failed
	if tx.FromWallet == glossary.SystemWallet || tx.ToWallet == glossary.SystemWallet {
		glogger.GetInstance().Errorf(ctx, "TxHandler - Exchange - Transaction (%s) has from/to wallet Id is system type", tx.Id)
		tx.Status = transaction.Rejected
		return tx, errors.New("From/To wallet id invalidate")
	}

	if err := t.SubAmount(ctx, mapBalanceToken, tx.FromWallet, tx.FromTokenId, tx.Amount); err != nil {
		glogger.GetInstance().Errorf(ctx, "TxHandler - Exchange - Transaction (%s): Unable to sub temp amount of From wallet", tx.Id)
		tx.Status = transaction.Rejected
		return tx, errors.WithMessage(err, "Sub balance of from wallet failed")
	}

	if err := t.AddAmount(ctx, mapBalanceToken, tx.ToWallet, tx.FromTokenId, tx.Amount); err != nil {
		glogger.GetInstance().Errorf(ctx, "TxHandler - Exchange - Transaction (%s): Unable to add temp amount of To wallet", tx.Id)
		tx.Status = transaction.Rejected
		return tx, errors.WithMessage(err, "Add balance of to wallet failed")
	}

	if err := t.SubAmount(ctx, mapBalanceToken, tx.ToWallet, tx.ToTokenId, tx.Amount); err != nil {
		glogger.GetInstance().Errorf(ctx, "TxHandler - Exchange - Transaction (%s): Unable to sub temp amount of To wallet", tx.Id)
		tx.Status = transaction.Rejected
		return tx, errors.WithMessage(err, "Sub balance of from wallet failed")
	}

	if err := t.AddAmount(ctx, mapBalanceToken, tx.FromWallet, tx.ToTokenId, tx.Amount); err != nil {
		glogger.GetInstance().Errorf(ctx, "TxHandler - Exchange - Transaction (%s): Unable to add temp amount of From wallet", tx.Id)
		tx.Status = transaction.Rejected
		return tx, errors.WithMessage(err, "Add balance of to wallet failed")
	}
	tx.Status = transaction.Confirmed
	return tx, nil
}