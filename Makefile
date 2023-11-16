BINARY_NAME=$(shell basename "$(PWD)")
DEMO_DIR=./example

.PHONY: build gen

build:
	@go build -o bin/$(BINARY_NAME) .

gen:
	@go run . -gen $(DEMO_DIR)/out -pkg out -proto-gen-path $(DEMO_DIR)/gen/proto -type User -model User -O "suppress=Permissions" -O "no-embed=Friends" -slice
	@go run . -gen $(DEMO_DIR)/out -pkg out -proto-gen-path $(DEMO_DIR)/gen/proto -type Friend -model Friend

install:
	@go install .

clean:
	@rm -rf $(DEMO_DIR)/out/*.entity.go
