package burn

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

type txBurn struct {
	*tx.TxBase
}

func NewTxBurn() tx.Handler {
	return &txBurn{tx.NewTxBase()}
}

func (t *txBurn) AccountingTx(ctx contractapi.TransactionContextInterface, tx *entity.Transaction, mapBalanceToken map[string]string) (*entity.Transaction, error) {
	if tx.ToWallet != glossary.SystemWallet {
		glogger.GetInstance().Errorf(ctx, "TxBurn - Transaction (%s): has To wallet Id is not system type", tx.Id)
		tx.Status = transaction.Rejected
		return tx, errors.New("To wallet id invalidate")
	}

	// decrease total supply of token on the blockchain
	tokenType, err := t.GetTokenType(ctx, tx.FromTokenId)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "TxBurn - Transaction (%s): Unable to get token type", tx.Id)
		tx.Status = transaction.Rejected
		return nil, err
	}
	totalSupplyUpdated, err := helper.SubBalance(tokenType.TotalSupply, tx.FromTokenAmount)
	if err != nil {
		glogger.GetInstance().Errorf(ctx, "TxBurn - Transaction (%s): Unable to decrease total supply of token", tx.Id)
		tx.Status = transaction.Rejected
		return nil, err
	}
	tokenType.TotalSupply = totalSupplyUpdated
	if err := t.Repo.Update(ctx, tokenType, doc.Tokens, helper.TokenKey(tokenType.Id)); err != nil {
		glogger.GetInstance().Errorf(ctx, "TxBurn - Transaction (%s): Unable to update token type (%s)", tx.Id, err.Error())
		tx.Status = transaction.Rejected
		return tx, errors.New("Unable to decrease total of token on the blockchain")
	}

	if err := t.SubAmount(ctx, mapBalanceToken, tx.FromWallet, tx.FromTokenId, tx.FromTokenAmount); err != nil {
		glogger.GetInstance().Errorf(ctx, "TxBurn - Transaction (%s): add balance failed (%s)", tx.Id, err.Error())
		tx.Status = transaction.Rejected
		return tx, errors.New("Sub balance of to wallet failed")
	}

	tx.Status = transaction.Confirmed
	return tx, nil
}
