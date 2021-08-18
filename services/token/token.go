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

package token

import (
	"encoding/json"
	"fmt"
	"github.com/Akachain/gringotts/dto/token"
	"github.com/Akachain/gringotts/entity"
	"github.com/Akachain/gringotts/errorcode"
	"github.com/Akachain/gringotts/glossary"
	"github.com/Akachain/gringotts/glossary/doc"
	"github.com/Akachain/gringotts/glossary/transaction"
	"github.com/Akachain/gringotts/helper"
	"github.com/Akachain/gringotts/helper/glogger"
	"github.com/Akachain/gringotts/services/base"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"strings"
)

type tokenService struct {
	*base.Base
}

func NewTokenService() *tokenService {
	return &tokenService{
		base.NewBase(),
	}
}

func (t *tokenService) Transfer(ctx contractapi.TransactionContextInterface, utxoInputKeys []string, utxoOutputs []token.UTXODto, metadata string) error {
	glogger.GetInstance().Info(ctx, "-----------Token Service - Transfer-----------")

	// Handle UTXO of transaction
	lstUtxoInputs, lstUtxoOutputs, err := t.handlerUtxo(ctx, utxoInputKeys, utxoOutputs)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Handle UTXO of transaction failed with err (%s)", err.Error())
		return err
	}

	jsonInputs, _ := json.Marshal(lstUtxoInputs)
	jsonOutputs, _ := json.Marshal(lstUtxoOutputs)

	// create tx burn token
	txTransfer := entity.NewTransaction(ctx)
	txTransfer.InputUTXOs = string(jsonInputs)
	txTransfer.OutputUTXOs = string(jsonOutputs)
	txTransfer.Metadata = metadata
	txTransfer.TxType = transaction.Transfer

	if err := t.Repo.Create(ctx, txTransfer, doc.Transactions, helper.TransactionKey(txTransfer.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Burn - Create transfer transaction failed with error (%v)", err)
		return helper.RespError(errorcode.BizUnableCreateTX)
	}
	glogger.GetInstance().Infof(ctx, "-----------Token Service - Transfer succeed (%s)-----------", txTransfer.Id)

	return nil
}

