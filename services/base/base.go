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
	"encoding/json"
	"github.com/Akachain/gringotts/entity"
	"github.com/Akachain/gringotts/errorcode"
	"github.com/Akachain/gringotts/glossary"
	"github.com/Akachain/gringotts/glossary/doc"
	"github.com/Akachain/gringotts/helper"
	"github.com/Akachain/gringotts/helper/glogger"
	"github.com/Akachain/gringotts/pkg/query"
	"github.com/Akachain/gringotts/repository"
	"github.com/Akachain/gringotts/repository/base"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/mitchellh/mapstructure"
)

type Base struct {
	Repo repository.Repo
}

func NewBase() *Base {
	return &Base{
		Repo: base.NewRepository(),
	}
}

func (b *Base) GetTransaction(ctx contractapi.TransactionContextInterface, txId string) (*entity.Transaction, error) {
	txData, err := b.Repo.Get(ctx, doc.Transactions, helper.TransactionKey(txId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Get transaction (%s) failed with error (%v)", txId, err)
		return nil, helper.RespError(errorcode.BizUnableGetTx)
	}

	tx := entity.NewTransaction()
	if err = mapstructure.Decode(txData, &tx); err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Decode transaction failed with error (%v)", err)
		return nil, helper.RespError(errorcode.BizUnableMapDecode)
	}
	return tx, nil
}

func (b *Base) GetTokenType(ctx contractapi.TransactionContextInterface, tokenId string) (*entity.Token, error) {
	tokenData, err := b.Repo.Get(ctx, doc.Tokens, helper.TokenKey(tokenId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Get token type failed with error (%v)", err)
		return nil, helper.RespError(errorcode.BizUnableGetTokenType)
	}

	token := new(entity.Token)
	if err = mapstructure.Decode(tokenData, &token); err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Decode token type failed with error (%v)", err)
		return nil, helper.RespError(errorcode.BizUnableMapDecode)
	}
	return token, nil
}

