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

package nft_transfer

import (
	"github.com/Akachain/gringotts/entity"
	"github.com/Akachain/gringotts/errorcode"
	"github.com/Akachain/gringotts/glossary/doc"
	"github.com/Akachain/gringotts/glossary/transaction"
	"github.com/Akachain/gringotts/helper"
	"github.com/Akachain/gringotts/helper/glogger"
	"github.com/Akachain/gringotts/pkg/tx"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type txNftTransfer struct {
	*tx.TxBase
}

func NewTxNftTransfer() tx.Handler {
	return &txNftTransfer{tx.NewTxBase()}
}

func (t *txNftTransfer) AccountingTx(ctx contractapi.TransactionContextInterface, tx *entity.Transaction, mapBalanceToken map[string]string) (*entity.Transaction, error) {
	// TODO: handle rollback map current balance when sub/add failed
	txUpdate, err := t.TxHandlerTransfer(ctx, mapBalanceToken, tx)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "NftTransfer - Handler token transfer failed with err (%s)", err.Error())
		return txUpdate, err
	}

	// handler owner of nft
	nftToken, err := t.GetNFT(ctx, tx.ToTokenId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "NftTransfer - Get NftToken failed with error (%s)", err.Error())
		txUpdate.Status = transaction.Rejected
		return txUpdate, err
	}

	if nftToken.OwnerId != tx.ToWallet {
		glogger.GetInstance().Error(ctx, "NftTransfer - To wallet not match owner of nft token")
		txUpdate.Status = transaction.Rejected
		return txUpdate, helper.RespError(errorcode.BizNftNotPermission)
	}

	nftToken.OwnerId = tx.FromWallet
	if err := t.Repo.Update(ctx, nftToken, doc.NftToken, helper.NFTKey(nftToken.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "TransferFrom - Update NftToken failed with error (%v)", err)
		txUpdate.Status = transaction.Rejected
		return txUpdate, helper.RespError(errorcode.BizUnableUpdateNFT)
	}

	txUpdate.Status = transaction.Confirmed
	return txUpdate, nil
}
