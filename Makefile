goversion=$(word 3,$(shell go version))

all:
	go test -v ./...
	if [ "${goversion}" = "go1.6.2" ]; then \
		go get github.com/golang/lint/golint ; \
		golint ./... ; \
		go test -cover -race ./... ; \
		go vet ./... ; \
	fi
	test -z "$(go fmt -s ./...)"
	go generate ./...
	git diff --exit-code
