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

package helper

import (
	"github.com/Akachain/gringotts/glossary"
	"gotest.tools/assert"
	"math/big"
	"testing"
)

func TestSubBalance(t *testing.T) {
	sub, err := SubBalance("0", "1234")
	assert.NilError(t, err, "Fail to sub balance")
	assert.Equal(t, sub, "-1234")
}

func TestAddBalance(t *testing.T) {
	sub, err := SubBalance("0", "1234")
	assert.NilError(t, err, "Fail to sub balance")
	assert.Equal(t, sub, "-1234")

	add, err := AddBalance(sub, "1234")
	assert.NilError(t, err, "Fail to add balance")
	assert.Equal(t, add, "0")
}

func TestBigInt(t *testing.T) {
	f := 0.333333
	bInt := new(big.Int)
	bInt.SetInt64(int64(f * glossary.AkcBase))
	t.Log(bInt.Int64())

	currB := new(big.Int)
	currB.SetInt64(333333333)

	newB := currB.Mul(currB, bInt)
	t.Log(newB)

	divNumber := new(big.Int)
	divNumber.SetInt64(glossary.AkcBase)
	t.Log(divNumber.Int64())

	t.Log(newB.Div(newB, divNumber).Int64())
	t.Log(newB.Mod(newB, divNumber))
}
