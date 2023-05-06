SUBDIRS := $(wildcard */.)

test: $(SUBDIRS)

build: $(SUBDIRS)

update: $(SUBDIRS)

deploy: $(SUBDIRS)

undeploy: $(SUBDIRS)

$(SUBDIRS):
	$(MAKE) -C $@ $(MAKECMDGOALS)

.PHONY: test build deploy $(SUBDIRS)
