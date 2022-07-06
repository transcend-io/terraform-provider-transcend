default: testacc
HOSTNAME=transcend.com
NAMESPACE=cli
NAME=transcend-io
BINARY=terraform-provider-${NAME}
VERSION=0.1
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

build:
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/$(GOOS)_$(GOARCH)
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/$(GOOS)_$(GOARCH)
