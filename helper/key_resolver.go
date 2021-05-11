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

// WalletKey return list key of wallet will be compose in couch db key
func WalletKey(walletId string) []string {
	return []string{walletId}
}

// TransactionKey return list key of transaction will be compose in couch db key
func TransactionKey(txId string) []string {
	return []string{txId}
}

// TokenKey return list key of token will be compose in couch db key
func TokenKey(tokenId string) []string {
	return []string{tokenId}
}

// HealthCheckKey return list key of health check will be compose in couch db key
func HealthCheckKey(id string) []string {
	return []string{id}
}

// EnrollmentKey return list key of enrollment will be compose in couch db key
func EnrollmentKey(tokenId string) []string {
	return []string{tokenId}
}

// NFTKey return list key of NFT will be compose in couch db key
func NFTKey(nftId string) []string {
	return []string{nftId}
}

// ExchangeKey return list key of exchange docs will be compose in couch db key
func ExchangeKey(txId string) []string {
	return []string{txId}
}
