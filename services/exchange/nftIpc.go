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
	"github.com/Akachain/gringotts/pkg/ipc"
	"github.com/Akachain/gringotts/services/base"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/mitchellh/mapstructure"
)

type nftIpc struct {
	*base.Base
}

func NewNftIpc() ipc.Ipc {
	return &nftIpc{base.NewBase()}
}

func (n *nftIpc) TransactionCallback(ctx contractapi.TransactionContextInterface, txId string, txStatus transaction.Status) error {
	txTime, _ := ctx.GetStub().GetTxTimestamp()
	exchangeTx, err := n.getExchangeTx(ctx, txId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "TransactionCallback - Get exchange tx failed with error (%v)", err)
		return err
	}

	if exchangeTx.Status != transaction.Pending {
		glogger.GetInstance().Error(ctx, "TransactionCallback - Exchange transaction has status difference pending")
		return helper.RespError(errorcode.BizExchangeTxInvalidStatus)
	}

	if txStatus == transaction.Rejected {
		exchangeTx.Status = transaction.Rejected
		if err := n.updateExchangeTx(ctx, txId, exchangeTx); err != nil {
			return err
		}
		return nil
	}

	nftToken, err := n.GetNFT(ctx, exchangeTx.NftTokenId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "TransactionCallback - Get NftToken token failed with error (%v)", err)
		return err
	}

	// update owner of nft token
	nftToken.OwnerId = exchangeTx.ToWalletId
	nftToken.UpdatedAt = helper.TimestampISO(txTime.Seconds)
	if err := n.Repo.Update(ctx, nftToken, doc.NftToken, helper.NFTKey(nftToken.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "TransactionCallback - Update Nft token (%s) failed with err (%v)", nftToken.Id, err)
		return helper.RespError(errorcode.BizUnableUpdateNFT)
	}

	// update status of exchange transaction
	exchangeTx.Status = transaction.Confirmed
	if err := n.updateExchangeTx(ctx, txId, exchangeTx); err != nil {
		return err
	}

	return nil
}

func (n *nftIpc) getExchangeTx(ctx contractapi.TransactionContextInterface, txId string) (*entity.Exchange, error) {
	exchangeTxData, err := n.Repo.Get(ctx, doc.Exchange, helper.ExchangeKey(txId))
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "NFTIpc - Get NftToken failed with error (%v)", err)
		return nil, helper.RespError(errorcode.BizUnableGetExchange)
	}

	exchangeTx := new(entity.Exchange)
	if err = mapstructure.Decode(exchangeTxData, &exchangeTx); err != nil {
		glogger.GetInstance().Errorf(ctx, "NFTIpc - Decode Exchange Tx failed with error (%v)", err)
		return nil, helper.RespError(errorcode.BizUnableMapDecode)
	}

	return exchangeTx, nil
}

func (n *nftIpc) updateExchangeTx(ctx contractapi.TransactionContextInterface, txId string, exchangeTx *entity.Exchange) error {
	txTime, _ := ctx.GetStub().GetTxTimestamp()
	exchangeTx.UpdatedAt = helper.TimestampISO(txTime.Seconds)
	if err := n.Repo.Update(ctx, exchangeTx, doc.Exchange, helper.ExchangeKey(txId)); err != nil {
		glogger.GetInstance().Errorf(ctx, "TransactionCallback - Update exchange tx (%s) failed with err (%v)", txId, err)
		return helper.RespError(errorcode.BizUnableUpdateExchange)
	}

	return nil
}
