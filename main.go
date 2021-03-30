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

package main

import (
	"github.com/Akachain/gringotts/helper/glogger"
	"github.com/Akachain/gringotts/smartcontract/basic"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	glogger.GetInstance().Logger.Info("Starting chaincode")
	sc := basic.NewBaseToken()
	chaincode, err := contractapi.NewChaincode(sc)

	if err != nil {
		glogger.GetInstance().Logger.Errorf("Error create akc chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		glogger.GetInstance().Logger.Errorf("Error starting chaincode: %s", err.Error())
	}
}
