###
# Terraform Provider for CloudBolt Makefile
###

TF_EXAMPLES=$(shell find examples -not -path '*/\.*' -regex 'examples/.*' -type 'd')
HOSTNAME=registry.terraform.io
NAMESPACE=CloudBoltSoftware
NAME=cloudbolt
VERSION=1.0.0
BINARY=terraform-provider-${NAME}
TF_PLUGINS_DIR=$(HOME)/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
OS_ARCH=darwin_amd64

.PHONY : install
install: build
	mkdir -p $(TF_PLUGINS_DIR)
	mv -f $(BINARY) $(TF_PLUGINS_DIR)

.PHONY : build
build:
	go build -o $(BINARY)

.PHONY : test
test: $(TF_EXAMPLES)

.PHONY : $(TF_EXAMPLES)
$(TF_EXAMPLES): build
	$(MAKE) -C $@

.PHONY : clean
clean:
	go clean
	rm $(TF_PLUGINS_DIR)/$(BINARY)


