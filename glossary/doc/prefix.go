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

// Package doc contains prefixes for state database document. For example, a Token
// object will be saved in the state database with key \u0000Token....
// This helps searching for this particular object much easier.
package doc

const (
	Transactions     = "Transactions"
	Wallets          = "Wallets"
	Tokens           = "Tokens"
	HealthCheck      = "HealthCheck"
	Enrollments      = "Enrollments"
	NftToken         = "NftToken"
	Exchange         = "Exchange"
	SpotBalances     = "SpotBalances"
	IaoBalances      = "IaoBalances"
	ExchangeBalances = "ExchangeBalances"
	Iao              = "Iao"
	InvestorBook     = "InvestorBook"
	Asset            = "Asset"
	BuyIaoCache      = "BuyIaoCache"
)
