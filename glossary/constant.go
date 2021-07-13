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

package glossary

// default select pagination size for query using
var PaginationSize = int32(10)

// default wallet of system using for mint or burn token
var SystemWallet = "0000000000000000000000000000000000000000"

// default length of id from blockchain generate
var IdLength = 40

// The base unit conversion rate for tokens.
//
// In ETH, the base is 10^18
//
// In BTC, the base is 10^8
//
// By default, we take 10^8 as the base conversion rate similar to BTC
const AkcBase = 100000000

// The numeral system where our bigInt string represents.
// By default we use decimal (base 10)
const StringUnitBase = 10

// NftExchange work around to distinguish transfer token using for exchange nft
const NftExchange = "NftExchange"

// NumberWorker use to bulk put state
const NumberWorker = 20
