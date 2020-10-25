REPO=naman.io
NAMESPACE=namantest
NAME=virtualbox
BINARY=terraform-provider-${NAME}
VERSION=0.0.1
OS_ARCH=darwin_amd64

build:
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${REPO}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${REPO}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
