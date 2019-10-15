###
# Terraform Provider for CloudBolt Makefile
###

TF_EXAMPLES=$(shell find examples -not -path '*/\.*' -regex 'examples/.*' -type 'd')
PLUGIN_EXECUTABLE=terraform-provider-cloudbolt
TF_PLUGINS_DIR=$(HOME)/.terraform.d/plugins/

.PHONY : install
install: build
	mkdir -p $(TF_PLUGINS_DIR)
	mv -f $(PLUGIN_EXECUTABLE) $(TF_PLUGINS_DIR)

.PHONY : build
build:
	go build -o $(PLUGIN_EXECUTABLE)

.PHONY : test
test: $(TF_EXAMPLES)

.PHONY : $(TF_EXAMPLES)
$(TF_EXAMPLES): build
	$(MAKE) -C $@

.PHONY : clean
clean:
	go clean
	rm $(TF_PLUGINS_DIR)/$(PLUGINS_EXECUTABLE)


