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

package dto

import (
	"errors"
	"github.com/Akachain/gringotts/entity"
	"github.com/Akachain/gringotts/glossary/doc"
	"github.com/Akachain/gringotts/glossary/transaction"
	"github.com/Akachain/gringotts/helper"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type TransferToken struct {
	FromWalletId string `json:"fromWalletId"`
	ToWalletId   string `json:"toWalletId"`
	TokenId      string `json:"tokenId"`
	Amount       string `json:"amount"`
}

func (t TransferToken) ToEntity(ctx contractapi.TransactionContextInterface) *entity.Transaction {
	txTime, _ := ctx.GetStub().GetTxTimestamp()
	transactionEntity := new(entity.Transaction)
	transactionEntity.Id = helper.GenerateID(doc.Transactions, ctx.GetStub().GetTxID())
	transactionEntity.TxType = transaction.Transfer
	transactionEntity.Status = transaction.Pending
	transactionEntity.SpenderWallet = t.FromWalletId
	transactionEntity.FromWallet = t.FromWalletId
	transactionEntity.ToWallet = t.ToWalletId
	transactionEntity.FromTokenAmount = t.Amount
	transactionEntity.ToTokenAmount = t.Amount
	transactionEntity.CreatedAt = helper.TimestampISO(txTime.Seconds)
	transactionEntity.UpdatedAt = helper.TimestampISO(txTime.Seconds)
	transactionEntity.BlockChainId = ctx.GetStub().GetTxID()

	return transactionEntity
}

func (t TransferToken) IsValid() error {
	if t.FromWalletId == "" || t.ToWalletId == "" {
		return errors.New("From/To wallet id is empty")
	}

	if t.TokenId == "" {
		return errors.New("token id is empty")
	}

	if t.Amount == "" {
		return errors.New("the transfer amount is empty")
	}
	return nil
}
