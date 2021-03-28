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

// The entity package contains all structured entities that are defined in the
// Gringotts system.
package entity

// The base entity provides necessary fields that every other entities must
// contain.
type Base struct {
	// The document id in state database
	Id string `mapstructure:"id" json:"id"`

	// This is the Fabric transaction id where this entity is created.
	// Fabric allows us to search for specific block if we know the block number
	// or the transaction id in that block. This is very helpful if later we want
	// to query this block to check the Read/Write set for more details
	BlockChainId string `mapstructure:"blockChainId" json:"blockChainId"`

	// Timestamp of the transaction that created this entity
	CreatedAt string `mapstructure:"createdAt" json:"createdAt"`

	// Timestamp of the last transaction that updated this entity
	UpdatedAt string `mapstructure:"updatedAt" json:"updatedAt"`
}
