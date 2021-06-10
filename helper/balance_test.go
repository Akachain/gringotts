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

func TestCompareFloatBalance(t *testing.T) {
	res := CompareFloatBalance("100000000", 1)
	assert.Equal(t, res, 0)

	res = CompareFloatBalance("100000000", 2)
	assert.Equal(t, res, -1)

	res = CompareFloatBalance("1000000000", 1)
	assert.Equal(t, res, 1)
}

func TestCompareStringBalance(t *testing.T) {
	res := CompareStringBalance("100000000", "100000000")
	assert.Equal(t, res, 0)

	res = CompareStringBalance("100000000", "200000000")
	assert.Equal(t, res, -1)

	res = CompareStringBalance("1000000000", "100000000")
	assert.Equal(t, res, 1)
}
