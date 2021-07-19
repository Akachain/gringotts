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

// Package errorcode contains a list of pre-defined error code
// to return to client
package errorcode

type ErrorCode string

type ErrorStruct struct {
	ErrCode       ErrorCode
	MessageDetail string
}

const (
	// System error code
	Succeed ErrorCode = "0"

	// validate error code
	InvalidArg            ErrorCode = "100"
	InvalidParam          ErrorCode = "101"
	ValidationFail        ErrorCode = "102"
	InvalidWalletInActive ErrorCode = "103"

	// Business error code
	BizUnableParse              ErrorCode = "300"
	BizUnableCreateTX           ErrorCode = "301"
	BizUnableCreateWallet       ErrorCode = "302"
	BizUnableGetWallet          ErrorCode = "303"
	BizUnableUpdateWallet       ErrorCode = "304"
	BizUnableMintToken          ErrorCode = "305"
	BizUnableBurnToken          ErrorCode = "306"
	BizUnableCreateToken        ErrorCode = "307"
	BizBalanceNotEnough         ErrorCode = "308"
	BizUnableMapDecode          ErrorCode = "309"
	BizUnableGetTokenType       ErrorCode = "310"
	BizUnableGetTx              ErrorCode = "311"
	BizUnableApproveAllowance   ErrorCode = "312"
	BizUnableGetAllowance       ErrorCode = "313"
	BizAllowanceNotEnough       ErrorCode = "314"
	BizUnableUpdateAllowance    ErrorCode = "315"
	BizUnableTransferDiffType   ErrorCode = "316"
	BizUnableGetEnrollment      ErrorCode = "317"
	BizIssueNotPermission       ErrorCode = "318"
	BizUnableCreateEnrollment   ErrorCode = "319"
	BizUnableUpdateTX           ErrorCode = "320"
	BizUnableCreateNFT          ErrorCode = "321"
	BizUnableGetNFT             ErrorCode = "322"
	BizUnableUpdateNFT          ErrorCode = "323"
	BizUnableCreateExchange     ErrorCode = "324"
	BizUnableGetExchange        ErrorCode = "325"
	BizExchangeTxInvalidStatus  ErrorCode = "326"
	BizUnableUpdateExchange     ErrorCode = "327"
	BizUnableCreateBalance      ErrorCode = "328"
	BizUnableUpdateBalance      ErrorCode = "329"
	BizUnableGetBalance         ErrorCode = "330"
	BizNftNotPermission         ErrorCode = "331"
	BizUnableCreateAsset        ErrorCode = "332"
	BizUnableGetAsset           ErrorCode = "333"
	BizUnableCreateIao          ErrorCode = "334"
	BizUnableGetIao             ErrorCode = "335"
	BizUnableUpdateAsset        ErrorCode = "336"
	BizOverMaxSupply            ErrorCode = "337"
	BizUnableToIssue            ErrorCode = "338"
	BizUnableUpdateIao          ErrorCode = "339"
	BizUnableCreateInvestorBook ErrorCode = "340"
	BizUnableCreateBuyCache     ErrorCode = "341"
	BizUnableGetBuyCache        ErrorCode = "342"
	BizUnableGetInvestorBook    ErrorCode = "343"
	BizUnableUpdateInvestorBook ErrorCode = "344"
)

var mapErrorCode = map[ErrorCode]string{
	Succeed:                     "Invoke transaction succeed",
	InvalidArg:                  "Incorrect number of arguments",
	InvalidParam:                "Parameter input invalidate",
	InvalidWalletInActive:       "Wallet has status inactive",
	BizUnableParse:              "Unable to parse argument",
	BizUnableCreateTX:           "Unable to create transaction on blockchain",
	BizUnableCreateWallet:       "Unable to create wallet on blockchain",
	BizUnableGetWallet:          "Unable to get wallet on blockchain",
	BizUnableUpdateWallet:       "Unable to update wallet on blockchain",
	BizUnableMintToken:          "Unable to mint token for wallet on blockchain",
	BizUnableBurnToken:          "Unable to burn token for wallet on blockchain",
	BizUnableCreateToken:        "Unable to create token type on blockchain",
	BizBalanceNotEnough:         "Balance of wallet insufficient",
	BizUnableMapDecode:          "Unable to parse entity on blockchain",
	BizUnableGetTokenType:       "Unable to get token type on blockchain",
	BizUnableGetTx:              "Unable to get list transaction on blockchain",
	BizUnableApproveAllowance:   "Unable to approve allowance of spender wallet on blockchain",
	BizUnableGetAllowance:       "Unable to get allowance of spender wallet on blockchain",
	BizAllowanceNotEnough:       "Transfer amount greater than allowance of spender",
	BizUnableUpdateAllowance:    "Unable update allowance of spender wallet on blockchain",
	BizUnableTransferDiffType:   "Unable to transfer token between wallet have different token type",
	BizUnableGetEnrollment:      "Unable to get enrollment on blockchain",
	BizIssueNotPermission:       "From/To wallet do not have permission to issue new token",
	BizUnableCreateEnrollment:   "Unable to create/update wallet enrollment",
	BizUnableUpdateTX:           "Unable to update transaction on blockchain",
	BizUnableCreateNFT:          "Unable to create new NFT token on blockchain",
	BizUnableGetNFT:             "Unable to get NFT token on blockchain",
	BizUnableUpdateNFT:          "Unable to update NFT token on blockchain",
	BizUnableCreateExchange:     "Unable to create new exchange nft record on blockchain",
	BizUnableGetExchange:        "Unable to get exchange nft record on blockchain",
	BizExchangeTxInvalidStatus:  "Exchange transaction have status difference pending on the blockchain",
	BizUnableUpdateExchange:     "Unable to update exchange nft transaction on the blockchain",
	BizUnableCreateBalance:      "Unable to create new balance of token on the blockchain",
	BizUnableUpdateBalance:      "Unable to update balance of token on the blockchain",
	BizUnableGetBalance:         "Unable to get balance of token on the blockchain",
	BizNftNotPermission:         "Wallet Id not match owner of nft token",
	BizUnableCreateAsset:        "Unable to create new asset on the blockchain",
	BizUnableGetAsset:           "Unable to get asset on the blockchain",
	BizUnableCreateIao:          "Unable to create new iao campaign on the blockchain",
	BizUnableGetIao:             "Unable to get Iao campaign on the blockchain",
	BizUnableUpdateAsset:        "Unable to update asset on the blockchain",
	BizOverMaxSupply:            "New token issue over the max supply",
	BizUnableToIssue:            "Unable to issue new token on the blockchain",
	BizUnableUpdateIao:          "Unable to update IAO on the blockchain",
	BizUnableCreateInvestorBook: "Unable to create new Investor Book on the blockchain",
	BizUnableCreateBuyCache:     "Unable to create buy iao cache on the blockchain",
	BizUnableGetBuyCache:        "Unable to get buy iao cache on the blockchain",
	BizUnableGetInvestorBook:    "Unable to get investor book on the blockchain",
	BizUnableUpdateInvestorBook: "Unable to update investor book on the blockchain",
}

func (e ErrorCode) Message() string {
	return mapErrorCode[e]
}

func (e ErrorCode) Code() string {
	return string(e)
}

func (e ErrorStruct) Error() string {
	msg := "\"status\":\"" + e.ErrCode.Code() + "\", \"message\":\"" + e.ErrCode.Message() + "\", \"messageDetail\":\"" + e.MessageDetail + "\""
	return msg
}
