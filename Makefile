BIN:=godmb
CONTAINER:=godmb
VERSION:=latest

.PHONY: build
build:
	go mod download
	go mod tidy
	CGO_ENABLED=0 go build -v -ldflags='-s -w' -o $(BIN) .

.PHONY: run
run:
	go build -o ${BIN} . && ./${BIN}

.PHONY: docker-build
docker-build:
	@docker build -t $(CONTAINER):$(VERSION) .

.PHONY: docker-run
docker-run: docker-build
	@docker run --rm -d --name $(CONTAINER) $(CONTAINER):$(VERSION) 

.PHONY: fmt
fmt:
	@which gofmt &>/dev/null || gofmt -w -s  .

.PHONY: lint 
lint:
	@golangci-lint --version || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest \
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

.PHONY: install-ytdlp
install-ytdlp:
	mkdir -p .bin && curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o .bin/yt-dlp\
		&& chmod a+rx .bin/yt-dlp

.PHONY: install-ytdlp-win
install-ytdlp-win:
	mkdir -p .bin && curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe -o .bin/yt-dlp.exe\
		&& chmod a+rx .bin/yt-dlp.exe
