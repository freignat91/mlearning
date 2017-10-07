.PHONY: all clean install fmt check version build run test

SHELL := /bin/sh
BASEDIR := $(shell echo $${PWD})

# build variables (provided to binaries by linker LDFLAGS below)
VERSION := 0.0.1

LDFLAGS=-ldflags "-X=main.Version=$(VERSION)"

# ignore vendor directory for go files
SRC := $(shell find . -type f -name '*.go' -not -path './vendor/*' -not -path './.git/*')

# for walking directory tree (like for proto rule)
DIRS = $(shell find . -type d -not -path '.' -not -path './vendor' -not -path './vendor/*' -not -path './.git' -not -path './.git/*')

# generated files that can be cleaned
GENERATED := $(shell find . -type f -name '*.pb.go' -not -path './vendor/*' -not -path './.git/*')

# ignore generated files when formatting/linting/vetting
CHECKSRC := $(shell find . -type f -name '*.go' -not -name '*.pb.go' -not -path './vendor/*' -not -path './.git/*')

OWNER := freignat91
NAME :=  mlearning
TAG := latest

IMAGE := $(OWNER)/$(NAME):$(TAG)
IMAGETEST := $(OWNER)/$(NAME):test
REPO := github.com/$(OWNER)/$(NAME)

CLI := ml
ENGINE := mlserver
TESTS := tests

all: version check install

version:
	@echo "version: $(VERSION) (build: $(BUILD))"

clean:
	@rm -rf $(GENERATED)

install-cli:
	@go install $(LDFLAGS) $(REPO)/$(CLI)

install-engine:
	@go install $(LDFLAGS) $(REPO)/$(ENGINE)

install: install-cli install-engine

# format and simplify if possible (https://golang.org/cmd/gofmt/#hdr-The_simplify_command)
fmt:
	@gofmt -s -l -w $(CHECKSRC)

proto:
	@protoc mlserver/server/server.proto --go_out=plugins=grpc:.

check:
	@test -z $(shell gofmt -l ${CHECKSRC} | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d} | sed '/pb\.go/d'; done
	@go tool vet ${CHECKSRC}

run: 	build

test:
	@go test ./tests -v

install-deps:
	@glide install --strip-vcs --strip-vendor --update-vendored

update-deps:
	@glide update --strip-vcs --strip-vendor --update-vendored

build:	install-cli
	@docker build -t $(IMAGE) .

buildtest: install-cli
	@docker build -t $(IMAGETEST) .

start:
	@docker node inspect self > /dev/null 2>&1 || docker swarm inspect > /dev/null 2>&1 || (echo "> Initializing swarm" && docker swarm init --advertise-addr 127.0.0.1)
	@docker network ls | grep aNetwork || (echo "> Creating overlay network 'aNetwork'" && docker network create -d overlay aNetwork)
	@docker service create --network aNetwork --name mlearning \
	--publish 30107:30107 \
	--detach=true \
	--replicas=1 \
	$(IMAGE)


starttest:
	@docker node inspect self > /dev/null 2>&1 || docker swarm inspect > /dev/null 2>&1 || (echo "> Initializing swarm" && docker swarm init --advertise-addr 127.0.0.1)
	@docker network ls | grep aNetwork || (echo "> Creating overlay network 'aNetwork'" && docker network create -d overlay aNetwork)
	@docker service create --network aNetwork --name mlearning \
	--publish 30107:30107 \
	--detach=true \
	--replicas=1 \
	$(IMAGETEST)

stop:
	@docker service rm mlearning || true

init:
	@docker service rm mlearning || true
	@rm -f ./logs/*
