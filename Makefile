# SPDX-FileCopyrightText: 2023 Christoph Mewes
# SPDX-License-Identifier: MIT

GO_TEST_FLAGS ?=

# thank you, https://stackoverflow.com/a/62283258
define FOREACH_MODULE
	for dir in $(dir $(wildcard */go.mod)); do \
	  (set -x; cd $$dir && $(1)); \
	done
endef

.DEFAULT: all
all: gimps test spellcheck lint

.PHONY: test
test:
	@$(call FOREACH_MODULE, CGO_ENABLED=1 go test $(GO_TEST_FLAGS) ./...)

.PHONY: gimps
gimps:
	@$(call FOREACH_MODULE, gimps --config ../.gimps.yaml .)

.PHONY: lint
lint:
	@$(call FOREACH_MODULE, golangci-lint run ./...)

.PHONY: spellcheck
spellcheck:
	typos
