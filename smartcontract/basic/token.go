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

package basic

import (
	"github.com/Akachain/gringotts/dto"
	"github.com/Akachain/gringotts/handler"
	"github.com/Akachain/gringotts/helper/glogger"
	"github.com/Akachain/gringotts/pkg/ipc"
	"github.com/Akachain/gringotts/smartcontract"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type baseToken struct {
	contractapi.Contract
	tokenHandler       *handler.TokenHandler
	walletHandler      *handler.WalletHandler
	healthCheckHandler handler.HealthCheckHandler
	accountingHandler  handler.AccountingHandler
}

func NewBaseToken(nftIpc ...ipc.Ipc) smartcontract.BasicToken {
	return &baseToken{
		tokenHandler:       handler.NewTokenHandler(),
		walletHandler:      handler.NewWalletHandler(),
		healthCheckHandler: handler.NewHealthCheckHandler(),
		accountingHandler:  handler.NewAccountingHandler(nftIpc...),
	}
}

func (b *baseToken) InitLedger(ctx contractapi.TransactionContextInterface) error {
	glogger.GetInstance().Info(ctx, "------------Init ChainCode------------")
	return nil
}

// Wallet feature
func (b *baseToken) CreateWallet(ctx contractapi.TransactionContextInterface, createWallet dto.CreateWallet) (string, error) {
	glogger.GetInstance().Info(ctx, "------------CreateWallet ChainCode------------")
	return b.walletHandler.CreateWallet(ctx, createWallet)
}

func (b *baseToken) UpdateWallet(ctx contractapi.TransactionContextInterface, updateWallet dto.UpdateWallet) error {
	glogger.GetInstance().Info(ctx, "------------UpdateWallet ChainCode------------")
	return b.walletHandler.UpdateWallet(ctx, updateWallet)
}

func (b *baseToken) GetBalance(ctx contractapi.TransactionContextInterface, balance dto.Balance) (string, error) {
	glogger.GetInstance().Info(ctx, "------------GetBalance ChainCode------------")
	return b.walletHandler.BalanceOf(ctx, balance)
}

// Token feature
func (b *baseToken) Mint(ctx contractapi.TransactionContextInterface, mintDto dto.MintToken) error {
	glogger.GetInstance().Info(ctx, "------------Mint ChainCode------------")
	return b.tokenHandler.Mint(ctx, mintDto)
}

func (b *baseToken) Burn(ctx contractapi.TransactionContextInterface, burnDto dto.BurnToken) error {
	glogger.GetInstance().Info(ctx, "------------Burn ChainCode------------")
	return b.tokenHandler.Burn(ctx, burnDto)
}

func (b *baseToken) Transfer(ctx contractapi.TransactionContextInterface, transferDto dto.TransferToken) error {
	glogger.GetInstance().Info(ctx, "------------Transfer ChainCode------------")
	return b.tokenHandler.Transfer(ctx, transferDto)
}

func (b *baseToken) CreateTokenType(ctx contractapi.TransactionContextInterface, createTokenTypeDto dto.CreateTokenType) (string, error) {
	glogger.GetInstance().Info(ctx, "------------CreateTokenType ChainCode------------")
	return b.tokenHandler.CreateTokenType(ctx, createTokenTypeDto)
}

func (b *baseToken) CreateHealthCheck(ctx contractapi.TransactionContextInterface, arg string) (string, error) {
	glogger.GetInstance().Info(ctx, "------------CreateHealthCheck ChainCode------------")
	return b.healthCheckHandler.CreateHealthCheck(ctx)
}

// Accounting feature
func (b *baseToken) GetAccountingTx(ctx contractapi.TransactionContextInterface) ([]string, error) {
	glogger.GetInstance().Info(ctx, "------------GetAccountingTx ChainCode------------")
	return b.accountingHandler.GetAccountingTx(ctx)
}

func (b *baseToken) CalculateBalance(ctx contractapi.TransactionContextInterface, accountingDto dto.AccountingBalance) error {
	glogger.GetInstance().Info(ctx, "------------CalculateBalance ChainCode------------")
	return b.accountingHandler.CalculateBalance(ctx, accountingDto)
}

func (b *baseToken) Swap(ctx contractapi.TransactionContextInterface, swapDto dto.SwapToken) error {
	glogger.GetInstance().Info(ctx, "------------Swap ChainCode------------")
	return b.tokenHandler.Swap(ctx, swapDto)
}

func (b *baseToken) Issue(ctx contractapi.TransactionContextInterface, issueDto dto.IssueToken) error {
	glogger.GetInstance().Info(ctx, "------------Issue ChainCode------------")
	return b.tokenHandler.Issue(ctx, issueDto)
}

func (b *baseToken) EnrollToken(ctx contractapi.TransactionContextInterface, enrollmentDto dto.Enrollment) error {
	glogger.GetInstance().Info(ctx, "------------EnrollToken ChainCode------------")
	return b.walletHandler.EnrollToken(ctx, enrollmentDto)
}
