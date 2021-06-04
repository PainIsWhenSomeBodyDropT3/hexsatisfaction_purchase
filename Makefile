.PHONY lint:
	golangci-lint run --config .golangci.yml

.PHONY swagger:swagger-spec

swagger-spec:
	swag init -g cmd/main.go

gen-mocks:
	mockery --all --keeptree
run:
	go run cmd/main.go

test-coverage:
	go test ./... -coverprofile coverage.out
	go tool cover -func coverage.out
	go tool cover -html=coverage.out -o coverage.html

start : lint  swagger  run
