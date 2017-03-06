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

test-ci: tools test

test: build validate godo-test-prepare
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

tools: $(GOTOOLS) godo-rebuild

validate: fmt vet

godo-rebuild: $(GODO)
	godo --rebuild

godo-test-prepare: godo-rebuild
	BIRDHOUSE_ENV=test godo db:testPrepare

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

clean: clean-logs clean-zk
	-rm -f ./main
	-rm -f ./keds
	-rm -f $(GOPATH)/bin/keds
	-rm -rf ./vendor

.PHONY: all build test-ci test vet lint fmt tools validate godo-rebuild \
	godo-test-prepare $(GOTOOLS) $(TEST_TARGETS) run glide-reset glide-install \
	clean clean-zk clean-logs start-containers stop-containers