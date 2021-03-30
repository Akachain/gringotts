module github.com/Akachain/gringotts

go 1.15

replace golang.org/x/sys => github.com/golang/sys v0.0.0-20210124154548-22da62e12c0c

replace github.com/satori/go.uuid v1.2.0 => github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b

require (
	github.com/Akachain/akc-go-sdk v1.0.12
	github.com/davecgh/go-spew v1.1.1
	github.com/golang/protobuf v1.3.3
	github.com/hyperledger/fabric v2.1.1+incompatible
	github.com/hyperledger/fabric-chaincode-go v0.0.0-20201119163726-f8ef75b17719
	github.com/hyperledger/fabric-contract-api-go v1.1.1
	github.com/mitchellh/mapstructure v1.1.2
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.4.2
	github.com/sykesm/zap-logfmt v0.0.4 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	gotest.tools v2.2.0+incompatible
)
