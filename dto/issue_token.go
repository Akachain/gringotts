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

import "github.com/pkg/errors"

type IssueToken struct {
	TokenId      string  `json:"tokenId"`
	FromWalletId string  `json:"fromWalletId"`
	ToWalletId   string  `json:"toWalletId"`
	Amount       float64 `json:"amount"`
}

func (i IssueToken) IsValid() error {
	if i.TokenId == "" {
		return errors.New("Token Id is empty")
	}
	if i.FromWalletId == "" || i.ToWalletId == "" {
		return errors.New("From/To wallet id is empty")
	}

	if i.Amount <= 0 {
		return errors.New("the issue amount is a negative/zero number")
	}
	return nil
}
