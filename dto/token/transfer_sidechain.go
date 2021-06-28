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
	"github.com/Akachain/gringotts/glossary/sidechain"
	"github.com/Akachain/gringotts/helper"
	"github.com/pkg/errors"
)

type TransferSideChain struct {
	WalletId  string             `json:"fromWalletId"`
	TokenId   string             `json:"tokenId"`
	FromChain sidechain.SideName `json:"fromChain"`
	ToChain   sidechain.SideName `json:"toChain"`
	Amount    string             `json:"amount"`
}

func (t TransferSideChain) IsValid() error {
	if t.WalletId == "" {
		return errors.New("Wallet Id is empty")
	}
	if t.TokenId == "" {
		return errors.New("Token id is empty")
	}

	if !t.FromChain.IsValidate() {
		return errors.New("From chain name is invalidate")
	}

	if !t.ToChain.IsValidate() {
		return errors.New("To chain name is invalidate")
	}

	if helper.CompareStringBalance(t.Amount, "0") <= 0 {
		return errors.New("Amount is zero or negative number")
	}
	return nil
}

func (t *TransferSideChain) String() string {
	return helper.MarshalStruct(t)
}
