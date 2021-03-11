###
# Terraform Provider for CloudBolt Makefile
###

TF_EXAMPLES=$(shell find examples -not -path '*/\.*' -regex 'examples/.*' -type 'd')
PLUGIN_EXECUTABLE=terraform-provider-cloudbolt
ifeq ($(OS),Windows_NT)
	PLUGIN_RELEASE_EXECUTABLE := $(strip $(PLUGIN_EXECUTABLE)_v$(VERSION)).exe
else
	PLUGIN_RELEASE_EXECUTABLE := $(strip $(PLUGIN_EXECUTABLE)_v$(VERSION))
endif
PKGNAME := cloudbolt
VERSION_NUM := 1.0.0
HOSTOS := $$(go env GOHOSTOS)
HOSTARCH := $$(go env GOHOSTARCH)
PLUGIN_RELEASE_EXECUTABLE := $(PLUGIN_EXECUTABLE)_v$(VERSION_NUM)
TF_PLUGINS_DIR := $(HOME)/.terraform.d/plugins/registry.terraform.io/CloudBoltSoftware/$(PKGNAME)/$(VERSION_NUM)/$(HOSTOS)_$(HOSTARCH)

.PHONY : install
install: build
	mkdir -p $(TF_PLUGINS_DIR)
	mv -f $(PLUGIN_RELEASE_EXECUTABLE) $(TF_PLUGINS_DIR)

.PHONY : build
build:
	go build -o $(PLUGIN_RELEASE_EXECUTABLE)

.PHONY : test
test: $(TF_EXAMPLES)

.PHONY : $(TF_EXAMPLES)
$(TF_EXAMPLES): build
	$(MAKE) -C $@

.PHONY : clean
clean:
	go clean
	rm $(TF_PLUGINS_DIR)/$(PLUGINS_EXECUTABLE)


