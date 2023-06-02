CONTAINER:=godmb
VERSION:=latest

fmt:
	gofmt -w -s  .

docker-build:
	docker build . -t $(CONTAINER):$(VERSION)

docker-run:
	docker run $(CONTAINER):$(VERSION)

