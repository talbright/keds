PACKAGES=./utils/config ./utils/token ./utils/system ./server ./events

keds: gen plugins
	go build -a -o keds main.go

plugins: example-plugin

example-plugin: plugin/builtin/example/main.go
	go build -a -o plugin/builtin/example/example plugin/builtin/example/*.go

test: keds
	ginkgo -v -race $(PACKAGES)

ginkgo-watch:
	ginkgo watch -r $(PACKAGES)

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
		proto/keds.proto

gen-grpc-swagger: gen-grpc
	cd gen && protoc \
		-I /usr/local/include \
		-I. \
		-I $(GOPATH)/src \
		-I $(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--swagger_out=logtostderr=true:. \
		proto/keds.proto

.PHONY: keds
