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
	"github.com/Akachain/gringotts/glossary/transaction"
	"github.com/pkg/errors"
)

type BuyAsset struct {
	ReqId    string `json:"reqId"`
	IaoId    string `json:"iaoId"`
	WalletId string `json:"walletId"`
	TokenId  string `json:"tokenId"`
	NumberAT string `json:"numberAT"`
}

type ResultHandle struct {
	Status         transaction.Status `json:"status"`
	NumberATFilled string             `json:"numberATFilled"`
	NumberST       string             `json:"numberST"`
	ReqId          string             `json:"reqId"`
}

type BuyBatchAsset struct {
	Requests []BuyAsset `json:"requests"`
}

func (b BuyBatchAsset) IsValid() error {
	if len(b.Requests) <= 0 {
		return errors.New("Input invalidate")
	}
	return nil
}

func (b BuyAsset) CloneToResult() ResultHandle {
	return ResultHandle{
		ReqId: b.ReqId,
	}
}