func (t *tokenService) Mint(ctx contractapi.TransactionContextInterface, walletId, tokenId, amount, metadata string) error {
	glogger.GetInstance().Info(ctx, "-----------Token Service - Mint-----------")

	// validate wallet exited
	if _, err := t.GetActiveWallet(ctx, walletId); err != nil {
		glogger.GetInstance().Errorf(ctx, "Mint - Get wallet mint failed with error (%s)", err.Error())
		return err
	}

	// check total supply with max supply

	// create new UTXO
	id := fmt.Sprintf("%s.%d", ctx.GetStub().GetTxID(), 0)
	utxoEntity := entity.NewUTXO(ctx)
	utxoEntity.WalletId = walletId
	utxoEntity.TokenId = tokenId
	utxoEntity.Amount = amount
	utxoEntity.Id = id
	if err := t.Repo.Create(ctx, utxoEntity, doc.Utxo, helper.UxtoKey(walletId, tokenId, id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Mint - Create UTXO failed with error (%s)", err.Error())
		return helper.RespError(errorcode.BizUnableCreateUTXO)
	}

	utxoJson, err := utxoEntity.ToJsonString()
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Mint - Marshall struct UTXO failed with error (%s)", err.Error())
		return helper.RespError(errorcode.BizUnableCreateTX)
	}

	// create tx mint token
	txMint := entity.NewTransaction(ctx)
	txMint.InputUTXOs = glossary.SystemWallet
	txMint.OutputUTXOs = utxoJson
	txMint.Metadata = metadata
	txMint.TxType = transaction.Mint

	if err := t.Repo.Create(ctx, txMint, doc.Transactions, helper.TransactionKey(txMint.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Mint - Create mint transaction failed with error (%s)", err.Error())
		return helper.RespError(errorcode.BizUnableCreateTX)
	}
	glogger.GetInstance().Infof(ctx, "-----------Token Service - Mint succeed (%s)-----------", txMint.Id)

	return nil
}

func (t *tokenService) Burn(ctx contractapi.TransactionContextInterface, utxoInputKeys []string, utxoOutputs []token.UTXODto, metadata string) error {
	glogger.GetInstance().Info(ctx, "-----------Token Service - Burn-----------")

	// Handle UTXO of transaction
	lstUtxoInputs, lstUtxoOutputs, err := t.handlerUtxo(ctx, utxoInputKeys, utxoOutputs)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Handle UTXO of transaction failed with err (%s)", err.Error())
		return err
	}

	jsonInputs, _ := json.Marshal(lstUtxoInputs)
	jsonOutputs, _ := json.Marshal(lstUtxoOutputs)

	// create tx burn token
	txBurn := entity.NewTransaction(ctx)
	txBurn.InputUTXOs = string(jsonInputs)
	txBurn.OutputUTXOs = string(jsonOutputs)
	txBurn.Metadata = metadata
	txBurn.TxType = transaction.Burn

	if err := t.Repo.Create(ctx, txBurn, doc.Transactions, helper.TransactionKey(txBurn.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Burn - Create burn transaction failed with error (%v)", err)
		return helper.RespError(errorcode.BizUnableCreateTX)
	}
	glogger.GetInstance().Infof(ctx, "-----------Token Service - Burn succeed (%s)-----------", txBurn.Id)

	return nil
}

func (t *tokenService) CreateType(ctx contractapi.TransactionContextInterface, name, tickerToken, maxSupply string) (string, error) {
	glogger.GetInstance().Info(ctx, "-----------Token Service - CreateType-----------")

	tokenEntity := entity.NewToken(ctx)
	tokenEntity.Name = name
	tokenEntity.TickerToken = tickerToken
	tokenEntity.MaxSupply = maxSupply

	if err := t.Repo.Create(ctx, tokenEntity, doc.Tokens, helper.TokenKey(tokenEntity.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "CreateType - Create token type failed with error (%s)", err.Error())
		return "", helper.RespError(errorcode.BizUnableCreateToken)
	}
	glogger.GetInstance().Infof(ctx, "-----------Token Service - CreateType succeed (%s)-----------", tokenEntity.Id)

	return tokenEntity.Id, nil
}

func (t *tokenService) Exchange(ctx contractapi.TransactionContextInterface, pairs []token.ItemPair, metadata string) error {
	glogger.GetInstance().Info(ctx, "-----------Token Service - Exchange-----------")

	// calculate hash and check cache
	stringInput, _ := json.Marshal(pairs)
	inputHash := helper.CalculateHash(string(stringInput) + metadata)
	cacheEntity, isExisted, err := t.GetTxCache(ctx, inputHash)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Exchange - Unable to get transaction cache in database (%s)", err.Error())
		return err
	}
	if isExisted {
		if cacheEntity.Status == transaction.Confirmed {
			return nil
		}
		return errors.New(cacheEntity.ErrorMessage)
	}
	txHandlerCache := entity.NewTxCache(ctx)
	txHandlerCache.Hash = inputHash

	lstUtxoInputs := make([]*entity.UTXO, 0)
	lstUtxoOutputs := make([]*entity.UTXO, 0)
	for _, pair := range pairs {
		// Handle UTXO of transaction
		utxoInputs, utxoOutputs, err := t.handlerUtxo(ctx, pair.Inputs, pair.Outputs)
		if err != nil {
			txHandlerCache.Status = transaction.Rejected
			txHandlerCache.ErrorMessage = err.Error()
			glogger.GetInstance().Errorf(ctx, "Exchange - Handle UTXO of transaction failed with err (%s)", err.Error())
			if err := t.Repo.Update(ctx, txHandlerCache, doc.TxCache, helper.ResultCacheKey(txHandlerCache.Hash)); err != nil {
				glogger.GetInstance().Errorf(ctx, "Exchange - Create cache buy Iao failed with err (%v)", err.Error())
				return helper.RespError(errorcode.BizUnableCreateTxCache)
			}
			return err
		}
		lstUtxoInputs = append(lstUtxoInputs, utxoInputs...)
		lstUtxoOutputs = append(lstUtxoOutputs, utxoOutputs...)
	}

	jsonInputs, _ := json.Marshal(lstUtxoInputs)
	jsonOutputs, _ := json.Marshal(lstUtxoOutputs)

	// create new swap transaction
	txEntity := entity.NewTransaction(ctx)
	txEntity.InputUTXOs = string(jsonInputs)
	txEntity.OutputUTXOs = string(jsonOutputs)
	txEntity.TxType = transaction.Exchange
	txEntity.Metadata = metadata

	if err := t.Repo.Create(ctx, txEntity, doc.Transactions, helper.TransactionKey(txEntity.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Exchange - Create exchange transaction failed with error (%v)", err)
		txHandlerCache.Status = transaction.Rejected
		txHandlerCache.ErrorMessage = err.Error()
		if err := t.Repo.Update(ctx, txHandlerCache, doc.TxCache, helper.ResultCacheKey(txHandlerCache.Hash)); err != nil {
			glogger.GetInstance().Errorf(ctx, "Exchange - Create cache buy Iao failed with err (%v)", err.Error())
			return helper.RespError(errorcode.BizUnableCreateTxCache)
		}
		return helper.RespError(errorcode.BizUnableCreateTX)
	}

	txHandlerCache.Status = transaction.Confirmed
	if err := t.Repo.Update(ctx, txHandlerCache, doc.TxCache, helper.ResultCacheKey(txHandlerCache.Hash)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Exchange - Create cache buy Iao failed with err (%v)", err.Error())
		return helper.RespError(errorcode.BizUnableCreateTxCache)
	}
	glogger.GetInstance().Infof(ctx, "-----------Token Service - Exchange succeed (%s)-----------", txEntity.Id)

	return nil
}

func (t *tokenService) Issue(ctx contractapi.TransactionContextInterface, uxtoInputKeys []string, utxoOutputs []token.UTXODto, metadata string) error {
	glogger.GetInstance().Info(ctx, "-----------Token Service - Issue-----------")

	// validate and summarize amount
	lstUtxoInputs, totalInputAmount, err := t.ValidateAndSummarizeInputs(ctx, uxtoInputKeys)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Issue - Validate issue inputs key failed with err (%s)", err.Error())
		return err
	}

	// get information of issuer
	walletId := lstUtxoInputs[0].WalletId
	tokenId := lstUtxoInputs[0].TokenId

	// check enrollment policy
	enrollment, isExisted, err := t.GetAndCheckExistEnrollment(ctx, tokenId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Issue - Check exit enrollment failed with err (%s)", err.Error())
		return helper.RespError(errorcode.BizUnableGetEnrollment)
	}
	if isExisted {
		// check permission to issue new token
		if enrollment.FromWalletId != "" {
			if !strings.Contains(enrollment.FromWalletId, walletId) {
				glogger.GetInstance().Errorf(ctx, "Issue - From wallet do not have permission issue token (%s)", tokenId)
				return helper.RespError(errorcode.BizIssueNotPermission)
			}
		}
		if enrollment.ToWalletId != "" {
			if !strings.Contains(enrollment.ToWalletId, walletId) {
				glogger.GetInstance().Errorf(ctx, "Issue - To wallet do not have permission issue token (%s)", tokenId)
				return helper.RespError(errorcode.BizIssueNotPermission)
			}
		}
	}

	// get information of asset token
	totalOutputAmount := "0"
	atTokenAmount := "0"
	atTokenId := ""
	lstUtxoOutputs := make([]*entity.UTXO, 0, len(utxoOutputs))
	for i, utxoOutput := range utxoOutputs {
		id := fmt.Sprintf("%s.%d", ctx.GetStub().GetTxID(), i)
		if utxoOutput.TokenId == tokenId {
			totalOutputAmount = utxoOutput.Amount
		} else {
			atTokenId = utxoOutput.TokenId
			atTokenAmount = utxoOutput.Amount
		}
		utxoEntity := entity.NewUTXO(ctx)
		utxoEntity.Id = id
		utxoEntity.WalletId = utxoOutput.WalletId
		utxoEntity.TokenId = utxoOutput.TokenId
		utxoEntity.Amount = utxoOutput.Amount
		lstUtxoOutputs = append(lstUtxoOutputs, utxoEntity)
	}
	if helper.CompareStringBalance(totalOutputAmount, totalInputAmount) >= 0 {
		glogger.GetInstance().Errorf(ctx, "The total amount of outputs (%s) greater than total amount of inputs (%s)", totalOutputAmount, totalInputAmount)
		return helper.RespError(errorcode.BizBalanceNotEnough)
	}

	// validate total supply of AT token
	atToken, err := t.GetTokenType(ctx, atTokenId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Issue - Get AT token failed with err (%s)", err.Error())
		return err
	}

	if atToken.MaxSupply != "" {
		newTotalSupply, err := helper.AddBalance(atToken.TotalSupply, atTokenAmount)
		if err != nil {
			glogger.GetInstance().Errorf(ctx, "Issue - Calculate new total supply failed with err (%s)", err.Error())
			return helper.RespError(errorcode.BizUnableToIssue)
		}

		if helper.CompareStringBalance(atToken.MaxSupply, newTotalSupply) < 0 {
			glogger.GetInstance().Error(ctx, "Issue - Number of new token over the max supply of the token")
			return helper.RespError(errorcode.BizOverMaxSupply)
		}
	}

	// update spent utxo input
	if err := t.UpdateSpentUTXO(ctx, lstUtxoInputs); err != nil {
		glogger.GetInstance().Errorf(ctx, "Issue - Update spent UTXO inputs failed with err (%s)", err.Error())
		return err
	}

	// create new utxo outputs
	if err := t.CreateUTXO(ctx, lstUtxoOutputs); err != nil {
		glogger.GetInstance().Errorf(ctx, "Issue - Create UTXO outputs failed with err (%s)", err.Error())
		return err
	}

	jsonInputs, _ := json.Marshal(lstUtxoInputs)
	jsonOutputs, _ := json.Marshal(lstUtxoOutputs)

	// create tx burn token
	txEntity := entity.NewTransaction(ctx)
	txEntity.InputUTXOs = string(jsonInputs)
	txEntity.OutputUTXOs = string(jsonOutputs)
	txEntity.Metadata = metadata
	txEntity.TxType = transaction.Issue

	if err := t.Repo.Create(ctx, txEntity, doc.Transactions, helper.TransactionKey(txEntity.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "Issue - Create issue transaction failed with error (%v)", err)
		return helper.RespError(errorcode.BizUnableCreateTX)
	}
	glogger.GetInstance().Infof(ctx, "-----------Token Service - Issue succeed (%s)-----------", txEntity.Id)

	return nil
}

