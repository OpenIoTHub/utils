module github.com/mDNSService/utils

go 1.12

require (
	github.com/grandcat/zeroconf v0.0.0-20190424104450-85eadb44205c
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
)

replace (
	golang.org/x/net => github.com/golang/net master
    golang.org/x/sync => github.com/golang/sync master
    golang.org/x/sys => github.com/golang/sys master
    golang.org/x/tools => github.com/golang/tools master
    golang.org/x/crypto => github.com/golang/crypto master
    golang.org/x/text => github.com/golang/text master
)
