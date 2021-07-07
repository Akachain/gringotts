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

package base

import (
	"github.com/Akachain/gringotts/entity"
	"github.com/Akachain/gringotts/errorcode"
	"github.com/Akachain/gringotts/glossary"
	"github.com/Akachain/gringotts/glossary/doc"
	"github.com/Akachain/gringotts/glossary/sidechain"
	"github.com/Akachain/gringotts/glossary/transaction"
	"github.com/Akachain/gringotts/helper"
	"github.com/Akachain/gringotts/helper/glogger"
	"github.com/Akachain/gringotts/repository"
	"github.com/Akachain/gringotts/repository/base"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type Base struct {
	Repo repository.Repo
}

func NewBase() *Base {
	return &Base{
		Repo: base.NewRepository(),
	}
}
func (b *Base) GetInvestorBook(ctx contractapi.TransactionContextInterface, investorId string) (*entity.InvestorBook, error) {
	isExisted, investorBookData, err := b.Repo.GetAndCheckExist(ctx, doc.InvestorBook, helper.InvestorBookKey(investorId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Get Investor Book (%s) failed with error (%s)", investorId, err.Error())
		return nil, helper.RespError(errorcode.BizUnableGetInvestorBook)
	}

	if !isExisted {
		return nil, nil
	}

	investorBookEntity := entity.NewInvestorBook()
	if err = mapstructure.Decode(investorBookData, &investorBookEntity); err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Decode Investor Book failed with error  (%s)", err.Error())
		return nil, helper.RespError(errorcode.BizUnableMapDecode)
	}
	return investorBookEntity, nil
}

func (b *Base) GetBuyIaoCache(ctx contractapi.TransactionContextInterface, buyIaoId string) (*entity.BuyIaoCache, bool, error) {
	isExisted, buyIaoData, err := b.Repo.GetAndCheckExist(ctx, doc.BuyIaoCache, helper.ResultCacheKey(buyIaoId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Get Buy Cache IAO (%s) failed with error (%s)", buyIaoId, err.Error())
		return nil, isExisted, helper.RespError(errorcode.BizUnableGetBuyCache)
	}

	buyIaoEntity := entity.NewIaoCache()
	if err = mapstructure.Decode(buyIaoData, &buyIaoEntity); err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Decode IAO failed with error  (%s)", err.Error())
		return nil, isExisted, helper.RespError(errorcode.BizUnableMapDecode)
	}
	return buyIaoEntity, isExisted, nil
}

func (b *Base) GetIao(ctx contractapi.TransactionContextInterface, iaoId string) (*entity.Iao, error) {
	iaoData, err := b.Repo.Get(ctx, doc.Iao, helper.IaoKey(iaoId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Get IAO (%s) failed with error (%s)", iaoId, err.Error())
		return nil, helper.RespError(errorcode.BizUnableGetIao)
	}

	iaoEntity := entity.NewIao()
	if err = mapstructure.Decode(iaoData, &iaoEntity); err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Decode IAO failed with error  (%s)", err.Error())
		return nil, helper.RespError(errorcode.BizUnableMapDecode)
	}
	return iaoEntity, nil
}

func (b *Base) GetAsset(ctx contractapi.TransactionContextInterface, assetId string) (*entity.Asset, error) {
	assetData, err := b.Repo.Get(ctx, doc.Asset, helper.AssetKey(assetId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Get asset (%s) failed with error (%s)", assetId, err.Error())
		return nil, helper.RespError(errorcode.BizUnableGetAsset)
	}

	asset := entity.NewAsset()
	if err = mapstructure.Decode(assetData, &asset); err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Decode asset failed with error  (%s)", err.Error())
		return nil, helper.RespError(errorcode.BizUnableMapDecode)
	}
	return asset, nil
}

func (b *Base) GetTransaction(ctx contractapi.TransactionContextInterface, txId string) (*entity.Transaction, error) {
	txData, err := b.Repo.Get(ctx, doc.Transactions, helper.TransactionKey(txId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Get transaction (%s) failed with error (%s)", txId, err.Error())
		return nil, helper.RespError(errorcode.BizUnableGetTx)
	}

	tx := entity.NewTransaction()
	if err = mapstructure.Decode(txData, &tx); err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Decode transaction failed with error  (%s)", err.Error())
		return nil, helper.RespError(errorcode.BizUnableMapDecode)
	}
	return tx, nil
}

func (b *Base) GetTokenType(ctx contractapi.TransactionContextInterface, tokenId string) (*entity.Token, error) {
	tokenData, err := b.Repo.Get(ctx, doc.Tokens, helper.TokenKey(tokenId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Get token type failed with error  (%s)", err.Error())
		return nil, helper.RespError(errorcode.BizUnableGetTokenType)
	}

	token := new(entity.Token)
	if err = mapstructure.Decode(tokenData, &token); err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Decode token type failed with error  (%s)", err.Error())
		return nil, helper.RespError(errorcode.BizUnableMapDecode)
	}
	return token, nil
}

func (b *Base) GetWallet(ctx contractapi.TransactionContextInterface, walletId string) (*entity.Wallet, error) {
	walletData, err := b.Repo.Get(ctx, doc.Wallets, helper.WalletKey(walletId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Get wallet failed with error  (%s)", err.Error())
		return nil, helper.RespError(errorcode.BizUnableGetWallet)
	}

	wallet := new(entity.Wallet)
	if err = mapstructure.Decode(walletData, &wallet); err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Decode wallet failed with error  (%s)", err.Error())
		return nil, helper.RespError(errorcode.BizUnableMapDecode)
	}
	return wallet, nil
}

func (b *Base) GetActiveWallet(ctx contractapi.TransactionContextInterface, walletId string) (*entity.Wallet, error) {
	wallet, err := b.GetWallet(ctx, walletId)
	if err != nil {
		return nil, err
	}

	if wallet.Status != glossary.Active {
		return nil, helper.RespError(errorcode.InvalidWalletInActive)
	}
	return wallet, nil
}

func (b *Base) GetEnrollment(ctx contractapi.TransactionContextInterface, tokenId string) (*entity.Enrollment, error) {
	enrollmentData, err := b.Repo.Get(ctx, doc.Enrollments, helper.EnrollmentKey(tokenId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Get enrollment failed with error  (%s)", err.Error())
		return nil, helper.RespError(errorcode.BizUnableGetEnrollment)
	}

	enrollment := new(entity.Enrollment)
	if err = mapstructure.Decode(enrollmentData, &enrollment); err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Decode enrollment failed with error  (%s)", err.Error())
		return nil, helper.RespError(errorcode.BizUnableMapDecode)
	}
	return enrollment, nil
}

func (b *Base) GetAndCheckExistEnrollment(ctx contractapi.TransactionContextInterface, tokenId string) (*entity.Enrollment, bool, error) {
	isExisted, enrollmentData, err := b.Repo.GetAndCheckExist(ctx, doc.Enrollments, helper.EnrollmentKey(tokenId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Get enrollment failed with error  (%s)", err.Error())
		return nil, false, helper.RespError(errorcode.BizUnableGetEnrollment)
	}

	if !isExisted {
		return nil, isExisted, nil
	}

	enrollment := new(entity.Enrollment)
	if err = mapstructure.Decode(enrollmentData, &enrollment); err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Decode enrollment failed with error  (%s)", err.Error())
		return nil, isExisted, helper.RespError(errorcode.BizUnableMapDecode)
	}
	return enrollment, isExisted, nil
}

func (b *Base) GetNFT(ctx contractapi.TransactionContextInterface, nftTokenId string) (*entity.NFT, error) {
	nftData, err := b.Repo.Get(ctx, doc.NftToken, helper.NFTKey(nftTokenId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "GetNFT - Get NftToken failed with error (%s)", err.Error())
		return nil, helper.RespError(errorcode.BizUnableGetNFT)
	}

	nftToken := new(entity.NFT)
	if err = mapstructure.Decode(nftData, &nftToken); err != nil {
		glogger.GetInstance().Errorf(ctx, "GetNFT - Decode NftToken token failed with error (%s)", err.Error())
		return nil, helper.RespError(errorcode.BizUnableMapDecode)
	}

	return nftToken, nil
}

func (b *Base) GetBalanceOfToken(ctx contractapi.TransactionContextInterface, domain, walletId string, tokenId string) (*entity.Balance, error) {
	balanceData, err := b.Repo.Get(ctx, domain, helper.BalanceKey(walletId, tokenId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Get balance of token failed with error (%s)", err.Error())
		return nil, helper.RespError(errorcode.BizUnableGetBalance)
	}

	balance := new(entity.Balance)
	if err = mapstructure.Decode(balanceData, &balance); err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Decode balance of token failed with error (%s)", err.Error())
		return nil, helper.RespError(errorcode.BizUnableMapDecode)
	}
	return balance, nil
}

func (b *Base) GetAndCheckBalanceOfToken(ctx contractapi.TransactionContextInterface, domain, walletId string, tokenId string) (*entity.Balance, bool, error) {
	isExisted, balanceData, err := b.Repo.GetAndCheckExist(ctx, domain, helper.BalanceKey(walletId, tokenId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Get balance of token failed with error (%s)", err.Error())
		return nil, false, helper.RespError(errorcode.BizUnableGetBalance)
	}

	if !isExisted {
		glogger.GetInstance().Info(ctx, "Base - Balance of token do not exist in the system")
		return nil, isExisted, nil
	}

	balance := new(entity.Balance)
	if err = mapstructure.Decode(balanceData, &balance); err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Decode balance of token failed with error (%s)", err.Error())
		return nil, isExisted, helper.RespError(errorcode.BizUnableMapDecode)
	}
	return balance, isExisted, nil
}

func (b *Base) ValidatePairWallet(ctx contractapi.TransactionContextInterface, fromWalletId,
	toWalletId string) (walletFrom, walletTo *entity.Wallet, err error) {
	// validate to wallet exist
	walletTo, err = b.GetActiveWallet(ctx, toWalletId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Get To Wallet failed with error (%v)", err)
		return nil, nil, err
	}

	// validate from wallet exist
	walletFrom, err = b.GetActiveWallet(ctx, fromWalletId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Get From Wallet failed with error (%v)", err)
		return nil, nil, err
	}
	return walletFrom, walletTo, err
}

// AddAmount to add amount of balance token
func (b *Base) AddAmount(ctx contractapi.TransactionContextInterface,
	mapCurrentBalance map[string]*entity.BalanceCache, domain, walletId string, tokenId string, amount string) error {
	key := domain + "_" + walletId + "_" + tokenId
	// Load current balance of wallet into memory
	if _, ok := mapCurrentBalance[key]; !ok {
		balanceToken, isExisted, err := b.GetAndCheckBalanceOfToken(ctx, domain, walletId, tokenId)
		if err != nil {
			return err
		}
		balanceCache := new(entity.BalanceCache)
		balanceCache.Domain = domain
		if isExisted {
			balanceCache.IsNew = false
			balanceCache.BalanceEntity = balanceToken
			mapCurrentBalance[key] = balanceCache
		} else {
			balanceEntity := entity.NewBalance(sidechain.Spot, ctx)
			balanceEntity.WalletId = walletId
			balanceEntity.TokenId = tokenId
			balanceEntity.Balances = "0"

			balanceCache.IsNew = true
			balanceCache.BalanceEntity = balanceEntity
			mapCurrentBalance[key] = balanceCache
		}
	}

	// update current balance
	updateCurrentBalance, err := helper.AddBalance(mapCurrentBalance[key].BalanceEntity.Balances, amount)
	if err != nil {
		return err
	}
	mapCurrentBalance[key].BalanceEntity.Balances = updateCurrentBalance

	return nil
}

// SubAmount to sub amount of balance token
func (b *Base) SubAmount(ctx contractapi.TransactionContextInterface,
	mapCurrentBalance map[string]*entity.BalanceCache, domain, walletId string, tokenId string, amount string) error {
	key := domain + "_" + walletId + "_" + tokenId
	// Load current balance of wallet into memory
	if _, ok := mapCurrentBalance[key]; !ok {
		balanceToken, err := b.GetBalanceOfToken(ctx, domain, walletId, tokenId)
		if err != nil {
			return err
		}

		balanceCache := new(entity.BalanceCache)
		balanceCache.IsNew = false
		balanceCache.Domain = domain
		balanceCache.BalanceEntity = balanceToken
		mapCurrentBalance[key] = balanceCache
	}

	// checking current balance with amount
	if helper.CompareStringBalance(mapCurrentBalance[key].BalanceEntity.Balances, amount) < 0 {
		return errors.Errorf("Wallet (%s) do not have enough balance", key)
	}

	// update current balance
	updateCurrentBalance, err := helper.SubBalance(mapCurrentBalance[key].BalanceEntity.Balances, amount)
	if err != nil {
		return err
	}
	mapCurrentBalance[key].BalanceEntity.Balances = updateCurrentBalance

	return nil
}

// RollbackTxHandler to rollback balance of wallet that was updated
func (b *Base) RollbackTxHandler(ctx contractapi.TransactionContextInterface, tx *entity.Transaction,
	mapCurrentBalance map[string]*entity.BalanceCache, step transaction.Step) error {
	// currently only rollback in case transfer or swap
	if step == transaction.SubFromWallet {
		if err := b.AddAmount(ctx, mapCurrentBalance, doc.SpotBalances, tx.FromWallet, tx.FromTokenId, tx.FromTokenAmount); err != nil {
			glogger.GetInstance().Errorf(ctx, "RollbackTxHandler - Add balance of wallet (%s) with transaction (%s) failed", tx.FromWallet, tx.Id)
			return err
		}
	}
	return nil
}

// UpdateBalance to update balance of wallet after handle transaction
func (b *Base) UpdateBalance(ctx contractapi.TransactionContextInterface, mapCurrentBalance map[string]*entity.BalanceCache) error {
	txTime, _ := ctx.GetStub().GetTxTimestamp()
	for key, balanceItem := range mapCurrentBalance {
		if balanceItem.IsNew {
			if err := b.Repo.Create(ctx, balanceItem.BalanceEntity, balanceItem.Domain, helper.BalanceKey(balanceItem.BalanceEntity.WalletId, balanceItem.BalanceEntity.TokenId)); err != nil {
				glogger.GetInstance().Errorf(ctx, "Base - Create new balance of wallet (%s) failed with err (%s)", key, err.Error())
				return helper.RespError(errorcode.BizUnableCreateBalance)
			}
		} else {
			balanceItem.BalanceEntity.UpdatedAt = helper.TimestampISO(txTime.Seconds)
			if err := b.Repo.Update(ctx, balanceItem.BalanceEntity, balanceItem.Domain, helper.BalanceKey(balanceItem.BalanceEntity.WalletId, balanceItem.BalanceEntity.TokenId)); err != nil {
				glogger.GetInstance().Errorf(ctx, "Base - Update balance of wallet (%s) failed with err (%s)", key, err.Error())
				return helper.RespError(errorcode.BizUnableUpdateBalance)
			}
		}
		glogger.GetInstance().Infof(ctx, "WalletId (%s) - Balances update succeed", key)
	}
	return nil
}
