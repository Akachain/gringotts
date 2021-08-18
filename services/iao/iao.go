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
	"fmt"
	"github.com/Akachain/gringotts/entity"
	"github.com/Akachain/gringotts/errorcode"
	"github.com/Akachain/gringotts/glossary"
	"github.com/Akachain/gringotts/glossary/doc"
	statusIao "github.com/Akachain/gringotts/glossary/iao"
	"github.com/Akachain/gringotts/helper"
	"github.com/Akachain/gringotts/helper/glogger"
	"github.com/Akachain/gringotts/services"
	"github.com/Akachain/gringotts/services/base"
	"github.com/Akachain/gringotts/services/token"
	"github.com/Akachain/gringotts/services/wallet"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type iaoService struct {
	*base.Base
	tokenService  services.Token
	walletService services.Wallet
}

func NewIaoService() services.Iao {
	return &iaoService{
		base.NewBase(),
		token.NewTokenService(),
		wallet.NewWalletService(),
	}
}

func (i *iaoService) CreateAsset(ctx contractapi.TransactionContextInterface, code, name, ownerWallet, tokenName, tickerToken, maxSupply, totalValue, documentUrl string) (string, error) {
	glogger.GetInstance().Info(ctx, "-----------Iao Service - CreateAsset-----------")

	tokenId, err := i.tokenService.CreateType(ctx, tokenName, tickerToken, maxSupply)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Iao Service - Create Token type of asset failed with err (%s)", err.Error())
		return "", err
	}

	assetEntity := entity.NewAsset(ctx)
	assetEntity.Name = name
	assetEntity.Code = code
	assetEntity.OwnerWallet = ownerWallet
	assetEntity.TokenId = tokenId
	assetEntity.Documents = documentUrl
	assetEntity.TokenAmount = maxSupply
	assetEntity.TotalValue = totalValue
	assetEntity.RemainingToken = maxSupply

	if err := i.Repo.Create(ctx, assetEntity, doc.Asset, helper.AssetKey(assetEntity.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Iao Service - Create asset failed with error (%s)", err.Error())
		return "", helper.RespError(errorcode.BizUnableCreateAsset)
	}

	result := fmt.Sprintf("{\"assetId\":\"%s\",\"tokenId\":\"%s\"}", assetEntity.Id, tokenId)
	glogger.GetInstance().Infof(ctx, "-----------Iao Service - CreateAsset succeed (%s)-----------", result)

	return result, nil
}

func (i *iaoService) CreateIao(ctx contractapi.TransactionContextInterface, assetId, assetTokenAmount, startDate, endDate string) (string, error) {
	glogger.GetInstance().Info(ctx, "-----------Iao Service - CreateIao-----------")

	assetEntity, err := i.GetAsset(ctx, assetId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Iao Service - Get asset failed with err (%s)", err.Error())
		return "", err
	}

	if helper.CompareStringBalance(assetEntity.RemainingToken, assetTokenAmount) < 0 {
		glogger.GetInstance().Error(ctx, "Iao Service - Asset do not have enough token for creat Iao campaign")
		return "", helper.RespError(errorcode.BizUnableCreateIao)
	}

	// create main wallet and holder AT wallet
	wallets, err := i.walletService.CreateMultipleWallet(ctx, glossary.Active, 2)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Iao Service - Create multiple wallet failed with err (%s)", err.Error())
		return "", err
	}
	if len(wallets) != 2 {
		glogger.GetInstance().Error(ctx, "Iao Service - Create multiple wallet less than 2")
		return "", helper.RespError(errorcode.BizUnableCreateIao)
	}

	iaoEntity := entity.NewIao(ctx)
	iaoEntity.AssetId = assetId
	iaoEntity.AssetTokenId = assetEntity.TokenId
	iaoEntity.AssetTokenAmount = "0"
	iaoEntity.MainWallet = wallets[0]
	iaoEntity.HolderWallet = wallets[1]
	iaoEntity.StartDate = startDate
	iaoEntity.EndDate = endDate

	if err := i.Repo.Create(ctx, iaoEntity, doc.Iao, helper.IaoKey(iaoEntity.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Iao Service - Create iao campaign of asset failed with error (%s)", err.Error())
		return "", helper.RespError(errorcode.BizUnableCreateIao)
	}

	result := fmt.Sprintf("{\"mainWalletId\":\"%s\",\"holderWalletId\":\"%s\"}", wallets[0], wallets[1])
	glogger.GetInstance().Infof(ctx, "-----------Iao Service - CreateIao succeed (%s)-----------", result)

	return result, nil
}

func (i *iaoService) UpdateStatusIao(ctx contractapi.TransactionContextInterface, iaoId string, status statusIao.Status) error {
	txTime, _ := ctx.GetStub().GetTxTimestamp()
	iaoEntity, err := i.GetIao(ctx, iaoId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "UpdateStatus - GetIao failed with err (%s)", err.Error())
		return err
	}

	if iaoEntity.Status != status {
		iaoEntity.UpdatedAt = helper.TimestampISO(txTime.Seconds)
		iaoEntity.Status = status
		if err := i.Repo.Update(ctx, iaoEntity, doc.Iao, helper.IaoKey(iaoEntity.Id)); err != nil {
			glogger.GetInstance().Errorf(ctx, "UpdateStatus - Update Iao (%s) failed with err (%v)", err.Error())
			return helper.RespError(errorcode.BizUnableUpdateIao)
		}
	}

	return nil
}
