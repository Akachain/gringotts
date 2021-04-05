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
	BizUnableParse            ErrorCode = "300"
	BizUnableCreateTX         ErrorCode = "301"
	BizUnableCreateWallet     ErrorCode = "302"
	BizUnableGetWallet        ErrorCode = "303"
	BizUnableUpdateWallet     ErrorCode = "304"
	BizUnableMintToken        ErrorCode = "305"
	BizUnableBurnToken        ErrorCode = "306"
	BizUnableCreateToken      ErrorCode = "307"
	BizBalanceNotEnough       ErrorCode = "308"
	BizUnableMapDecode        ErrorCode = "309"
	BizUnableGetTokenType     ErrorCode = "310"
	BizUnableGetTx            ErrorCode = "311"
	BizUnableApproveAllowance ErrorCode = "312"
	BizUnableGetAllowance     ErrorCode = "313"
	BizAllowanceNotEnough     ErrorCode = "314"
	BizUnableUpdateAllowance  ErrorCode = "315"
	BizUnableTransferDiffType ErrorCode = "316"
	BizUnableGetEnrollment    ErrorCode = "317"
	BizIssueNotPermission     ErrorCode = "318"
	BizUnableCreateEnrollment ErrorCode = "319"
)

var mapErrorCode = map[ErrorCode]string{
	Succeed:                   "Invoke transaction succeed",
	InvalidArg:                "Incorrect number of arguments",
	InvalidParam:              "Parameter input invalidate",
	InvalidWalletInActive:     "Wallet has status inactive",
	BizUnableParse:            "Unable to parse argument",
	BizUnableCreateTX:         "Unable to create transaction on blockchain",
	BizUnableCreateWallet:     "Unable to create wallet on blockchain",
	BizUnableGetWallet:        "Unable to get wallet on blockchain",
	BizUnableUpdateWallet:     "Unable to update wallet on blockchain",
	BizUnableMintToken:        "Unable to mint token for wallet on blockchain",
	BizUnableBurnToken:        "Unable to burn token for wallet on blockchain",
	BizUnableCreateToken:      "Unable to create token type on blockchain",
	BizBalanceNotEnough:       "Balance of wallet insufficient",
	BizUnableMapDecode:        "Unable to parse entity on blockchain",
	BizUnableGetTokenType:     "Unable to get token type on blockchain",
	BizUnableGetTx:            "Unable to get list transaction on blockchain",
	BizUnableApproveAllowance: "Unable to approve allowance of spender wallet on blockchain",
	BizUnableGetAllowance:     "Unable to get allowance of spender wallet on blockchain",
	BizAllowanceNotEnough:     "Transfer amount greater than allowance of spender",
	BizUnableUpdateAllowance:  "Unable update allowance of spender wallet on blockchain",
	BizUnableTransferDiffType: "Unable to transfer token between wallet have different token type",
	BizUnableGetEnrollment:    "Unable to get enrollment on blockchain",
	BizIssueNotPermission:     "From/To wallet do not have permission to issue new token",
	BizUnableCreateEnrollment: "Unable to create/update wallet enrollment",
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
