test:
	go test -v $(shell go list ./... | grep -v /vendor/)

build: clean
	gox -osarch="linux/amd64 windows/amd64 darwin/amd64" \
	-output="releases/{{.OS}}_{{.Arch}}/tf-inventory" .

#release:

deps:
	go get -u github.com/mitchellh/gox

clean:
	rm -rf releases/
