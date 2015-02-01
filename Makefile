VERSION=1.0.2

all: deps
	go build .

deps:
	go get -u github.com/gorilla/pat
	go get -u github.com/crowdmob/goamz/aws

dist: deps
	-rm -rf ./.dist-build
	gox -osarch="linux/amd64" -output="build/imagestore/{{.Dir}}_{{.OS}}_{{.Arch}}" .
	mkdir ./.dist-build
	cp ./build/imagestore/imagestore_linux_amd64 ./.dist-build/imagestore
	cd ./.dist-build; zip imagestore.zip imagestore
	mv ./.dist-build/imagestore.zip imagestore-$(VERSION).zip
	rm -rf ./.dist-build

.PHONY: all dist deps
