buildDate = '2024-07-21'
buildVersion = '0.0.1'
secretKey = 'my-secret'

all: build

install:
	./install-protoc.sh
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

protoc:
	protoc --go_out=internal/adapter/handler/proto --go_opt=paths=source_relative \
	--go-grpc_out=internal/adapter/handler/proto --go-grpc_opt=paths=source_relative \
	--proto_path=internal/adapter/handler/proto \
	service.proto user.proto vault.proto

cert:
	openssl req -x509 -newkey rsa:4096 -keyout ./cert/private.pem -out ./cert/cert.pem -passout pass:MPm!uaMf63wN*f9 -sha256 -days 365
	openssl rsa -in ./cert/private.pem -out ./cert/private.pem -passin pass:MPm!uaMf63wN*f9

build: build_server build_client

build_server: protoc cert
	pushd ./cmd/server
	go build -o ../../bin/server -ldflags='-X main.buildDate=$(buildDate) -X main.buildVersion=$(buildVersion)'
	popd

build_client: protoc
	pushd ./cmd/client
	GOOS=darwin GOARCH=arm64 go build -o ../../bin/gophkeeper-osx-arm64 -ldflags='-X main.buildDate=$(buildDate) -X main.buildVersion=$(buildVersion) -X main.secretKey=$('secretKey')'
	GOOS=linux GOARCH=amd64 go build -o ../../bin/gophkeeper-linux-amd64 -ldflags='-X main.buildDate=$(buildDate) -X main.buildVersion=$(buildVersion) -X main.secretKey=$('secretKey')'
	GOOS=windows GOARCH=amd64 go build -o ../../bin/gophkeeper-windows-amd64.exe -ldflags='-X main.buildDate=$(buildDate) -X main.buildVersion=$(buildVersion) -X main.secretKey=$('secretKey')'
	popd

run_server: build_server
	./bin/server

run_client: build_client
	./bin/gophkeeper