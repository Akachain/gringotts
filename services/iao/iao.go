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
	"encoding/json"
	"fmt"
	"github.com/Akachain/gringotts/dto/iao"
	"github.com/Akachain/gringotts/entity"
	"github.com/Akachain/gringotts/errorcode"
	"github.com/Akachain/gringotts/glossary"
	"github.com/Akachain/gringotts/glossary/doc"
	"github.com/Akachain/gringotts/glossary/transaction"
	"github.com/Akachain/gringotts/helper"
	"github.com/Akachain/gringotts/helper/glogger"
	"github.com/Akachain/gringotts/services"
	"github.com/Akachain/gringotts/services/base"
	"github.com/Akachain/gringotts/services/token"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/panjf2000/ants/v2"
	"sync"
)

type iaoService struct {
	*base.Base
	tokenService services.Token
}

func NewIaoService() services.Iao {
	return &iaoService{
		base.NewBase(),
		token.NewTokenService(),
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

func (i *iaoService) CreateIao(ctx contractapi.TransactionContextInterface, assetId, assetTokenAmount, startDate, endDate string, rate int64) (string, error) {
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

	iaoEntity := entity.NewIao(ctx)
	iaoEntity.AssetId = assetId
	iaoEntity.AssetTokenId = assetEntity.TokenId
	iaoEntity.AssetTokenAmount = "0"
	iaoEntity.StableTokenAmount = "0"
	iaoEntity.StartDate = startDate
	iaoEntity.EndDate = endDate
	iaoEntity.Rate = rate

	if err := i.Repo.Create(ctx, iaoEntity, doc.Iao, helper.IaoKey(iaoEntity.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Iao Service - Create iao campaign of asset failed with error (%s)", err.Error())
		return "", helper.RespError(errorcode.BizUnableCreateIao)
	}

	// create tx transfer AT
	txEntity := entity.NewTransaction(ctx)
	txEntity.SpenderWallet = iaoEntity.Id
	txEntity.FromWallet = assetEntity.OwnerWallet
	txEntity.ToWallet = iaoEntity.Id
	txEntity.FromTokenId = assetEntity.TokenId
	txEntity.ToTokenId = assetEntity.TokenId
	txEntity.FromTokenAmount = assetTokenAmount
	txEntity.ToTokenAmount = assetTokenAmount
	txEntity.TxType = transaction.IaoDepositAT
	txEntity.Note = iaoEntity.Id

	if err := i.Repo.Create(ctx, txEntity, doc.Transactions, helper.TransactionKey(txEntity.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Iao Service - Create transaction failed with error (%s)", err.Error())
		return "", helper.RespError(errorcode.BizUnableCreateTX)
	}

	glogger.GetInstance().Infof(ctx, "-----------Iao Service - CreateIao succeed (%s)-----------", iaoEntity.Id)

	return iaoEntity.Id, nil
}

func (i *iaoService) BuyBatchAsset(ctx contractapi.TransactionContextInterface, batchReq []iao.BuyAsset) (string, error) {
	glogger.GetInstance().Info(ctx, "Start BuyBatchAsset")
	var wg sync.WaitGroup
	// calculate hash and check cache
	stringInput, _ := json.Marshal(batchReq)
	inputHash := helper.CalculateHash(string(stringInput))
	cacheEntity, isExisted, err := i.GetBuyIaoCache(ctx, inputHash)
	if err != nil {
		return "", err
	}

	if isExisted {
		return cacheEntity.Result, nil
	}

	// cache iao
	iaoMap := make(map[string]*entity.Iao, 0)
	balanceMap := make(map[string]*entity.BalanceCache, len(batchReq))
	resultHandle := make([]iao.ResultHandle, 0, len(batchReq))
	investorMap := make(map[string]*entity.InvestorBook, len(batchReq))

	for _, req := range batchReq {
		res := req.CloneToResult()
		iaoEntity, err := i.getIaoInfo(ctx, iaoMap, req.IaoId)
		if err != nil {
			glogger.GetInstance().Errorf(ctx, "BuyBatchAsset - Handle req (%s) failed get IAO", req.ReqId, err.Error())
			res.Status = transaction.Rejected
			resultHandle = append(resultHandle, res)
			continue
		}
		if helper.CompareStringBalance(iaoEntity.RemainingAssetToken, "0") <= 0 {
			glogger.GetInstance().Errorf(ctx, "BuyBatchAsset - Handle req (%s) remaining of iao is zero", req.ReqId)
			res.Status = transaction.Rejected
			resultHandle = append(resultHandle, res)
			// TODO: return not continue for optimize
			continue
		}

		var numberATBuy string
		if iaoEntity.RemainingAssetToken >= req.NumberAT {
			numberATBuy = req.NumberAT
		} else {
			numberATBuy = iaoEntity.RemainingAssetToken
		}

		stableToken, err := helper.MulBalance(numberATBuy, iaoEntity.Rate)
		if err != nil {
			glogger.GetInstance().Errorf(ctx, "BuyBatchAsset - Handle req (%s) failed to calculate stable token", req.ReqId)
			res.Status = transaction.Rejected
			resultHandle = append(resultHandle, res)
			continue
		}

		err = i.SubAmount(ctx, balanceMap, doc.IaoBalances, req.WalletId, req.TokenId, stableToken)
		if err != nil {
			glogger.GetInstance().Errorf(ctx, "BuyBatchAsset - Handle req (%s) failed get Balance", req.ReqId, err.Error())
			res.Status = transaction.Rejected
			resultHandle = append(resultHandle, res)
			continue
		}

		if err := i.addInvestorBook(ctx, investorMap, req, stableToken); err != nil {
			glogger.GetInstance().Errorf(ctx, "BuyBatchAsset - Handle req (%s) failed get create investor book", req.ReqId, err.Error())
			res.Status = transaction.Rejected
			resultHandle = append(resultHandle, res)
			continue
		}
		updateATRemain, _ := helper.SubBalance(iaoEntity.RemainingAssetToken, numberATBuy)
		updateST, _ := helper.AddBalance(iaoEntity.StableTokenAmount, stableToken)
		iaoEntity.RemainingAssetToken = updateATRemain
		iaoEntity.StableTokenAmount = updateST
		iaoMap[iaoEntity.Id] = iaoEntity

		res.Status = transaction.Confirmed
		res.NumberATFilled = numberATBuy
		resultHandle = append(resultHandle, res)
	}

	// log create Read/Write Set
	glogger.GetInstance().Info(ctx, "Start Create Read/Write")
	respErrors := make(chan error)
	wgDone := make(chan bool)

	resultJson, _ := json.Marshal(resultHandle)

	// update iao
	if err := i.updateIao(ctx, iaoMap); err != nil {
		return "", err
	}

	// update balance
	poolBalance, err := ants.NewPoolWithFunc(glossary.NumberWorker, func(input interface{}) {
		err := i.AsyncUpdateBalance(ctx, input)
		if err != nil {
			glogger.GetInstance().Errorf(ctx, "BuyBatchAsset - Update balance failed with err (%s)", err.Error())
			respErrors <- err
		}
		wg.Done()
	})
	defer poolBalance.Release()
	for _, balanceCache := range balanceMap {
		wg.Add(1)
		_ = poolBalance.Invoke(balanceCache)
	}

	// insert investor book
	pooInvestorBook, err := ants.NewPoolWithFunc(glossary.NumberWorker, func(input interface{}) {
		err := i.insertInvestorBook(ctx, input)
		if err != nil {
			glogger.GetInstance().Errorf(ctx, "BuyBatchAsset - Insert investor book failed with err (%s)", err.Error())
			respErrors <- err
		}
		wg.Done()
	})
	defer pooInvestorBook.Release()
	for _, investorBook := range investorMap {
		wg.Add(1)
		_ = pooInvestorBook.Invoke(investorBook)
	}
	go func() {
		wg.Wait()
		close(wgDone)
	}()

	select {
	case <-wgDone:
		buyCache := entity.NewIaoCache(ctx)
		buyCache.Hash = inputHash
		buyCache.Result = string(resultJson)
		if err := i.Repo.Create(ctx, buyCache, doc.BuyIaoCache, helper.ResultCacheKey(buyCache.Hash)); err != nil {
			glogger.GetInstance().Errorf(ctx, "BuyBatchAsset - Create cache buy Iao failed with err (%v)", err.Error())
			return "", helper.RespError(errorcode.BizUnableUpdateIao)
		}
	case err := <-respErrors:
		close(respErrors)
		glogger.GetInstance().Errorf(ctx, "BuyBatchAsset - Handle buy Iao failed with err (%v)", err.Error())
		return "", err
	}
	glogger.GetInstance().Info(ctx, "End BuyBatchAsset")

	return string(resultJson), nil
}

func (i *iaoService) getIaoInfo(ctx contractapi.TransactionContextInterface, iaoMap map[string]*entity.Iao, iaoId string) (*entity.Iao, error) {
	if _, ok := iaoMap[iaoId]; !ok {
		iaoEntity, err := i.GetIao(ctx, iaoId)
		if err != nil {
			return nil, err
		}

		iaoMap[iaoEntity.Id] = iaoEntity
	}

	return iaoMap[iaoId], nil
}

func (i *iaoService) getBalanceWalletInfo(ctx contractapi.TransactionContextInterface, balanceMap map[string]*entity.Balance,
	walletId, tokenId string) (*entity.Balance, error) {
	if _, ok := balanceMap[walletId]; !ok {
		balanceEntity, err := i.GetBalanceOfToken(ctx, doc.IaoBalances, walletId, tokenId)
		if err != nil {
			return nil, err
		}
		balanceMap[walletId] = balanceEntity
	}

	return balanceMap[walletId], nil
}

// updateIao update remaining token of IAO after handle transaction
func (i *iaoService) updateIao(ctx contractapi.TransactionContextInterface, iaoMap map[string]*entity.Iao) error {
	txTime, _ := ctx.GetStub().GetTxTimestamp()
	for _, itemIao := range iaoMap {
		itemIao.UpdatedAt = helper.TimestampISO(txTime.Seconds)
		if err := i.Repo.Update(ctx, itemIao, doc.Iao, helper.IaoKey(itemIao.Id)); err != nil {
			glogger.GetInstance().Errorf(ctx, "UpdateIao - Update Iao (%s) failed with err (%v)", itemIao.Id, err.Error())
			return helper.RespError(errorcode.BizUnableUpdateIao)
		}
		glogger.GetInstance().Debugf(ctx, "UpdateIao - Update Iao (%s) succeed", itemIao.Id)
	}

	return nil
}

func (i *iaoService) insertInvestorBook(ctx contractapi.TransactionContextInterface, input interface{}) error {
	investorBook := input.(*entity.InvestorBook)
	if err := i.Repo.Update(ctx, investorBook, doc.InvestorBook, helper.InvestorBookKey(investorBook.WalletId, investorBook.IaoId)); err != nil {
		glogger.GetInstance().Errorf(ctx, "UpdateIao - Update Iao (%s) failed with err (%v)", investorBook.Id, err.Error())
		return helper.RespError(errorcode.BizUnableCreateInvestorBook)
	}
	return nil
}

func (i *iaoService) addInvestorBook(ctx contractapi.TransactionContextInterface,
	investorBookMap map[string]*entity.InvestorBook, req iao.BuyAsset, amountST string) (err error) {
	var investorBookEntity *entity.InvestorBook
	key := req.IaoId + "_" + req.WalletId
	if _, ok := investorBookMap[key]; !ok {
		investorBookEntity, err = i.GetInvestorBook(ctx, req.WalletId, req.IaoId)
		if err != nil {
			return err
		}

		if investorBookEntity == nil {
			investorBookId := helper.GenerateID(doc.InvestorBook, req.ReqId)
			investorBookEntity = entity.NewInvestorBook(ctx)
			investorBookEntity.IaoId = req.IaoId
			investorBookEntity.WalletId = req.WalletId
			investorBookEntity.Id = investorBookId
		}
	} else {
		investorBookEntity = investorBookMap[key]
	}

	balanceAT, err := helper.AddBalance(investorBookEntity.AssetTokenAmount, req.NumberAT)
	if err != nil {
		return err
	}

	balanceST, err := helper.AddBalance(investorBookEntity.StableTokenAmount, amountST)
	if err != nil {
		return err
	}
	investorBookEntity.AssetTokenAmount = balanceAT
	investorBookEntity.StableTokenAmount = balanceST

	investorBookMap[key] = investorBookEntity

	return nil
}
