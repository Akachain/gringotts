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

package unit

import (
	"errors"
	"github.com/Akachain/gringotts/glossary"
	"math/big"
)

type BalanceUnit struct {
	*big.Int
}

func NewBalanceUnit() *BalanceUnit {
	return &BalanceUnit{
		new(big.Int),
	}
}

func NewBalanceUnitFromString(balanceS string) *BalanceUnit {
	bUnit := NewBalanceUnit()
	bUnit.SetStringUnit(balanceS)
	return bUnit
}

func NewBalanceUnitFromFloat(balanceF float64) *BalanceUnit {
	bUnit := NewBalanceUnit()
	bUnit.SetFloatUnit(balanceF)
	return bUnit
}

// SetFloatUnit to set amount or balance float to base unit
func (b *BalanceUnit) SetFloatUnit(value float64) {
	valueToUnit := value * glossary.AkcBase
	// TODO: need refactor
	if b.Int == nil {
		b.Int = new(big.Int)
	}
	b.SetInt64(int64(valueToUnit))
}

func (b *BalanceUnit) SetStringUnit(value string) {
	// TODO: need refactor
	if b.Int == nil {
		b.Int = new(big.Int)
	}
	b.SetString(value, glossary.StringUnitBase)
}

// AddBalance add amount into current balance
func (b *BalanceUnit) AddBalance(amount *BalanceUnit) error {
	if amount == nil {
		return errors.New("invalidate amount unit")
	}
	b.Add(b.Int, amount.Int)
	return nil
}

// SubBalance subtracts the current balance with amount
func (b *BalanceUnit) SubBalance(amount *BalanceUnit) error {
	if amount == nil {
		return errors.New("invalidate amount unit")
	}
	b.Sub(b.Int, amount.Int)
	return nil
}