func (t *tokenService) handlerUtxo(ctx contractapi.TransactionContextInterface, utxoInputKeys []string,
	utxoOutputs []token.UTXODto) (lstUtxoInputs []*entity.UTXO, lstUtxoOutputs []*entity.UTXO, err error) {
	totalInputAmount := "0"
	totalOutputAmount := "0"
	// Validate and summarize utxo inputs
	lstUtxoInputs, totalInputAmount, err = t.ValidateAndSummarizeInputs(ctx, utxoInputKeys)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Validate UTXO inputs failed with err (%s)", err.Error())
		return nil, nil, err
	}

	// Validate and summarize utxo outputs
	lstUtxoOutputs, totalOutputAmount, err = t.ValidateAndSummarizeOutputs(ctx, utxoOutputs)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "Validate UTXO outputs failed with err (%s)", err.Error())
		return nil, nil, err
	}

	// Validate total inputs equals total outputs
	if helper.CompareStringBalance(totalInputAmount, totalOutputAmount) != 0 {
		glogger.GetInstance().Errorf(ctx, "Total utxoInput amount %d does not equal total utxoOutput amount %d", totalInputAmount, totalOutputAmount)
		return nil, nil, helper.RespError(errorcode.BizAmountNotEqual)
	}

	// update spent utxo input
	if err := t.UpdateSpentUTXO(ctx, lstUtxoInputs); err != nil {
		glogger.GetInstance().Errorf(ctx, "Update spent UTXO inputs failed with err (%s)", err.Error())
		return nil, nil, err
	}

	// create new utxo outputs
	if err := t.CreateUTXO(ctx, lstUtxoOutputs); err != nil {
		glogger.GetInstance().Errorf(ctx, "Create UTXO outputs failed with err (%s)", err.Error())
		return nil, nil, err
	}
	return lstUtxoInputs, lstUtxoOutputs, nil
}
