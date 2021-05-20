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

package entity

import (
	"github.com/Akachain/gringotts/glossary/doc"
	"github.com/Akachain/gringotts/glossary/transaction"
	"github.com/Akachain/gringotts/helper"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Exchange struct {
	OwnerWalletId string
	ToWalletId    string
	NftTokenId    string
	Price         float64
	Status        transaction.Status
	Base          `mapstructure:",squash"`
}

func NewExchange(ctx ...contractapi.TransactionContextInterface) *Exchange {
	if len(ctx) <= 0 {
		return &Exchange{}
	}
	txTime, _ := ctx[0].GetStub().GetTxTimestamp()
	return &Exchange{
		Base: Base{
			Id:           helper.GenerateID(doc.Exchange, ctx[0].GetStub().GetTxID()),
			CreatedAt:    helper.TimestampISO(txTime.Seconds),
			UpdatedAt:    helper.TimestampISO(txTime.Seconds),
			BlockChainId: ctx[0].GetStub().GetTxID(),
		},
	}
}