module github.com/Akachain/gringotts

go 1.15

replace golang.org/x/sys => github.com/golang/sys v0.0.0-20210124154548-22da62e12c0c

require (
	github.com/Akachain/akc-go-sdk-v2 v1.0.0
	github.com/davecgh/go-spew v1.1.1
	github.com/hyperledger/fabric v2.1.1+incompatible
	github.com/hyperledger/fabric-chaincode-go v0.0.0-20210319203922-6b661064d4d9
	github.com/hyperledger/fabric-contract-api-go v1.1.1
	github.com/mitchellh/mapstructure v1.3.2
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	gotest.tools v2.2.0+incompatible
)
