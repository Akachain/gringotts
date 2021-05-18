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
	"github.com/Akachain/gringotts/services/nft"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type NftHandler struct {
	nftService services.NFT
}

func NewNftHandler() NftHandler {
	return NftHandler{nft.NewNftService()}
}

func (n *NftHandler) Mint(ctx contractapi.TransactionContextInterface, mintNFT dto.MintNFT) (string, error) {
	glogger.GetInstance().Info(ctx, "-----------NFT Handler - Mint-----------")

	// checking dto validate
	if err := mintNFT.IsValid(); err != nil {
		glogger.GetInstance().Errorf(ctx, "NFT Handler - Mint Input invalidate %v", err)
		return "", helper.RespError(errorcode.InvalidParam)
	}

	return n.nftService.Mint(ctx, mintNFT.GS1Number, mintNFT.OwnerWalletId, mintNFT.Metadata, mintNFT.HashData)
}

func (n *NftHandler) OwnerOf(ctx contractapi.TransactionContextInterface, ownerNFT dto.OwnerNFT) (string, error) {
	glogger.GetInstance().Info(ctx, "-----------NFT Handler - OwnerOf-----------")

	// checking dto validate
	if err := ownerNFT.IsValid(); err != nil {
		glogger.GetInstance().Errorf(ctx, "NFT Handler - Owner Of NFT Input invalidate %v", err)
		return "", helper.RespError(errorcode.InvalidParam)
	}

	return n.nftService.OwnerOf(ctx, ownerNFT.NFTTokenId)
}

func (n *NftHandler) BalanceOf(ctx contractapi.TransactionContextInterface, balanceOfNFT dto.BalanceOfNFT) (int, error) {
	glogger.GetInstance().Info(ctx, "-----------NFT Handler - BalanceOf-----------")

	// checking dto validate
	if err := balanceOfNFT.IsValid(); err != nil {
		glogger.GetInstance().Errorf(ctx, "NFT Handler - BalanceOf Input invalidate %v", err)
		return -1, helper.RespError(errorcode.InvalidParam)
	}

	return n.nftService.BalanceOf(ctx, balanceOfNFT.OwnerWalletId)
}

func (n *NftHandler) TransferNFT(ctx contractapi.TransactionContextInterface, transferNFT dto.TransferNFT) error {
	glogger.GetInstance().Info(ctx, "-----------NFT Handler - TransferNFT-----------")

	// checking dto validate
	if err := transferNFT.IsValid(); err != nil {
		glogger.GetInstance().Errorf(ctx, "NFT Handler - TransferNFT Input invalidate %v", err)
		return helper.RespError(errorcode.InvalidParam)
	}

	return n.nftService.TransferFrom(ctx, transferNFT.OwnerWalletId, transferNFT.ToWalletId, transferNFT.NFTTokenId)
}
