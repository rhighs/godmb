CONTAINER:=godmb
VERSION:=latest

docker-build:
	docker build . -t $(CONTAINER):$(VERSION)

docker-run:
	docker run $(CONTAINER):$(VERSION)
