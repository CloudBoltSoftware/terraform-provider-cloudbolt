TF_EXAMPLES=$(shell find examples -not -path '*/\.*' -regex 'examples/.*' -type 'd')
PLUGIN_EXECUTABLE=terraform-provider-cloudbolt
TF_PLUGINS_DIR=$(HOME)/.terraform.d/plugins/

build:
	go build -o $(PLUGIN_EXECUTABLE)
	mv -f $(PLUGIN_EXECUTABLE) $(TF_PLUGINS_DIR)

test: $(TF_EXAMPLES)

# TODO: Running examples

$(TF_EXAMPLES): build
	$(MAKE) -C $@

clean:
	go clean
	rm $(TF_PLUGINS_DIR)/$(PLUGINS_EXECUTABLE)
