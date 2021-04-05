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
	"github.com/Akachain/gringotts/dto"
	"github.com/Akachain/gringotts/errorcode"
	"github.com/Akachain/gringotts/helper"
	"github.com/Akachain/gringotts/helper/glogger"
	"github.com/Akachain/gringotts/services"
	"github.com/Akachain/gringotts/services/token"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type TokenHandler struct {
	tokenService services.Token
}

func NewTokenHandler() *TokenHandler {
	tokenService := token.NewTokenService()
	return &TokenHandler{tokenService: tokenService}
}

// Transfer to transfer token between wallet.
func (t *TokenHandler) Transfer(ctx contractapi.TransactionContextInterface, transferDto dto.TransferToken) error {
	glogger.GetInstance().Info(ctx, "-----------Token Handler - Transfer-----------")

	// checking dto validate
	if err := transferDto.IsValid(); err != nil {
		glogger.GetInstance().Errorf(ctx, "TokenHandler - Transfer Input invalidate %v", err)
		return helper.RespError(errorcode.InvalidParam)
	}

	return t.tokenService.Transfer(ctx, transferDto.FromWalletId, transferDto.ToWalletId, transferDto.Amount)
}

// Mint generate new token for wallet.
func (t *TokenHandler) Mint(ctx contractapi.TransactionContextInterface, mintDto dto.MintToken) error {
	glogger.GetInstance().Info(ctx, "-----------Token Handler - Mint-----------")

	// checking dto validate
	if err := mintDto.IsValid(); err != nil {
		glogger.GetInstance().Errorf(ctx, "TokenHandler - Mint Input invalidate %v", err)
		return helper.RespError(errorcode.InvalidParam)
	}

	return t.tokenService.Mint(ctx, mintDto.WalletId, mintDto.Amount)
}

// Burn to burn token existed in the system.
func (t *TokenHandler) Burn(ctx contractapi.TransactionContextInterface, burnDto dto.BurnToken) error {
	glogger.GetInstance().Info(ctx, "-----------Token Handler - Burn-----------")

	// checking dto validate
	if err := burnDto.IsValid(); err != nil {
		glogger.GetInstance().Errorf(ctx, "TokenHandler - Burn Input invalidate %v", err)
		return helper.RespError(errorcode.InvalidParam)
	}

	return t.tokenService.Burn(ctx, burnDto.WalletId, burnDto.Amount)
}

// CreateTokenType to create new token type.
func (t *TokenHandler) CreateTokenType(ctx contractapi.TransactionContextInterface, tokenTypeDto dto.CreateTokenType) (string, error) {
	glogger.GetInstance().Info(ctx, "-----------Token Handler - CreateTokenType-----------")

	// checking dto validate
	if err := tokenTypeDto.IsValid(); err != nil {
		glogger.GetInstance().Errorf(ctx, "TokenHandler - Create Token Type Input invalidate %v", err)
		return "", helper.RespError(errorcode.InvalidParam)
	}

	return t.tokenService.CreateType(ctx, tokenTypeDto.Name, tokenTypeDto.Rate)
}

// Swap to swap between different token type.
func (t *TokenHandler) Swap(ctx contractapi.TransactionContextInterface, swapDto dto.SwapToken) error {
	glogger.GetInstance().Info(ctx, "-----------Token Handler - Swap-----------")

	// checking dto validate
	if err := swapDto.IsValid(); err != nil {
		glogger.GetInstance().Errorf(ctx, "TokenHandler - Swap Input invalidate %v", err)
		return helper.RespError(errorcode.InvalidParam)
	}

	return t.tokenService.Swap(ctx, swapDto.FromWalletId, swapDto.ToWalletId, swapDto.Amount)
}

// Issue to issue new token type form stable token.
func (t *TokenHandler) Issue(ctx contractapi.TransactionContextInterface, issueDto dto.IssueToken) error {
	glogger.GetInstance().Info(ctx, "-----------Token Handler - Swap-----------")

	// checking dto validate
	if err := issueDto.IsValid(); err != nil {
		glogger.GetInstance().Errorf(ctx, "TokenHandler - Issue Input invalidate %v", err)
		return helper.RespError(errorcode.InvalidParam)
	}

	return t.tokenService.Issue(ctx, issueDto.TokenId, issueDto.FromWalletId, issueDto.ToWalletId, issueDto.Amount)
}
