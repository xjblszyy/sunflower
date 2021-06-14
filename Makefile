# workdir info
PACKAGE=sunflower
PREFIX=$(shell pwd)
CMD_PACKAGE=${PACKAGE}
OUTPUT_DIR=${PREFIX}/bin
OUTPUT_FILE=${OUTPUT_DIR}/sunflower
COMMIT_ID=$(shell git rev-parse --short HEAD)
VERSION=$(shell git describe --tags || echo "v0.0.0")
VERSION_IMPORT_PATH=cmd/version
BUILD_TIME=$(shell date '+%Y-%m-%dT%H:%M:%S%Z')
VCS_BRANCH=$(shell git symbolic-ref --short -q HEAD)

# which golint
GOLINT=$(shell which golangci-lint || echo '')

# build args
BUILD_ARGS := \
    -ldflags "-X $(VERSION_IMPORT_PATH).appName=$(PACKAGE) \
    -X $(VERSION_IMPORT_PATH).version=$(VERSION) \
    -X $(VERSION_IMPORT_PATH).revision=$(COMMIT_ID) \
    -X $(VERSION_IMPORT_PATH).branch=$(VCS_BRANCH) \
    -X $(VERSION_IMPORT_PATH).buildDate=$(BUILD_TIME)"
EXTRA_BUILD_ARGS=

export GOCACHE=

.PONY: lint test
default: lint test build

dep:
	go get goa.design/goa/v3 && go get goa.design/goa/v3/... && go get github.com/fatih/gomodifytags

lint:
	@echo "+ $@"
	@$(if $(GOLINT), , \
		$(error Please install golint: `go get -u github.com/golangci/golangci-lint/cmd/golangci-lint`))
	golangci-lint run --deadline=10m ./...

test:
	@echo "+ test"
	go test -cover $(EXTRA_BUILD_ARGS) ./...

build:
	@echo "+ build"
	go build $(BUILD_ARGS) $(EXTRA_BUILD_ARGS) -o ${OUTPUT_FILE} $(CMD_PACKAGE)

clean:
	@echo "+ $@"
	@rm -r "${OUTPUT_DIR}"

GOA_PB_FILES = $(shell find gen -type f -name '*.pb.go')
GOA_MOCK_TAR = $(patsubst %.pb.go, %.mock.go, $(GOA_PB_FILES))

pkg/api/gen/grpc/%.mock.go: pkg/api/gen/grpc/%.pb.go
	@echo "+" $(subst pb,pb/mock, $@)
	@mockgen --source $< --destination $(subst pb,pb/mock, $@) --package=mock

goa-mock: $(GOA_MOCK_TAR)

gen-apidoc:
	@echo "need: npm -g install redoc-cli"
	@redoc-cli bundle pkg/api/gen/http/openapi3.json -o docs/apidoc.html

GOA_DESIGN_PATH=pkg/api/design
GOA_GEN_OUTPUT=pkg/api/
goa:
	echo "+ $@"
ifeq ($(GOA_GEN_OUTPUT), .)
	goa gen $(PACKAGE)/$(GOA_DESIGN_PATH)
else
	goa gen $(PACKAGE)/$(GOA_DESIGN_PATH) -o $(GOA_GEN_OUTPUT)
endif
	@make goa-mock GOA_GEN_OUTPUT=$(GOA_GEN_OUTPUT)
	@make gen-apidoc

release:
	@echo "+ $@"
	@make build
	@mkdir -p dist/
	@tar -zcvf dist/sunflower-${VERSION}.tar.gz README.md -C bin/ sunflower -C ../config/ config.sample.yml sunflower.service
