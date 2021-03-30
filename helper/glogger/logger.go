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

// Package Gringotts Logger (glogger) contains a singleton logging instance
// to use throughout the project.
package glogger

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric/common/flogging"
	"sync"
)

var lock = &sync.Mutex{}

type gLogger struct {
	Logger *flogging.FabricLogger
}

var instance *gLogger = nil

func newLogger() *gLogger {
	logging, _ := flogging.New(flogging.Config{})
	lg := new(gLogger)
	lg.Logger = logging.Logger("GringottsSmartContract")

	return lg
}

// Return
func GetInstance() *gLogger {
	lock.Lock()
	defer lock.Unlock()
	if instance == nil {
		instance = newLogger()
	}

	return instance
}

func (tl *gLogger) Error(ctx contractapi.TransactionContextInterface, arg ...interface{}) {
	tl.Logger.Errorf("TxBlockchain (%s) - Message: (%v)", ctx.GetStub().GetTxID(), arg)
}

func (tl *gLogger) Errorf(ctx contractapi.TransactionContextInterface, template string, arg ...interface{}) {
	reNewTemplate := "TxBlockchain (%s) - " + template
	tl.Logger.Errorf(reNewTemplate, ctx.GetStub().GetTxID(), arg)
}

func (tl *gLogger) Info(ctx contractapi.TransactionContextInterface, arg ...interface{}) {
	tl.Logger.Infof("TxBlockchain (%s) - Message: (%v)", ctx.GetStub().GetTxID(), arg)
}

func (tl *gLogger) Infof(ctx contractapi.TransactionContextInterface, template string, arg ...interface{}) {
	reNewTemplate := "TxBlockchain (%s) - " + template
	tl.Logger.Infof(reNewTemplate, ctx.GetStub().GetTxID(), arg)
}

func (tl *gLogger) Debug(ctx contractapi.TransactionContextInterface, arg ...interface{}) {
	tl.Logger.Debugf("TxBlockchain (%s) - Message: (%v)", ctx.GetStub().GetTxID(), arg)
}

func (tl *gLogger) Debugf(ctx contractapi.TransactionContextInterface, template string, arg ...interface{}) {
	reNewTemplate := "TxBlockchain (%s) - " + template
	tl.Logger.Debugf(reNewTemplate, ctx.GetStub().GetTxID(), arg)
}
