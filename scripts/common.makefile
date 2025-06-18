# common.makefile

# Updated: <2025/06/05 12:19:43>
.PHONY: scripts/common.makefile

VERSION := $(shell ! [ -f VERSION ] && printf "0.0.1" >VERSION; cat VERSION)

help: ## show help
	@echo "The following build targets have help summaries:"
	@gawk 'BEGIN{FS=":.*[#][#]"} /[#][#]/ && !/^#/ {h[$$1":"]=$$2}END{n=asorti(h,d);for (i=1;i<=n;i++){printf "%-26s%s\n", d[i], h[d[i]]}}' $(MAKEFILE_LIST)
	@echo

versionbump:  VERSION ## increment BUILD number in VERSION file
	echo "$(VERSION)" | gawk -F. '{printf "%d.%d.%d",$$1, $$2, $$3+1}' > VERSION

versionbumpminor: VERSION ## increment MINOR number in VERSION file
	echo "$(VERSION)" | gawk -F. '{printf "%d.%d.0", $$1, $$2+1}' > VERSION

versionbumpmajor: VERSION ## increment MAJOR number in VERSION file
	echo "$(VERSION)" | gawk -F. '{printf "%d.0.0", $$1+1}' > VERSION

.PHONY: help versionbump minorversionbump majorversionbump

VERSION:
	@printf "0.0.1" >VERSION

vardump::  ## echo make variables
	@echo "common.makefile: VERSION: $(VERSION)"

clean:: ## clean ephemeral build resources

realclean:: clean  ## clean all resources that can be re-made (implies clean)

distclean:: realclean  ## clean all resources that can be installed or re-made (implies realclean)

.PHONY: vardump clean realclean distclean
