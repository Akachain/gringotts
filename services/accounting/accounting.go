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

package accounting

import (
	"encoding/json"
	"github.com/Akachain/gringotts/dto"
	"github.com/Akachain/gringotts/entity"
	"github.com/Akachain/gringotts/errorcode"
	"github.com/Akachain/gringotts/glossary"
	"github.com/Akachain/gringotts/glossary/doc"
	"github.com/Akachain/gringotts/glossary/transaction"
	"github.com/Akachain/gringotts/helper"
	"github.com/Akachain/gringotts/helper/glogger"
	"github.com/Akachain/gringotts/pkg/query"
	"github.com/Akachain/gringotts/services/base"
	"github.com/davecgh/go-spew/spew"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"strings"
)

type accountingService struct {
	*base.Base
}

func NewAccountingService() *accountingService {
	return &accountingService{base.NewBase()}
}

func (a *accountingService) GetTx(ctx contractapi.TransactionContextInterface) ([]string, error) {
	glogger.GetInstance().Debugf(ctx, "GetTx - Get Query String %s", query.GetPendingTransactionQueryString())
	var txIdList []string
	resultsIterator, err := a.Repo.GetQueryStringWithPagination(ctx, query.GetPendingTransactionQueryString())
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "GetTx - Get query with paging failed with error (%v)", err)
		return nil, helper.RespError(errorcode.BizUnableGetTx)
	}
	defer resultsIterator.Close()

	// Check data response after query in database
	if !resultsIterator.HasNext() {
		// Return with list transaction id empty
		return txIdList, nil
	}

	mapFilterExchangeTx := make(map[string]string, 0)
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			glogger.GetInstance().Errorf(ctx, "GetTx - Start query failed with error (%v)", err)
			return nil, helper.RespError(errorcode.BizUnableGetTx)
		}
		tx := entity.NewTransaction()
		err = json.Unmarshal(queryResponse.Value, tx)
		if err != nil {
			glogger.GetInstance().Error(ctx, "GetTx - Unable to unmarshal transaction")
			continue
		}

		// check transaction is exchange or not
		if tx.TxType == transaction.Exchange {
			if _, ok := mapFilterExchangeTx[tx.BlockChainId]; !ok {
				mapFilterExchangeTx[tx.BlockChainId] = tx.Id
				txIdList = append(txIdList, tx.Id)
			}
			continue
		}

		// add transaction for otherwise
		txIdList = append(txIdList, tx.Id)
	}
	glogger.GetInstance().Infof(ctx, "GetTx - List transaction have status pending: (%s)", strings.Join(txIdList, ","))

	return txIdList, nil
}

func (a *accountingService) CalculateBalance(ctx contractapi.TransactionContextInterface, accountingDto dto.AccountingBalance) error {
	glogger.GetInstance().Infof(ctx, "CalculateBalance - List transaction: (%s)", strings.Join(accountingDto.TxId, ","))
	txTime, _ := ctx.GetStub().GetTxTimestamp()
	// map temp balance
	mapCurrentBalance := make(map[string]string, len(accountingDto.TxId)*2)
	lstTx := make([]*entity.Transaction, 0, len(accountingDto.TxId))

	for _, id := range accountingDto.TxId {
		tx, err := a.GetTransaction(ctx, id)
		if err != nil {
			glogger.GetInstance().Errorf(ctx, "CalculateBalance - Get transaction (%s) failed with error (%v)", id, err)
			continue
		}
		glogger.GetInstance().Debugf(ctx, "CalculateBalance - Transaction get from db: %s", spew.Sdump(tx))

		// check status of transaction
		if tx.Status != transaction.Pending {
			glogger.GetInstance().Errorf(ctx, "CalculateBalance - Transaction (%s) has status (%s)", id, tx.Status)
			continue
		}

		// update time when update transaction
		tx.UpdatedAt = helper.TimestampISO(txTime.Seconds)

		// check type of transaction
		txHandleStep, err := a.txHandler(ctx, tx, mapCurrentBalance)
		lstTx = append(lstTx, tx)
		if err != nil {
			glogger.GetInstance().Errorf(ctx, "CalculateBalance - Handle transaction (%s) failed with error (%v)", id, err)
			if txHandleStep != transaction.Validation {
				if err := a.rollbackTxHandler(ctx, tx, mapCurrentBalance, txHandleStep); err != nil {
					glogger.GetInstance().Errorf(ctx, "CalculateBalance - Rollback handle transaction (%s) failed with error (%v)", id, err)
				}
			}
			continue
		}
	}

	// Update transaction status
	if err := a.updateTransaction(ctx, lstTx); err != nil {
		return err
	}

	// Update balance of wallet after calculate total balance update of wallet
	if err := a.updateBalance(ctx, mapCurrentBalance); err != nil {
		return err
	}

	return nil
}

