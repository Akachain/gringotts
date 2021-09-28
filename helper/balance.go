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

// Package helper contains common library functions that are shared among business
// related functions.
package helper

import (
	"github.com/Akachain/gringotts/pkg/unit"
	"github.com/pkg/errors"
	"math/big"
)

// AddBalance add amount into current balance
func AddBalance(currentBalance string, amount string) (string, error) {
	if amount == "" {
		return "", errors.New("invalidate amount")
	}

	// convert balance to akc unit
	curBalanceUnit := unit.NewBalanceUnitFromString(currentBalance)
	if curBalanceUnit.Cmp(big.NewInt(0)) < 0 {
		return "", errors.New("Current balance less than zero")
	}
	// convert amount to akc unit
	amountUnit := unit.NewBalanceUnitFromString(amount)
	if amountUnit.Cmp(big.NewInt(0)) < 0 {
		return "", errors.New("Unable to add negative amount number")
	}

	if err := curBalanceUnit.AddBalance(amountUnit); err != nil {
		return "", err
	}
	return curBalanceUnit.String(), nil
}

// SubBalance sub current balance with amount
func SubBalance(currentBalance string, amount string) (string, error) {
	if amount == "" {
		return "", errors.New("invalidate amount")
	}

	// convert balance to akc unit
	curBalanceUnit := unit.NewBalanceUnitFromString(currentBalance)
	if curBalanceUnit.Cmp(big.NewInt(0)) < 0 {
		return "", errors.New("Current balance less than zero")
	}
	// convert amount to akc unit
	amountUnit := unit.NewBalanceUnitFromString(amount)
	if amountUnit.Cmp(big.NewInt(0)) < 0 {
		return "", errors.New("Unable to sub negative amount number")
	}

	if err := curBalanceUnit.SubBalance(amountUnit); err != nil {
		return "", err
	}
	return curBalanceUnit.String(), nil
}

// MulBalance mul current balance with amount
func MulBalance(currentBalance string, amount int64) (string, error) {
	if amount <= 0 {
		return "", errors.New("invalidate amount")
	}

	// convert balance to akc unit
	curBalanceUnit := unit.NewBalanceUnitFromString(currentBalance)

	if err := curBalanceUnit.MulBalance(amount); err != nil {
		return "", err
	}
	return curBalanceUnit.String(), nil
}

// CompareStringBalance to compare between current balance and amount.
// Amount is string type.
// Return 1 if current balance greater than amount. Otherwise return -1
func CompareStringBalance(currentBalance string, amount string) int {
	// convert balance to akc unit
	curBalanceUnit := unit.NewBalanceUnitFromString(currentBalance)
	// convert amount to akc unit
	amountUnit := unit.NewBalanceUnitFromString(amount)

	return curBalanceUnit.Cmp(amountUnit.Int)
}

// CompareFloatBalance to compare between current balance and amount.
// Amount is float64 type.
// Return 1 if current balance greater than amount. Otherwise return -1
func CompareFloatBalance(currentBalance string, amount float64) int {
	// convert balance to akc unit
	curBalanceUnit := unit.NewBalanceUnitFromString(currentBalance)
	// convert amount to akc unit
	amountUnit := unit.NewBalanceUnitFromFloat(amount)

	return curBalanceUnit.Cmp(amountUnit.Int)
}
