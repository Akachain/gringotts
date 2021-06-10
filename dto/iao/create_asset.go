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

import "github.com/pkg/errors"

type CreateAsset struct {
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	TokenName   string `json:"tokenName"`
	TickerToken string `json:"tickerToken"`
	MaxSupply   string `json:"maxSupply"`
	TotalValue  string `json:"totalValue"`
	ExpireDate  string `json:"expireDate"`
}

func (c CreateAsset) IsValid() error {
	if c.Name == "" || c.Owner == "" {
		return errors.New("Name/Owner is empty")
	}
	if c.TokenName == "" || c.TickerToken == "" {
		return errors.New("TokenName/TickerToken id is empty")
	}

	if c.MaxSupply == "" {
		return errors.New("The max supply of asset token is empty")
	}

	if c.TotalValue == "" {
		return errors.New("Total value of new token is empty")
	}
	return nil
}
