# Makefile  -*-Makefile-*-

BINDIR := $(PWD)/bin

.PHONY: Makefile
include scripts/common.makefile
include scripts/go.makefile
include scripts/openapi.makefile
export

OPENAPI_SPECS_FILES := secrets-hub-api.yaml
OPENAPI_SPECS := $(addprefix api/,$(OPENAPI_SPECS_FILES))

.PHONY: docs
docs: docs/secrets-hub-api.html

docs/secrets-hub-api.html: api/secrets-hub-api.yaml
	$(REDOCLY_CLI) build-docs public@v0 -o $@

.PHONY: gen
gen: VERSION secretshub/secretshub-client.gen.go secretshub/secretshub-types.gen.go

secretshub/secretshub-types.gen.go: api/secrets-hub-api.yaml | oapi-codegen-installed
	$(OAPI_CODEGEN) -generate types -package secretshub $< > $@

secretshub/secretshub-client.gen.go: api/secrets-hub-api.yaml | oapi-codegen-installed
	$(OAPI_CODEGEN) -generate client -package secretshub $< > $@

$(BINDIR)/secretshub-client: VERSION secretshub/secretshub-client.gen.go secretshub/secretshub-types.gen.go examples/secretshub-client/main.go
	$(GO) build -o $@ $(LDFLAGS) examples/secretshub-client/main.go

.PHONY: get-all-stores
get-all-stores: $(BINDIR)/get-all-stores  ## Build the get-all-stores binary

$(BINDIR)/get-all-stores: VERSION examples/get-all-stores/main.go secretshub/secretshub-client.gen.go secretshub/secretshub-types.gen.go
	$(GO) build -o $@ $(LDFLAGS) examples/get-all-stores/main.go

.PHONY: get-secrets
get-secrets: $(BINDIR)/get-secrets  ## Build the get-secrets binary

$(BINDIR)/get-secrets: VERSION examples/get-secrets/main.go secretshub/secretshub-client.gen.go secretshub/secretshub-types.gen.go
	$(GO) build -o $@ $(LDFLAGS) examples/get-secrets/main.go

.PHONY: get-policies
get-policies: $(BINDIR)/get-policies  ## Build the get-policies binary

$(BINDIR)/get-policies: VERSION examples/get-policies/main.go secretshub/secretshub-client.gen.go secretshub/secretshub-types.gen.go
	$(GO) build -o $@ $(LDFLAGS) examples/get-policies/main.go


.PHONY: secretshub-client
secretshub-client: $(BINDIR)/secretshub-client ## Build the secretshub-client binary

.PHONY: build-all
build-all: gen docs secretshub-client get-all-stores get-secrets get-policies ## Build all binaries

.PHONY: deps vardump clean realclean
deps: oapi-codegen redocly-cli identity-client ## Install dependencies

clean::
	rm -f $(BINDIR)/secretshub-client
	rm -f $(BINDIR)/get-all-stores
	rm -f $(BINDIR)/get-secrets
	rm -f $(BINDIR)/get-policies

realclean:: clean
	rm -f secretshub/*.gen.go
	rm -f docs/secrets-hub-api.html

distclean:: realclean
	rm -f $(BINDIR)/identity-client
	rm -f $(BINDIR)/oapi-codegen

vardump::
	@echo "Makefile: BINDIR: $(BINDIR)"
	@echo "Makefile: OPENAPI_SPECS: $(OPENAPI_SPECS)"

# Helpers
.PHONY: identity-client
identity-client: $(BINDIR)/identity-client ## Install the identity-client binary from the identityadmin-sdk-go examples

$(BINDIR)/identity-client:
	GOBIN=$(BINDIR) $(GO) install github.com/davidh-cyberark/identityadmin-sdk-go/examples/identity-client@latest

