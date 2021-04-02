package entity

import (
	"github.com/Akachain/gringotts/glossary/doc"
	"github.com/Akachain/gringotts/helper"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Enrollment struct {
	TokenId      string `mapstructure:"tokenId" json:"tokenId"`
	FromWalletId string `mapstructure:"fromWallet" json:"fromWallet"`
	ToWalletId   string `mapstructure:"toWallet" json:"toWallet"`
	Base         `mapstructure:",squash"`
}

func NewEnrollment(ctx ...contractapi.TransactionContextInterface) *Enrollment {
	if len(ctx) <= 0 {
		return &Enrollment{}
	}
	txTime, _ := ctx[0].GetStub().GetTxTimestamp()
	return &Enrollment{
		Base: Base{
			Id:           helper.GenerateID(doc.Enrollments, ctx[0].GetStub().GetTxID()),
			CreatedAt:    helper.TimestampISO(txTime.Seconds),
			UpdatedAt:    helper.TimestampISO(txTime.Seconds),
			BlockChainId: ctx[0].GetStub().GetTxID(),
		},
	}
}
