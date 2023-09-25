.PHONEY:
all: build

build:
	go build -o dist/file-server-manager ./cmd/manager

.PHONEY:
release: clean build
	cp file-server-manager.service dist/
	tar -zcvf file-server-manager.tar.gz -C dist/ .

.PHONEY:
clean:
	rm -rf dist
	rm -rf *.tar.gz

lint:
	find -name *.go | xargs gofmt -w