// txHandler to handle transaction and update calculate current balance of each wallet
func (a *accountingService) txHandler(ctx contractapi.TransactionContextInterface, tx *entity.Transaction,
	mapCurrentBalance map[string]string) (transaction.Step, error) {
	switch tx.TxType {
	case transaction.Transfer, transaction.Swap, transaction.Issue:
		if tx.FromWallet == glossary.SystemWallet || tx.ToWallet == glossary.SystemWallet {
			glogger.GetInstance().Errorf(ctx, "TxHandler - Transfer - Transaction (%s) has from/to wallet Id is system type", tx.Id)
			tx.Status = transaction.Rejected
			return transaction.Validation, errors.New("From/To wallet id invalidate")
		}

		if err := a.subAmount(ctx, mapCurrentBalance, tx.FromWallet, tx.Amount); err != nil {
			glogger.GetInstance().Errorf(ctx, "TxHandler - Transfer - Transaction (%s): Unable to sub temp amount of From wallet", tx.Id)
			tx.Status = transaction.Rejected
			return transaction.Validation, errors.WithMessage(err, "Sub balance of from wallet failed")
		}

		if err := a.addAmount(ctx, mapCurrentBalance, tx.ToWallet, tx.Amount); err != nil {
			glogger.GetInstance().Errorf(ctx, "TxHandler - Transfer - Transaction (%s): Unable to add temp amount of To wallet", tx.Id)
			tx.Status = transaction.Rejected
			return transaction.SubFromWallet, errors.WithMessage(err, "Add balance of to wallet failed")
		}

		tx.Status = transaction.Confirmed
		break
	case transaction.Mint:
		if tx.FromWallet != glossary.SystemWallet {
			glogger.GetInstance().Errorf(ctx, "TxHandler - Mint - Transaction (%s): has From wallet Id is not system type", tx.Id)
			tx.Status = transaction.Rejected
			return transaction.Validation, errors.New("From wallet id invalidate")
		}
		if err := a.addAmount(ctx, mapCurrentBalance, tx.ToWallet, tx.Amount); err != nil {
			glogger.GetInstance().Errorf(ctx, "TxHandler - Mint - Transaction (%s): add balance failed (%v)", tx.Id, err)
			tx.Status = transaction.Rejected
			return transaction.Validation, errors.New("Add balance of to wallet failed")
		}
		tx.Status = transaction.Confirmed
		break
	case transaction.Burn:
		if tx.ToWallet != glossary.SystemWallet {
			glogger.GetInstance().Errorf(ctx, "TxHandler - Burn - Transaction (%s): has To wallet Id is not system type", tx.Id)
			tx.Status = transaction.Rejected
			return transaction.Validation, errors.New("To wallet id invalidate")
		}
		if err := a.subAmount(ctx, mapCurrentBalance, tx.FromWallet, tx.Amount); err != nil {
			glogger.GetInstance().Errorf(ctx, "TxHandler - Burn - Transaction (%s) sub balance failed (%v)", tx.Id, err)
			tx.Status = transaction.Rejected
			return transaction.Validation, errors.New("Sub balance of from wallet failed")
		}
		tx.Status = transaction.Confirmed
		break
	case transaction.Exchange:
		txList, err := a.GetExchangeTxByBlockchainId(ctx, tx.BlockChainId)
		if err != nil {
			glogger.GetInstance().Errorf(ctx, "TxHandler - Exchange - Get tx with blockchain id (%s) failed (%v)", tx.BlockChainId, err)
			return transaction.Validation, err
		}

		if len(txList) <= 1 {
			glogger.GetInstance().Error(ctx, "TxHandler - Exchange - Get exchange tx less than one")
			return transaction.Validation, errors.New("Exchange transaction need have a pair")
		}

		for _, subTx := range txList {
			if err := a.subAmount(ctx, mapCurrentBalance, subTx.FromWallet, subTx.Amount); err != nil {
				glogger.GetInstance().Errorf(ctx, "TxHandler - Exchange - Transaction (%s): Unable to sub temp amount of From wallet", subTx.Id)
				tx.Status = transaction.Rejected
				return transaction.Validation, errors.WithMessage(err, "Sub balance of from wallet failed")
			}

			if err := a.addAmount(ctx, mapCurrentBalance, subTx.ToWallet, subTx.Amount); err != nil {
				glogger.GetInstance().Errorf(ctx, "TxHandler - Exchange - Transaction (%s): Unable to add temp amount of To wallet", subTx.Id)
				tx.Status = transaction.Rejected
				return transaction.Validation, errors.WithMessage(err, "Add balance of to wallet failed")
			}
			subTx.Status = transaction.Confirmed
		}
	default:
		glogger.GetInstance().Errorf(ctx, "TxHandler - Transaction (%s) has type (%s) not support", tx.Id, string(tx.TxType))
		tx.Status = transaction.Rejected
	}

	return transaction.Validation, nil
}

