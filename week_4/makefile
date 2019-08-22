test:
	go test -v ./...

install:
	go get -u google.golang.org/grpc
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go get github.com/pact-foundation/pact-go

set-up:
	export PATH=$(PATH):$(GOPATH)/bin

generate:
	go generate ./practice/todo/

run:
	docker-compose down
	docker-compose up -d
	go run ./practice/todo/server/main.go

set-up-test-env:
	docker-compose down && docker-compose up -d

clean:
	docker-compose down

integration-test: set-up-test-env
	 go test -count=1 -v -tags=integration ./practice/todo/server...

integration-test-persistence: set-up-test-env
	 go test -count=1 -v -tags=persistence ./practice/todo/server...

integration-test-contract:
	 go test -count=1 -v -tags=contract ./practice/todo/server...

integration-test-consumer:
	 go test -count=1 -v -tags=integration ../practice/todo/consumer/...
