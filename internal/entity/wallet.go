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
	"github.com/Akachain/gringotts/glossary"
	"github.com/Akachain/gringotts/glossary/doc"
	"github.com/Akachain/gringotts/helper"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// A wallet only contains 1 type of token and its balance.
type Wallet struct {
	TokenId  string          `mapstructure:"tokenId" json:"tokenId"`
	Status   glossary.Status `mapstructure:"status" json:"status"`
	Balances string          `mapstructure:"balances" json:"balances"`
	Base     `mapstructure:",squash"`
}

func NewWallet(ctx ...contractapi.TransactionContextInterface) *Wallet {
	if len(ctx) <= 0 {
		return &Wallet{}
	}
	txTime, _ := ctx[0].GetStub().GetTxTimestamp()
	return &Wallet{
		Base: Base{
			Id:           helper.GenerateID(doc.Wallets, ctx[0].GetStub().GetTxID()),
			CreatedAt:    helper.TimestampISO(txTime.Seconds),
			UpdatedAt:    helper.TimestampISO(txTime.Seconds),
			BlockChainId: ctx[0].GetStub().GetTxID(),
		},
	}
}
