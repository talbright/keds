PACKAGES=./utils/config ./utils/token ./utils/system ./server ./events ./client ./cmd
VENDOR=$(shell pwd)/vendor
GOTOOLS := \
		github.com/onsi/ginkgo/ginkgo \
		github.com/golang/protobuf/protoc-gen-go \
		github.com/golang/protobuf/proto

keds: gen plugins
	go build -a -o keds main.go

plugins: example-plugin notifications-plugin

example-plugin: plugin/builtin/example/main.go
	go build -a -o plugin/builtin/example/example plugin/builtin/example/*.go

notifications-plugin: plugin/builtin/notifications/main.go
	go build -a -o plugin/builtin/notifications/notifications plugin/builtin/notifications/*.go

test: keds
	ginkgo -v -race $(PACKAGES)

ginkgo-watch:
	ginkgo watch -r $(PACKAGES)

gen: gen-grpc gen-grpc-reverse-proxy gen-grpc-swagger

gen-grpc:
	cd gen && protoc \
		-I proto \
		-I $(VENDOR)/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:proto \
		proto/keds.proto

gen-grpc-reverse-proxy: gen-grpc
	cd gen && protoc \
		-I /usr/local/include \
		-I. \
		-I $(VENDOR)/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--grpc-gateway_out=logtostderr=true:. \
		proto/keds.proto

gen-grpc-swagger: gen-grpc
	cd gen && protoc \
		-I /usr/local/include \
		-I. \
		-I $(VENDOR)/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--swagger_out=logtostderr=true:. \
		proto/keds.proto

go-tools:
	go get -u -v $(GOTOOLS)

.PHONY: keds go-tools
