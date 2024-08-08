all: run

run:
	go run ./cmd/... -config=dev.yml

build:
	go build -o automated-stratus-red-team ./cmd/...

test:
	go test -v ./...

clean:
	rm -r dist/ automated-stratus-red-team || true

update:
	go get -u ./...
	go mod tidy
