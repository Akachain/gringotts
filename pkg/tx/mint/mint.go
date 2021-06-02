package mint

import (
	"github.com/Akachain/gringotts/entity"
	"github.com/Akachain/gringotts/glossary"
	"github.com/Akachain/gringotts/glossary/doc"
	"github.com/Akachain/gringotts/glossary/transaction"
	"github.com/Akachain/gringotts/helper"
	"github.com/Akachain/gringotts/helper/glogger"
	"github.com/Akachain/gringotts/pkg/tx"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
)

type txMint struct {
	*tx.TxBase
}

func NewTxMint() tx.Handler {
	return &txMint{tx.NewTxBase()}
}

func (t *txMint) AccountingTx(ctx contractapi.TransactionContextInterface, tx *entity.Transaction, mapBalanceToken map[string]string) (*entity.Transaction, error) {
	if tx.FromWallet != glossary.SystemWallet {
		glogger.GetInstance().Errorf(ctx, "TxMint - Transaction (%s): has From wallet Id is not system type", tx.Id)
		tx.Status = transaction.Rejected
		return tx, errors.New("From wallet id invalidate")
	}

	// increase total supply of token on the blockchain
	tokenType, err := t.GetTokenType(ctx, tx.ToTokenId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "TxMint - Transaction (%s): Unable to get token type", tx.Id)
		tx.Status = transaction.Rejected
		return nil, err
	}
	totalSupplyUpdated, err := helper.AddBalance(tokenType.TotalSupply, tx.ToTokenAmount)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "TxMint - Transaction (%s): Unable to increase total supply of token", tx.Id)
		tx.Status = transaction.Rejected
		return nil, err
	}
	tokenType.TotalSupply = totalSupplyUpdated
	if err := t.Repo.Update(ctx, tokenType, doc.Tokens, helper.TokenKey(tokenType.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "TxMint - Transaction (%s): Unable to update token type (%s)", tx.Id, err.Error())
		tx.Status = transaction.Rejected
		return tx, errors.New("Unable to increase total of token on the blockchain")
	}

	if err := t.AddAmount(ctx, mapBalanceToken, tx.ToWallet, tx.ToTokenId, tx.ToTokenAmount); err != nil {
		glogger.GetInstance().Errorf(ctx, "TxMint - Transaction (%s): add balance failed (%s)", tx.Id, err.Error())
		tx.Status = transaction.Rejected
		return tx, errors.New("Add balance of to wallet failed")
	}

	tx.Status = transaction.Confirmed
	return tx, nil
}
