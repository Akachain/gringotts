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
	// wallet use to issue new token
	WalletId string `json:"walletId"`

	// token will use to issue to other token. example is stable token
	FromTokenId string `json:"fromTokenId"`

	// token will be issued base on stable token
	ToTokenId string `json:"toTokenId"`

	// number of token will be use to issue new token. Use base unit (ax10^8)
	FromTokenAmount string `json:"fromTokenAmount"`

	// number of new token will be issue. Use base unit (ax10^8)
	ToTokenAmount string `json:"toTokenAmount"`
}

func (i IssueToken) IsValid() error {
	if i.WalletId == "" {
		return errors.New("Wallet Id is empty")
	}
	if i.FromTokenId == "" || i.ToTokenId == "" {
		return errors.New("From/To token id is empty")
	}

	if i.FromTokenAmount == "" {
		return errors.New("the amount of from token is empty")
	}

	if i.ToTokenAmount == "" {
		return errors.New("the amount of new token is empty")
	}
	return nil
}
