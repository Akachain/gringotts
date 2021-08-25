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
	"fmt"
	"github.com/Akachain/gringotts/dto/token"
	"github.com/Akachain/gringotts/entity"
	"github.com/Akachain/gringotts/errorcode"
	"github.com/Akachain/gringotts/glossary"
	"github.com/Akachain/gringotts/glossary/doc"
	"github.com/Akachain/gringotts/helper"
	"github.com/Akachain/gringotts/helper/glogger"
	"github.com/Akachain/gringotts/repository"
	"github.com/Akachain/gringotts/repository/base"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/mitchellh/mapstructure"
	"strings"
)

type Base struct {
	Repo repository.Repo
}

func NewBase() *Base {
	return &Base{
		Repo: base.NewRepository(),
	}
}

func (b *Base) GetUTXOFullKey(ctx contractapi.TransactionContextInterface, fullKey string) (*entity.UTXO, error) {
	compositeKeyPath := strings.Split(fullKey, glossary.SplitKey)
	return b.GetUTXO(ctx, compositeKeyPath[2], compositeKeyPath[3], compositeKeyPath[4])
}

func (b *Base) GetUTXO(ctx contractapi.TransactionContextInterface, walletId, tokenId, txId string) (*entity.UTXO, error) {
	utxoData, err := b.Repo.Get(ctx, doc.Utxo, helper.UxtoKey(walletId, tokenId, txId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Get UTXO (%s) failed with error (%s)", walletId+"_"+tokenId+"_"+txId, err.Error())
		return nil, helper.RespError(errorcode.BizUnableGetUTXO)
	}

	utxoEntity := entity.NewUTXO()
	if err = mapstructure.Decode(utxoData, &utxoEntity); err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Decode UTXO failed with error  (%s)", err.Error())
		return nil, helper.RespError(errorcode.BizUnableMapDecode)
	}
	return utxoEntity, nil
}

func (b *Base) GetTxCache(ctx contractapi.TransactionContextInterface, txCacheId string) (*entity.TxCache, bool, error) {
	isExisted, txCacheData, err := b.Repo.GetAndCheckExist(ctx, doc.TxCache, helper.ResultCacheKey(txCacheId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Get transaction cache (%s) failed with error (%s)", txCacheId, err.Error())
		return nil, isExisted, helper.RespError(errorcode.BizUnableGetTxCache)
	}

	txCacheEntity := entity.NewTxCache()
	if err = mapstructure.Decode(txCacheData, &txCacheEntity); err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Decode IAO failed with error  (%s)", err.Error())
		return nil, isExisted, helper.RespError(errorcode.BizUnableMapDecode)
	}
	return txCacheEntity, isExisted, nil
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

	tokenEntity := new(entity.Token)
	if err = mapstructure.Decode(tokenData, &tokenEntity); err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Decode token type failed with error  (%s)", err.Error())
		return nil, helper.RespError(errorcode.BizUnableMapDecode)
	}
	return tokenEntity, nil
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

// ValidateAndSummarizeInputs to validate and summarize utxo inputs
func (b *Base) ValidateAndSummarizeInputs(ctx contractapi.TransactionContextInterface, utxoInputKeys []string) ([]*entity.UTXO, string, error) {
	utxoInputs := make(map[string]*entity.UTXO)
	totalInputAmount := "0"
	for _, utxoInputKey := range utxoInputKeys {
		if utxoInputs[utxoInputKey] != nil {
			glogger.GetInstance().Error(ctx, "The same utxo input can not be spend twice")
			return nil, "", helper.RespError(errorcode.DuplicateUtxo)
		}

		utxoEntity, err := b.GetUTXOFullKey(ctx, utxoInputKey)
		if err != nil {
			glogger.GetInstance().Errorf(ctx, "Get UTXO failed with err (%s)", err.Error())
			return nil, "", err
		}
		if utxoEntity.Status != 0 {
			glogger.GetInstance().Errorf(ctx, "UTXO (%s) have status spent", utxoEntity.Id)
			return nil, "", helper.RespError(errorcode.InvalidateStatusUtxo)
		}

		tempAmount, _ := helper.AddBalance(totalInputAmount, utxoEntity.Amount)
		totalInputAmount = tempAmount
		utxoInputs[utxoInputKey] = utxoEntity
	}

	lstUtxo := make([]*entity.UTXO, 0, len(utxoInputs))
	for _, utxo := range utxoInputs {
		lstUtxo = append(lstUtxo, utxo)
	}
	return lstUtxo, totalInputAmount, nil
}

// ValidateAndSummarizeOutputs to validate and summarize utxo outputs
func (b *Base) ValidateAndSummarizeOutputs(ctx contractapi.TransactionContextInterface, utxoOutputs []token.UTXODto) ([]*entity.UTXO, string, error) {
	totalOutputAmount := "0"
	txID := ctx.GetStub().GetTxID()
	lstUtxoOutput := make([]*entity.UTXO, 0, len(utxoOutputs))
	for i, utxoOutput := range utxoOutputs {
		if helper.CompareStringBalance(utxoOutput.Amount, "0") <= 0 {
			glogger.GetInstance().Errorf(ctx, "UTXO (%s) outputs have amount negative", utxoOutput.WalletId+"_"+utxoOutput.TokenId)
			return nil, "", helper.RespError(errorcode.InvalidateAmountUtxo)
		}

		tempAmount, _ := helper.AddBalance(totalOutputAmount, utxoOutput.Amount)
		totalOutputAmount = tempAmount
		id := fmt.Sprintf("%s.%d", txID, i)
		utxoOutputEntity := entity.NewUTXO(ctx)
		utxoOutputEntity.Id = id
		utxoOutputEntity.WalletId = utxoOutput.WalletId
		utxoOutputEntity.TokenId = utxoOutput.TokenId
		utxoOutputEntity.Amount = utxoOutput.Amount
		if utxoOutput.WalletId == glossary.SystemWallet {
			utxoOutputEntity.Status = 1
		}
		lstUtxoOutput = append(lstUtxoOutput, utxoOutputEntity)
	}
	return lstUtxoOutput, totalOutputAmount, nil
}

func (b *Base) UpdateSpentUTXO(ctx contractapi.TransactionContextInterface, lstUtxo []*entity.UTXO) error {
	txTime, _ := ctx.GetStub().GetTxTimestamp()
	for _, utxo := range lstUtxo {
		utxo.UpdatedAt = helper.TimestampISO(txTime.Seconds)
		utxo.Status = 1
		if err := b.Repo.Update(ctx, utxo, doc.Utxo, helper.UxtoKey(utxo.WalletId, utxo.TokenId, utxo.Id)); err != nil {
			glogger.GetInstance().Errorf(ctx, "Base - Update UTXO (%s) failed with err (%s)", utxo.Id, err.Error())
			return helper.RespError(errorcode.BizUnableCreateUTXO)
		}
	}
	return nil
}

func (b *Base) CreateUTXO(ctx contractapi.TransactionContextInterface, lstUtxo []*entity.UTXO) error {
	for _, utxo := range lstUtxo {
		if err := b.Repo.Update(ctx, utxo, doc.Utxo, helper.UxtoKey(utxo.WalletId, utxo.TokenId, utxo.Id)); err != nil {
			glogger.GetInstance().Errorf(ctx, "Base - Create UTXO (%s) failed with err (%s)", utxo.Id, err.Error())
			return helper.RespError(errorcode.BizUnableCreateUTXO)
		}
	}
	return nil
}
