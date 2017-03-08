CGOENABLED:=1
PACKAGES=$(shell glide novendor)
PACKAGE_DIRS=$(shell for i in $$(glide novendor -x | egrep -v '^\.$$' ); do echo $${i%/}; done)
TEST_TARGETS = $(addprefix "test-",$(notdir $(PACKAGE_DIRS)))
GOLINT = github.com/golang/lint/golint
GODO = gopkg.in/godo.v1/cmd/godo
GOTOOLS = $(GODO) github.com/onsi/ginkgo/ginkgo github.com/GeertJohan/fgt
CWD = $(shell pwd)

define DOCKER_ENV
DATA_VOLUMES=$(CWD)/extras/docker/data
endef

export CGOENABLED

all: tools validate build

keds: main.go
	go build -a -o keds main.go

build: glide-install compile

test: build validate
	go test -p 1 -v -race $(PACKAGES)

vet:
	fgt go vet $(PACKAGES)

lint: $(GOLINT)
	@for p in $(PACKAGES) ; do \
		golint $$p ; \
	done

fmt:
	fgt gofmt -l $(PACKAGE_DIRS) *.go

$(GOLINT):
	go install ./vendor/$@

$(GOTOOLS):
	go install ./vendor/$@

tools: $(GOTOOLS)

gen-minipod: gen-minipod-grpc gen-minipod-reverse-proxy gen-minpod-swagger

gen-minipod-grpc:
	cd gen && protoc \
		-I minipod/ \
		-I $(GOPATH)/src \
		-I $(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:minipod \
		minipod/minipod.proto

gen-minipod-reverse-proxy:
	cd gen && protoc \
		-I /usr/local/include \
		-I. \
		-I $(GOPATH)/src \
		-I $(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--grpc-gateway_out=logtostderr=true:. \
		minipod/minipod.proto

gen-minipod-swagger:
	cd gen && protoc \
		-I /usr/local/include \
		-I. \
		-I $(GOPATH)/src \
		-I $(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--swagger_out=logtostderr=true:. \
		minipod/minipod.proto

minipod-server:
	cd minipod/server && go run main.go

minipod-client:
	cd minipod/client && go run main.go

gen: gen-minipod

validate: fmt vet

test-%: PCKG = ./$*/...
test-%: PCKG_DIR = ./$*
test-%:
	fgt gofmt -l $(PCKG_DIR)
	fgt go vet $(PCKG)
	ginkgo -v -r $(PCKG_DIR)

glide-install:
	glide install

glide-reset:
	-rm -rf ./vendor
	glide install

run:
	-go run -race main.go | tee run.log

clean-logs:
	-rm *.log

clean: clean-logs
	-rm -f ./main
	-rm -f ./keds
	-rm -f $(GOPATH)/bin/keds
	-rm -rf ./vendor

.PHONY: all build test-ci test vet lint fmt tools validate godo-rebuild \
	godo-test-prepare $(GOTOOLS) $(TEST_TARGETS) run glide-reset glide-install \
	clean clean-zk clean-logs start-containers stop-containers
