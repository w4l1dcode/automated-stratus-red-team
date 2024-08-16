all: run

run:
	# Run program
	go run ./cmd/... -config=config.yml

build:
	aws eks update-kubeconfig --region $(K8S_REGION) --name $(CLUSTER_NAME)
	go build -o stratus-red-team ./cmd/...

test:
	go test -v ./...

clean:
	rm -r dist/ stratus-red-team || true

update:
	go get -u ./...
	go mod tidy