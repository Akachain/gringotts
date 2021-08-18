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
	"encoding/json"
	"github.com/Akachain/gringotts/glossary/doc"
	"github.com/Akachain/gringotts/helper"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type UTXO struct {
	// WalletId is wallet Id of owner
	WalletId string
	// TokenId is token Id of wallet
	TokenId string
	// Amount number token of UTXO
	Amount string
	// Status is a status of UTXO: 0 - Unspent, 1 - Spent
	Status int
	// Base of entity
	Base `mapstructure:",squash"`
}

// NewUTXO return new entity UTXO with default status is Unspent
func NewUTXO(ctx ...contractapi.TransactionContextInterface) *UTXO {
	if len(ctx) <= 0 {
		return &UTXO{}
	}
	txTime, _ := ctx[0].GetStub().GetTxTimestamp()
	return &UTXO{
		Base: Base{
			Id:           helper.GenerateID(doc.Utxo, ctx[0].GetStub().GetTxID()),
			CreatedAt:    helper.TimestampISO(txTime.Seconds),
			UpdatedAt:    helper.TimestampISO(txTime.Seconds),
			BlockChainId: ctx[0].GetStub().GetTxID(),
		},
		Status: 0,
	}
}

// ToJsonString return json string UTXO struct
func (u *UTXO) ToJsonString() (string, error) {
	data, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
