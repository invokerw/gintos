GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
	API_PROTO_FILES=$(shell find api -name *.proto)
endif

.PHONY: init
# init env
init:
	go install github.com/google/wire/cmd/wire@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/favadi/protoc-go-inject-tag@latest
	cd ../cmd/protoc-gen-go-gin/; go install .; cd -


.PHONY: api
# generate api proto
api:
	protoc --experimental_allow_proto3_optional \
		   --proto_path=./api \
	       --proto_path=../third_party \
 	       --go_out=paths=source_relative:./api \
 	       --go-gin_out=paths=source_relative,rbac_path=./api,rbac_package_name=api:./api \
	       --openapi_out=fq_schema_naming=true,naming=proto,default_response=false:./assets \
	       $(API_PROTO_FILES)
	PB_GO_FILES=$$(find api -name '*.pb.go' | grep -v gin.pb.go); \
		for file in $$PB_GO_FILES; do \
			echo "Injecting tags into: $$file"; \
			protoc-go-inject-tag -input="$${file}"; \
		done

conf:
	protoc --proto_path=./internal \
		   --proto_path=../third_party \
		   --go_out=paths=source_relative:./internal \
		   $(INTERNAL_PROTO_FILES)

.PHONY: wire
# generate wire
wire:
	go run github.com/google/wire/cmd/wire

.PHONY: db
# migrate db
db:
	cd ./internal/data/ent && go generate ./...

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: all
# generate all
all:
	make conf;
	make db;
	make wire;
	make api;
	make build;

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
