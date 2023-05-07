test:
	go vet
	go test

build:
	GOOS=linux CGO_ENABLED=0 go build -ldflags -w

update:
	go get -u
	go mod tidy