// rollbackTxHandler to rollback balance of wallet that was updated
func (a *accountingService) rollbackTxHandler(ctx contractapi.TransactionContextInterface, tx *entity.Transaction,
	mapCurrentBalance map[string]string, step transaction.Step) error {
	// currently only rollback in case transfer or swap
	if step == transaction.SubFromWallet {
		if err := a.addAmount(ctx, mapCurrentBalance, tx.FromWallet, tx.Amount); err != nil {
			glogger.GetInstance().Errorf(ctx, "RollbackTxHandler - Add balance of wallet (%s) with transaction (%s) failed", tx.FromWallet, tx.Id)
			return err
		}
	}
	return nil
}

// addAmount to add amount for wallet
func (a *accountingService) addAmount(ctx contractapi.TransactionContextInterface,
	mapCurrentBalance map[string]string, walletId string, amount string) error {
	// Load current balance of wallet into memory
	if _, ok := mapCurrentBalance[walletId]; !ok {
		wallet, err := a.GetWallet(ctx, walletId)
		if err != nil {
			return err
		}
		mapCurrentBalance[walletId] = wallet.Balances
	}

	// update current balance
	updateCurrentBalance, err := helper.AddBalance(mapCurrentBalance[walletId], amount)
	if err != nil {
		return err
	}
	mapCurrentBalance[walletId] = updateCurrentBalance

	return nil
}

// subAmount to sub amount for wallet
func (a *accountingService) subAmount(ctx contractapi.TransactionContextInterface,
	mapCurrentBalance map[string]string, walletId string, amount string) error {
	// Load current balance of wallet into memory
	if _, ok := mapCurrentBalance[walletId]; !ok {
		wallet, err := a.GetWallet(ctx, walletId)
		if err != nil {
			return err
		}
		mapCurrentBalance[walletId] = wallet.Balances
	}

	// checking current balance with amount
	if helper.CompareStringBalance(mapCurrentBalance[walletId], amount) <= 0 {
		return errors.Errorf("Wallet (%s) do not have enough balance", walletId)
	}

	// update current balance
	updateCurrentBalance, err := helper.SubBalance(mapCurrentBalance[walletId], amount)
	if err != nil {
		return err
	}
	mapCurrentBalance[walletId] = updateCurrentBalance

	return nil
}

// updateBalance to update balance of wallet after handle transaction
func (a *accountingService) updateBalance(ctx contractapi.TransactionContextInterface, mapCurrentBalance map[string]string) error {
	for key, value := range mapCurrentBalance {
		wallet, err := a.GetWallet(ctx, key)
		if err != nil {
			glogger.GetInstance().Errorf(ctx, "updateBalance - Get wallet with id (%s) failed with err (%v)", key, err)
			return helper.RespError(errorcode.BizUnableGetWallet)
		}
		wallet.Balances = value
		if err := a.Repo.Update(ctx, wallet, doc.Wallets, helper.WalletKey(key)); err != nil {
			glogger.GetInstance().Errorf(ctx, "updateBalance - Update balance of wallet (%s) failed with err (%v)", key, err)
			return helper.RespError(errorcode.BizUnableUpdateWallet)
		}
		glogger.GetInstance().Infof(ctx, "WalletId (%s) - Balances update succeed", key)
	}
	return nil
}

// updateTransaction update status of transaction after handle transaction
func (a *accountingService) updateTransaction(ctx contractapi.TransactionContextInterface, lstTx []*entity.Transaction) error {
	for _, tx := range lstTx {
		if err := a.Repo.Update(ctx, tx, doc.Transactions, helper.TransactionKey(tx.Id)); err != nil {
			glogger.GetInstance().Errorf(ctx, "UpdateTransaction - Update transaction (%s) failed with err (%v)", tx.Id, err)
			return helper.RespError(errorcode.BizUnableUpdateTX)
		}
		glogger.GetInstance().Infof(ctx, "UpdateTransaction - Update transaction (%s) succeed", tx.Id)
	}

	return nil
}
