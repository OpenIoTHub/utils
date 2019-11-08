module github.com/mDNSService/utils

go 1.12

require (
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/iotdevice/zeroconf v0.0.0-20190424104450-85eadb44205c
	github.com/kr/pretty v0.1.0 // indirect
	github.com/miekg/dns v1.1.22 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20191011191535-87dc89f01550
	golang.org/x/net => github.com/golang/net v0.0.0-20191011234655-491137f69257
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys => github.com/golang/sys v0.0.0-20191010194322-b09406accb47
	golang.org/x/text => github.com/golang/text v0.3.2
	golang.org/x/tools => github.com/golang/tools v0.0.0-20191011211836-4c025a95b26e
	golang.org/x/xerrors => github.com/golang/xerrors v0.0.0-20191011141410-1b5146add898
)
