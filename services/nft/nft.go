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

package nft

import (
	"github.com/Akachain/gringotts/entity"
	"github.com/Akachain/gringotts/errorcode"
	"github.com/Akachain/gringotts/glossary/doc"
	"github.com/Akachain/gringotts/glossary/sidechain"
	"github.com/Akachain/gringotts/glossary/transaction"
	"github.com/Akachain/gringotts/helper"
	"github.com/Akachain/gringotts/helper/glogger"
	"github.com/Akachain/gringotts/pkg/unit"
	"github.com/Akachain/gringotts/services"
	"github.com/Akachain/gringotts/services/base"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type nftService struct {
	*base.Base
}

func NewNftService() services.NFT {
	return &nftService{
		base.NewBase(),
	}
}

func (n *nftService) Mint(ctx contractapi.TransactionContextInterface, gs1Number string, ownerWalletId string, metaData string, hashData string) (string, error) {
	glogger.GetInstance().Info(ctx, "-----------NftToken Service - Mint-----------")

	if _, err := n.GetActiveWallet(ctx, ownerWalletId); err != nil {
		glogger.GetInstance().Errorf(ctx, "Mint - Get owner wallet failed with error (%v)", err)
		return "", err
	}

	nftEntity := entity.NewNFT(ctx)
	nftEntity.GS1Number = gs1Number
	nftEntity.OwnerId = ownerWalletId
	nftEntity.MetaData = metaData
	nftEntity.HashData = hashData

	if err := n.Repo.Create(ctx, nftEntity, doc.NftToken, helper.NFTKey(nftEntity.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "NftToken Service - Mint NftToken failed with error (%v)", err)
		return "", helper.RespError(errorcode.BizUnableCreateNFT)
	}

	// create balance
	balanceEntity := entity.NewBalance(sidechain.Spot, ctx)
	balanceEntity.WalletId = ownerWalletId
	balanceEntity.TokenId = nftEntity.Id
	balanceEntity.Balances = "1"
	if err := n.Repo.Create(ctx, balanceEntity, doc.SpotBalances, helper.BalanceKey(ownerWalletId, doc.NftToken, nftEntity.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Create - Init balance of stable token failed with error (%v)", err)
		return "", helper.RespError(errorcode.BizUnableCreateBalance)
	}

	glogger.GetInstance().Infof(ctx, "-----------NftToken Service - Transfer succeed (%s)-----------", nftEntity.Id)

	return nftEntity.Id, nil
}

func (n *nftService) OwnerOf(ctx contractapi.TransactionContextInterface, nftTokenId string) (string, error) {
	glogger.GetInstance().Info(ctx, "-----------NftToken Service - OwnerOf-----------")

	nftToken, err := n.GetNFT(ctx, nftTokenId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "OwnerOf - Get NftToken failed with error (%v)", err)
		return "", err
	}

	return nftToken.OwnerId, nil
}

func (n *nftService) BalanceOf(ctx contractapi.TransactionContextInterface, ownerWalletId string) (int, error) {
	glogger.GetInstance().Info(ctx, "-----------NftToken Service - BalanceOf-----------")
	// TODO implementation logic get number nft of owner wallet
	return -1, nil
}

func (n *nftService) TransferFrom(ctx contractapi.TransactionContextInterface, fromWalletId string, toWalletId string,
	fromTokenId string, nftTokenId string, price float64) error {
	glogger.GetInstance().Info(ctx, "-----------NftToken Service - TransferFrom-----------")

	if _, _, err := n.ValidatePairWallet(ctx, fromWalletId, toWalletId); err != nil {
		glogger.GetInstance().Errorf(ctx, "TransferFrom - Validation pair wallet nft failed with error (%s)", err.Error())
		return err
	}

	// handler owner of nft
	nftToken, err := n.GetNFT(ctx, nftTokenId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "TransferFrom - Get NftToken failed with error (%s)", err.Error())
		return err
	}

	if nftToken.OwnerId != toWalletId {
		glogger.GetInstance().Error(ctx, "NftTransfer - To wallet not match owner of nft token")
		return helper.RespError(errorcode.BizNftNotPermission)
	}

	// convert price to akc base
	amountUnit := unit.NewBalanceUnitFromFloat(price)

	// create new swap transaction
	txEntity := entity.NewTransaction(ctx)
	txEntity.SpenderWallet = fromWalletId
	txEntity.FromWallet = fromWalletId
	txEntity.ToWallet = toWalletId
	txEntity.FromTokenId = fromTokenId
	txEntity.ToTokenId = nftTokenId
	txEntity.TxType = transaction.TransferNft
	txEntity.FromTokenAmount = amountUnit.String()
	txEntity.ToTokenAmount = amountUnit.String()

	if err := n.Repo.Create(ctx, txEntity, doc.Transactions, helper.TransactionKey(txEntity.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "TransferFrom - Create transfer nft transaction failed with error (%v)", err)
		return helper.RespError(errorcode.BizUnableCreateTX)
	}
	glogger.GetInstance().Infof(ctx, "-----------NftToken Service - Transfer nft succeed (%s)-----------", txEntity.Id)

	return nil
}
