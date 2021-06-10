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

package handler

import (
	"github.com/Akachain/gringotts/dto/token"
	"github.com/Akachain/gringotts/errorcode"
	"github.com/Akachain/gringotts/helper"
	"github.com/Akachain/gringotts/helper/glogger"
	"github.com/Akachain/gringotts/services"
	"github.com/Akachain/gringotts/services/wallet"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type WalletHandler struct {
	walletService services.Wallet
}

func NewWalletHandler() *WalletHandler {
	return &WalletHandler{
		walletService: wallet.NewWalletService(),
	}
}

// CreateWallet generate new wallet with token type
func (w *WalletHandler) CreateWallet(ctx contractapi.TransactionContextInterface, createWalletDto token.CreateWallet) (string, error) {
	glogger.GetInstance().Info(ctx, "-----------Wallet Handler - CreateWallet-----------")

	// checking dto validate
	if err := createWalletDto.IsValid(); err != nil {
		glogger.GetInstance().Errorf(ctx, "Wallet Handler - Create Wallet Input invalidate %v", err)
		return "", helper.RespError(errorcode.InvalidParam)
	}
	return w.walletService.Create(ctx, createWalletDto.TokenId, createWalletDto.Status)
}

// UpdateWallet to update status of wallet
func (w *WalletHandler) UpdateWallet(ctx contractapi.TransactionContextInterface, updateWalletDto token.UpdateWallet) error {
	glogger.GetInstance().Info(ctx, "-----------Wallet Handler - UpdateWallet-----------")

	// checking dto validate
	if err := updateWalletDto.IsValid(); err != nil {
		glogger.GetInstance().Errorf(ctx, "Wallet Handler - Update Wallet Input invalidate %v", err)
		return helper.RespError(errorcode.InvalidParam)
	}
	return w.walletService.Update(ctx, updateWalletDto.WalletId, updateWalletDto.Status)
}

// BalanceOf to return balance of wallet
func (w *WalletHandler) BalanceOf(ctx contractapi.TransactionContextInterface, balanceDto token.Balance) (string, error) {
	glogger.GetInstance().Info(ctx, "-----------Wallet Handler - BalanceOf-----------")

	// checking dto validate
	if err := balanceDto.IsValid(); err != nil {
		glogger.GetInstance().Errorf(ctx, "Wallet Handler - Balance Input invalidate %v", err)
		return "-1", helper.RespError(errorcode.InvalidParam)
	}
	return w.walletService.BalanceOf(ctx, balanceDto.WalletId, balanceDto.TokenId)
}

// EnrollToken to create or update enrollment policy for token
func (w *WalletHandler) EnrollToken(ctx contractapi.TransactionContextInterface, enrollmentDto token.Enrollment) error {
	glogger.GetInstance().Info(ctx, "-----------Wallet Handler - EnrollToken-----------")

	// checking dto validate
	if err := enrollmentDto.IsValid(); err != nil {
		glogger.GetInstance().Errorf(ctx, "Wallet Handler - Enrollment Input invalidate %v", err)
		return helper.RespError(errorcode.InvalidParam)
	}
	return w.walletService.EnrollToken(ctx, enrollmentDto.TokenId, enrollmentDto.FromWalletId, enrollmentDto.ToWalletId)
}
