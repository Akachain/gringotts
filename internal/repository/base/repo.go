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

// The basic repository implementation use utility functions in akc-go-sdk
package base

import (
	"github.com/Akachain/akc-go-sdk/util"
	"github.com/Akachain/gringotts/glossary"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
)

type repo struct {
}

func NewRepository() *repo {
	return &repo{}
}

func (r *repo) GetQueryStringWithPagination(ctx contractapi.TransactionContextInterface, queryString string) (shim.StateQueryIteratorInterface, error) {
	res, _, err := ctx.GetStub().GetQueryResultWithPagination(queryString, glossary.PaginationSize, "")
	return res, err
}

func (r *repo) Create(ctx contractapi.TransactionContextInterface, entity interface{}, tableModel string, keys []string) error {
	return util.Createdata(ctx.GetStub(), tableModel, keys, entity)
}

func (r *repo) Update(ctx contractapi.TransactionContextInterface, entity interface{}, tableModel string, keys []string) error {
	return util.UpdateExistingData(ctx.GetStub(), tableModel, keys, entity)
}

func (r *repo) Get(ctx contractapi.TransactionContextInterface, tableModel string, keys []string) (interface{}, error) {
	return util.Getdatabyrowkeys(ctx.GetStub(), keys, tableModel)
}

func (r *repo) IsExist(ctx contractapi.TransactionContextInterface, docPrefix string, keys []string) (bool, error) {
	stub := ctx.GetStub()
	compositeKey, err := stub.CreateCompositeKey(docPrefix, keys)
	if err != nil {
		return false, errors.WithMessage(err, "IsExist - Create composite key fail")
	}

	var bytes []byte
	bytes, err = stub.GetState(compositeKey)
	if err != nil {
		return false, errors.WithMessage(err, "IsExist - Get document fail")
	}

	return bytes != nil, nil
}
