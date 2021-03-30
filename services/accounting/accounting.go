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
	"github.com/Akachain/gringotts/errorcode"
	"github.com/Akachain/gringotts/glossary"
	"github.com/Akachain/gringotts/glossary/doc"
	"github.com/Akachain/gringotts/glossary/transaction"
	"github.com/Akachain/gringotts/helper"
	"github.com/Akachain/gringotts/helper/glogger"
	"github.com/Akachain/gringotts/internal/entity"
	"github.com/Akachain/gringotts/pkg/query"
	"github.com/Akachain/gringotts/services/base"
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
		txIdList = append(txIdList, tx.Id)
	}
	glogger.GetInstance().Infof(ctx, "GetTx - List transaction have status pending: (%s)", strings.Join(txIdList, ","))

	return txIdList, nil
}

func (a *accountingService) CalculateBalance(ctx contractapi.TransactionContextInterface, accountingDto dto.AccountingBalance) error {
	glogger.GetInstance().Infof(ctx, "CalculateBalance - List transaction: (%s)", strings.Join(accountingDto.TxId, ","))
	// map temp balance
	mapCurrentBalance := make(map[string]string, len(accountingDto.TxId)*2)

	for _, id := range accountingDto.TxId {
		tx, err := a.GetTransaction(ctx, id)
		if err != nil {
			glogger.GetInstance().Errorf(ctx, "CalculateBalance - Get transaction (%s) failed with error (%v)", id, err)
			continue
		}
		glogger.GetInstance().Infof(ctx, "Transaction get from db: %v", tx)

		// check status of transaction
		if tx.Status != transaction.Pending {
			glogger.GetInstance().Errorf(ctx, "CalculateBalance - Transaction (%s) has status (%s)", id, tx.Status)
			continue
		}

		// check type of transaction
		switch tx.TxType {
		case transaction.Transfer:
			if tx.FromWallet == glossary.SystemWallet || tx.ToWallet == glossary.SystemWallet {
				glogger.GetInstance().Errorf(ctx, "CalculateBalance - Transfer - Transaction (%s) has from/to wallet Id is system type", id)
				tx.Status = transaction.Rejected
				goto UpdateTx
			}

			if err := a.subAmount(ctx, mapCurrentBalance, tx.FromWallet, tx.Amount); err != nil {
				glogger.GetInstance().Errorf(ctx, "CalculateBalance - Transfer - Transaction (%s): Unable to sub temp amount of From wallet", id)
				tx.Status = transaction.Rejected
				goto UpdateTx
			}

			if err := a.addAmount(ctx, mapCurrentBalance, tx.ToWallet, tx.Amount); err != nil {
				glogger.GetInstance().Errorf(ctx, "CalculateBalance - Transfer - Transaction (%s): Unable to add temp amount of To wallet", id)
				tx.Status = transaction.Rejected
				goto UpdateTx
			}

			tx.Status = transaction.Confirmed
			break
		case transaction.Mint:
			if tx.FromWallet != glossary.SystemWallet {
				glogger.GetInstance().Errorf(ctx, "CalculateBalance - Mint - Transaction (%s): has From wallet Id is not system type", id)
				tx.Status = transaction.Rejected
				goto UpdateTx
			}
			if err := a.addAmount(ctx, mapCurrentBalance, tx.ToWallet, tx.Amount); err != nil {
				glogger.GetInstance().Errorf(ctx, "CalculateBalance - Mint - Transaction (%s): add balance failed (%v)", id, err)
				tx.Status = transaction.Rejected
				goto UpdateTx
			}
			tx.Status = transaction.Confirmed
			break
		case transaction.Burn:
			if tx.ToWallet != glossary.SystemWallet {
				glogger.GetInstance().Errorf(ctx, "CalculateBalance - Burn - Transaction (%s): has To wallet Id is not system type", id)
				tx.Status = transaction.Rejected
				goto UpdateTx
			}
			if err := a.subAmount(ctx, mapCurrentBalance, tx.FromWallet, tx.Amount); err != nil {
				glogger.GetInstance().Errorf(ctx, "CalculateBalance - Burn - Transaction (%s) sub balance failed (%v)", id, err)
				tx.Status = transaction.Rejected
				goto UpdateTx
			}
			tx.Status = transaction.Confirmed
			break
		default:
			glogger.GetInstance().Errorf(ctx, "CalculateBalance - Transaction (%s) has type (%s) not support", id, tx.TxType)
			tx.Status = transaction.Rejected
		}
	UpdateTx:
		glogger.GetInstance().Infof(ctx, "CalculateBalance - Update Transaction get from db: %v", tx)
		if err := a.Repo.Update(ctx, tx, doc.Transactions, helper.TransactionKey(tx.Id)); err != nil {
			glogger.GetInstance().Errorf(ctx, "CalculateBalance - Update transaction (%s) failed with err (%v)", id, err)
		}
		glogger.GetInstance().Infof(ctx, "CalculateBalance - Update transaction (%s) succeed", id)
	}

	// Log debug
	glogger.GetInstance().Infof(ctx, "CalculateBalance - Length of mapCurrentBalance(%d)", len(mapCurrentBalance))

	// Update balance of wallet after calculate total balance update of wallet
	if err := a.updateBalance(ctx, mapCurrentBalance); err != nil {
		return err
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
