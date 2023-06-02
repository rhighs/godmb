BIN:=godmb
CONTAINER:=godmb
VERSION:=latest

all:
	go build  -o $(BIN) .

.PHONY: fmt lint docker-build docker-run

fmt:
	@ gofmt --version || gofmt -w -s  .

lint:
	@ golangci-lint --version || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest \
		&& golangci-lint run --disable-all \
							  	-E errcheck \
							  	-E gosimple \
					  			-E ineffassign \
					  			-E staticcheck \
					  			-E typecheck \
					  			-E unused \
					  			-E bodyclose \
					  			-E contextcheck \
								-v ./... \

docker-build:
	@ docker build . -t $(CONTAINER):$(VERSION)

docker-run:
	@ docker run $(CONTAINER):$(VERSION)
