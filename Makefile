MAKEFLAGS += --no-print-directory
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
HOSTGO := env -u GOOS -u GOARCH -u GOARM -- go
LDFLAGS := $(LDFLAGS)
commit ?= $(BUILD_VERSION)
# Go built-in race detector works only for 64 bits architectures.
ifneq ($(GOARCH), 386)
	race_detector := -race
endif

.PHONY: build
build: tidy
	$(HOSTGO) build -o .build/numa_exporter main.go

.PHONY: go-install
go-install:
	go install -mod=mod -ldflags "-w -s $(LDFLAGS)" -o .build/numa_exporter main.go -trimpath

.PHONY: tidy
tidy:
	go mod verify
	go mod tidy

.PHONY: check
check: fmtcheck vet

.PHONY: test-all
test-all: fmtcheck vet
	go test $(race_detector) ./...

.PHONY: clean
clean:
	rm -rf .build

.PHONY: nametest
nametest:
	@echo $(commit)


.PHONY: docker-image
docker-image: build
	docker build -f scripts/Dockerfile -t numa_exporter:$(commit) .