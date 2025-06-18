# openapi.makefile

# Updated: <2025/06/04 17:38:20>

.PHONY: scripts/openapi.makefile 
include scripts/go.makefile
include scripts/npm.makefile

ifndef BINDIR
$(error BINDIR is not set)
endif

#####
# Redocly - for linting and generating docs
# REDOCLY_CLI := $(NPM) exec -- @redocly/cli
REDOCLY_CLI := $(NPMBINRELDIR)/redocly

.PHONY: redocly-cli
redocly-cli: $(REDOCLY_CLI)  ## install redocly-cli - for linting and generating docs

$(REDOCLY_CLI): | npm-installed
	$(NPM) list @redocly/cli >/dev/null || $(NPM) install @redocly/cli

redocly-cli-uninstall: | npm-installed  ## uninstall redocly-cli
	$(NPM) uninstall @redocly/cli

#####
# Code generation
OAPI_CODEGEN_INSTALL_VERSION := v2.4.1
OAPI_CODEGEN := $(shell command -v $(BINDIR)/oapi-codegen 2> /dev/null)

.PHONY: oapi-codegen-installed
oapi-codegen-installed: ## check if oapi-codegen tool is installed
ifndef OAPI_CODEGEN
	$(error "OAPI_CODEGEN is not installed; try 'make oapi-codegen'")
endif

.PHONY: oapi-codegen
oapi-codegen: $(BINDIR)/oapi-codegen ## install oapi-codegen tool

$(BINDIR)/oapi-codegen: | go-installed
	GOBIN=$(BINDIR) $(GO) install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_INSTALL_VERSION)

.PHONY: vardump distclean
distclean::
	rm -f $(BINDIR)/oapi-codegen
	$(NPM) uninstall @redocly/cli@latest

vardump::
	@echo "openapi.makefile: REDOCLY_CLI: $(REDOCLY_CLI)"
	@echo "openapi.makefile: REDOCLY_CLI version: $(shell $(REDOCLY_CLI) --version)"
	@echo "openapi.makefile: OAPI_CODEGEN: $(OAPI_CODEGEN)"
	@echo "openapi.makefile: OAPI_CODEGEN version: $(shell $(OAPI_CODEGEN) --version | gawk '/^v/')"
	@echo "openapi.makefile: OAPI_CODEGEN_INSTALL_VERSION: $(OAPI_CODEGEN_INSTALL_VERSION) -- override this to change the version of oapi-codegen installed"
