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

package exchange

import (
	"github.com/Akachain/gringotts/entity"
	"github.com/Akachain/gringotts/errorcode"
	"github.com/Akachain/gringotts/glossary/doc"
	"github.com/Akachain/gringotts/glossary/transaction"
	"github.com/Akachain/gringotts/helper"
	"github.com/Akachain/gringotts/helper/glogger"
	"github.com/Akachain/gringotts/services"
	"github.com/Akachain/gringotts/services/base"
	"github.com/Akachain/gringotts/services/token"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type exchangeService struct {
	tokenService services.Token
	*base.Base
}

func NewExchangeService() services.Exchange {
	return &exchangeService{
		token.NewTokenService(),
		base.NewBase(),
	}
}

func (e *exchangeService) TransferNft(ctx contractapi.TransactionContextInterface, ownerWalletId string, toWalletId string, nftToken string, price float64) error {
	glogger.GetInstance().Info(ctx, "-----------Exchange Service - TransferNft-----------")

	// transfer token using buy nft token
	txId, err := e.tokenService.Transfer(ctx, toWalletId, ownerWalletId, price)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Exchange Service - Call token transfer failed with err(%v)", err)
		return err
	}
	// transfer nft token
	exchangeEntity := entity.NewExchange(ctx)
	exchangeEntity.OwnerWalletId = ownerWalletId
	exchangeEntity.ToWalletId = toWalletId
	exchangeEntity.Price = price
	exchangeEntity.NftTokenId = nftToken
	exchangeEntity.Status = transaction.Pending

	if err := e.Repo.Create(ctx, exchangeEntity, doc.Exchange, helper.ExchangeKey(txId)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Exchange Service - Mint NFT failed with error (%v)", err)
		return helper.RespError(errorcode.BizUnableCreateExchange)
	}
	glogger.GetInstance().Infof(ctx, "-----------Exchange Service - TransferNft succeed (%s)-----------", exchangeEntity.Id)

	return nil
}
