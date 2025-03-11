TEST?=./...
PKG_NAME=openvswitch

default: build

build: fmtcheck
	go install

test: fmtcheck
	go test $(TEST) -parallel=4

testacc:
	TF_ACC=1 go test ./$(PKG_NAME) -v $(TESTARGS) -timeout 120m

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./$(PKG_NAME)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

lint:
	@echo "==> Checking source code against linters..."
	@golangci-lint run ./$(PKG_NAME)
	@tfproviderlint \
		-c 1 \
		./$(PKG_NAME)

tools:
	@echo "==> installing required tooling..."
	GO111MODULE=on go install github.com/bflad/tfproviderlint/cmd/tfproviderlint
	GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: build fmt fmtcheck lint test testacc tools vet
