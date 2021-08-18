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
	"github.com/pkg/errors"
)

type ItemPair struct {
	TokenId string    `json:"tokenId"`
	Inputs  []string  `json:"inputs"`
	Outputs []UTXODto `json:"outputs"`
}

type ExchangeToken struct {
	Pairs    []ItemPair `json:"pairs"`
	Metadata string     `json:"metadata"`
}

func (ex ExchangeToken) IsValid() error {
	if ex.Pairs == nil || len(ex.Pairs) <= 0 {
		errors.New("List pair exchange empty")
	}

	for _, pair := range ex.Pairs {
		if pair.Inputs == nil || len(pair.Inputs) <= 0 {
			return errors.New("Inputs UTXO is nil or empty")
		}

		if pair.Outputs == nil || len(pair.Outputs) <= 0 {
			return errors.New("Outputs UTXO is nil or empty")
		}
	}

	return nil
}
