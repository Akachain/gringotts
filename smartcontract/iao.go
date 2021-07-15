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

package smartcontract

import (
	"github.com/Akachain/gringotts/dto/iao"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Iao interface {
	// CreateAsset to create new asset and token type
	CreateAsset(ctx contractapi.TransactionContextInterface, asset iao.CreateAsset) (string, error)

	// CreateIao to create new iao for asset. It will return address of Iao to investor buy asset token
	CreateIao(ctx contractapi.TransactionContextInterface, assetIao iao.AssetIao) (string, error)

	// UpdateStatusIao to update status of IAO
	UpdateStatusIao(ctx contractapi.TransactionContextInterface, updateIao iao.UpdateIao) error

	// BuyAssetToken investor call to buy asset token
	BuyAssetToken(ctx contractapi.TransactionContextInterface, asset iao.BuyBatchAsset) (string, error)
}
