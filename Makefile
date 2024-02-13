build: 
	cd cmd && go build -o hexsatisfaction_purchase .

docker: build
	# docker image prune -af
	docker build -t hexsatisfaction_purchase:1.0 .

docker-compose: docker
	# docker image prune -af
	# docker container prune -f
	docker-compose build --no-cache
	docker-compose up -d








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

