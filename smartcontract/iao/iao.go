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
	"github.com/Akachain/gringotts/dto/iao"
	"github.com/Akachain/gringotts/handler"
	"github.com/Akachain/gringotts/smartcontract"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type iaoSc struct {
	iaoHandler handler.IaoHandler
}

func NewIaoSc() smartcontract.Iao {
	return &iaoSc{iaoHandler: handler.NewIaoHandler()}
}

func (i *iaoSc) CreateAsset(ctx contractapi.TransactionContextInterface, asset iao.CreateAsset) (string, error) {
	return i.iaoHandler.CreateAsset(ctx, asset)
}

func (i *iaoSc) CreateIao(ctx contractapi.TransactionContextInterface, assetIao iao.AssetIao) (string, error) {
	return i.iaoHandler.CreateIao(ctx, assetIao)
}

func (i *iaoSc) BuyAssetToken(ctx contractapi.TransactionContextInterface, batchAsset iao.BuyBatchAsset) (string, error) {
	return i.iaoHandler.BuyBatchAsset(ctx, batchAsset)
}
