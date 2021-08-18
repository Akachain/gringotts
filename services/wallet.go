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

package services

import (
	"github.com/Akachain/gringotts/glossary"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Wallet interface {
	// Create to create new wallet. Each wallet belong to token type
	Create(ctx contractapi.TransactionContextInterface, status glossary.Status) (string, error)

	// Update to update status of wallet. Active or InActive
	Update(ctx contractapi.TransactionContextInterface, walletId string, status glossary.Status) error

	// EnrollToken to register wallet id that will be issue/mint token into.
	// Currently support add list from wallet id and to wallet id
	EnrollToken(ctx contractapi.TransactionContextInterface, tokenId string, fromWalletId []string, toWalletId []string) error

	// CreateMultipleWallet to create multiple wallet with one call
	CreateMultipleWallet(ctx contractapi.TransactionContextInterface, status glossary.Status, numberWallet int) ([]string, error)
}
