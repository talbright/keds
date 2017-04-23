PACKAGES=./utils/config ./utils/token

keds:
	go build -a -o keds main.go

plugins: plugin-main

plugin-main: plugin/example/main.go
	go build -a -o plugin/example/example plugin/example/*.go

test: keds
	go test -p 1 -v -race $(PACKAGES)

gen: gen-grpc gen-grpc-reverse-proxy gen-grpc-swagger

gen-grpc:
	cd gen && protoc \
		-I proto \
		-I $(GOPATH)/src \
		-I $(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:proto \
		proto/keds.proto

gen-grpc-reverse-proxy: gen-grpc
	cd gen && protoc \
		-I /usr/local/include \
		-I. \
		-I $(GOPATH)/src \
		-I $(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--grpc-gateway_out=logtostderr=true:. \
		gen/proto/keds.proto

gen-grpc-swagger: gen-grpc
	cd gen && protoc \
		-I /usr/local/include \
		-I. \
		-I $(GOPATH)/src \
		-I $(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--swagger_out=logtostderr=true:. \
		gen/proto/keds.proto

.PHONY: keds

