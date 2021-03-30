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

// Package query contains all query string used throughout the project.
// While it may make sense if we put this part inside each business service
// package, we find that multiple services may re-use the same query. Thus,
// we put it here so that the same query can be re-used among all services if needed.
package query

// GetPendingTransactionQueryString return query list transaction have status pending.
// Based on the state db document prefix, we basically get all documents that has prefix
// greater \u0000Transactions and less than \u0000Transaction + inf as well as pending status
// It is then sorted ascending based on the created timestamp
//
// The corresponding index is defined in the META-INF directory
func GetPendingTransactionQueryString() string {
	return `
		{ "selector": 
			{ 	
				"Status": 
					{ "$eq": "Pending" },
				"_id": 
					{"$gt": "\u0000Transactions",
					"$lt": "\u0000Transactions\uFFFF"}			
			},
			"sort": [
				"CreatedAt"
			],
			"use_index":["indexPendingTxDoc","indexPendingTx"]
		}`
}
