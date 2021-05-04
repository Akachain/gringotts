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
	"fmt"
	"github.com/Akachain/gringotts/glossary"
	"golang.org/x/crypto/sha3"
	"sort"
)

// GenerateID return id of docs base on document prefix name and Fabric transaction ID
// This way, each document in the state database will have this type of ID format:
func GenerateID(docPrefix string, txID string) string {
	h := sha3.New256()
	h.Write([]byte(docPrefix + txID))
	shaString := fmt.Sprintf("%x", h.Sum(nil))

	return shaString[len(shaString)-glossary.IdLength:]
}

// ArrayContains return true if array contain search string, otherwise return fail
func ArrayContains(s []string, searchString string) bool {
	i := sort.SearchStrings(s, searchString)
	return i < len(s) && s[i] == searchString
}
