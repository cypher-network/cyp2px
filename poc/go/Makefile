NAME   := cypher/cyp2px
TAG    := $(shell git log -1 --pretty=%H)
IMG    := ${NAME}:${TAG}
LATEST := ${NAME}:latest

define print
	@echo ""
	@echo "### $(1)"
	@echo ""
endef

image:
	$(call print,getting revision)
	$(eval REVISION=$(shell git log -1 --pretty=%h))
	@echo "Revision: ${REVISION}"
	
	$(call print,building docker image with tag ${IMG})
	@docker build --tag ${IMG} .
	
	$(call print,getting images: ${NAME})
	@docker images ${NAME}
	
	@echo ""