func (b *Base) GetWallet(ctx contractapi.TransactionContextInterface, walletId string) (*entity.Wallet, error) {
	walletData, err := b.Repo.Get(ctx, doc.Wallets, helper.WalletKey(walletId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Get wallet failed with error (%v)", err)
		return nil, helper.RespError(errorcode.BizUnableGetWallet)
	}

	wallet := new(entity.Wallet)
	if err = mapstructure.Decode(walletData, &wallet); err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Decode wallet failed with error (%v)", err)
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

func (b *Base) AddBalance(ctx contractapi.TransactionContextInterface, walletId string, amount string) error {
	txTime, _ := ctx.GetStub().GetTxTimestamp()
	//  update balance of wallet
	walletData, err := b.Repo.Get(ctx, doc.Wallets, helper.WalletKey(walletId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Get wallet failed with error (%v)", err)
		return helper.RespError(errorcode.BizUnableGetWallet)
	}

	wallet := new(entity.Wallet)
	if err = mapstructure.Decode(walletData, &wallet); err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Decode wallet failed with error (%v)", err)
		return helper.RespError(errorcode.BizUnableMapDecode)
	}

	updatedBalance, err := helper.AddBalance(wallet.Balances, amount)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Sub balance of wallet failed with error (%v)", err)
		return helper.RespError(errorcode.BizUnableUpdateWallet)
	}
	wallet.Balances = updatedBalance
	wallet.UpdatedAt = helper.TimestampISO(txTime.Seconds)
	if err := b.Repo.Update(ctx, wallet, doc.Wallets, helper.WalletKey(wallet.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Transfer - Update wallet failed with error (%v)", err)
		return helper.RespError(errorcode.BizUnableUpdateWallet)
	}
	return nil
}

func (b *Base) SubBalance(ctx contractapi.TransactionContextInterface, walletId string, amount string) error {
	txTime, _ := ctx.GetStub().GetTxTimestamp()
	//  update balance of wallet
	walletData, err := b.Repo.Get(ctx, doc.Wallets, helper.WalletKey(walletId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Get wallet failed with error (%v)", err)
		return helper.RespError(errorcode.BizUnableGetWallet)
	}

	wallet := new(entity.Wallet)
	if err = mapstructure.Decode(walletData, &wallet); err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Decode wallet failed with error (%v)", err)
		return helper.RespError(errorcode.BizUnableMapDecode)
	}

	// check balance insufficient
	if helper.CompareStringBalance(wallet.Balances, amount) < 0 {
		glogger.GetInstance().Error(ctx, "Base - Balance of wallet  insufficient")
		return helper.RespError(errorcode.BizBalanceNotEnough)
	}

	updatedBalance, err := helper.SubBalance(wallet.Balances, amount)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Sub balance of wallet failed with error (%v)", err)
		return helper.RespError(errorcode.BizUnableUpdateWallet)
	}
	wallet.Balances = updatedBalance
	wallet.UpdatedAt = helper.TimestampISO(txTime.Seconds)
	if err := b.Repo.Update(ctx, wallet, doc.Wallets, helper.WalletKey(wallet.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Transfer - Update from wallet failed with error (%v)", err)
		return helper.RespError(errorcode.BizUnableUpdateWallet)
	}
	return nil
}

func (b *Base) GetEnrollment(ctx contractapi.TransactionContextInterface, tokenId string) (*entity.Enrollment, error) {
	enrollmentData, err := b.Repo.Get(ctx, doc.Enrollments, helper.EnrollmentKey(tokenId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Get enrollment failed with error (%v)", err)
		return nil, helper.RespError(errorcode.BizUnableGetEnrollment)
	}

	enrollment := new(entity.Enrollment)
	if err = mapstructure.Decode(enrollmentData, &enrollment); err != nil {
		glogger.GetInstance().Errorf(ctx, "Base - Decode enrollment failed with error (%v)", err)
		return nil, helper.RespError(errorcode.BizUnableMapDecode)
	}
	return enrollment, nil
}

func (b *Base) GetNFT(ctx contractapi.TransactionContextInterface, nftTokenId string) (*entity.NFT, error) {
	nftData, err := b.Repo.Get(ctx, doc.NFT, helper.NFTKey(nftTokenId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "GetNFT - Get NFT failed with error (%v)", err)
		return nil, helper.RespError(errorcode.BizUnableGetNFT)
	}

	nftToken := new(entity.NFT)
	if err = mapstructure.Decode(nftData, &nftToken); err != nil {
		glogger.GetInstance().Errorf(ctx, "GetNFT - Decode NFT token failed with error (%v)", err)
		return nil, helper.RespError(errorcode.BizUnableMapDecode)
	}

	return nftToken, nil
}

func (b *Base) GetExchangeTxByBlockchainId(ctx contractapi.TransactionContextInterface, blockchainId string) ([]*entity.Transaction, error) {
	txList := make([]*entity.Transaction, 0)
	resultsIterator, err := b.Repo.GetQueryString(ctx, query.GetTransactionByBlockchainId(blockchainId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "GetExchangeTxByBlockchainId - Get query exchange tx failed with error (%v)", err)
		return nil, helper.RespError(errorcode.BizUnableGetTx)
	}
	defer resultsIterator.Close()

	// Check data response after query in database
	if !resultsIterator.HasNext() {
		// Return with list transaction id empty
		return txList, nil
	}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			glogger.GetInstance().Errorf(ctx, "GetExchangeTxByBlockchainId - Start query failed with error (%v)", err)
			return nil, helper.RespError(errorcode.BizUnableGetTx)
		}
		tx := entity.NewTransaction()
		err = json.Unmarshal(queryResponse.Value, tx)
		if err != nil {
			glogger.GetInstance().Errorf(ctx, "GetExchangeTxByBlockchainId - Unable to unmarshal transaction error (%v)", err)
			return txList, helper.RespError(errorcode.BizUnableGetTx)
		}

		txList = append(txList, tx)
	}

	return txList, nil
}
