BINARIES := pgdecrypt pgencrypt
PREFIX := /usr/local
SOURCES := go.mod go.sum util.go $(patsubst %,cmd/%/*.go,$(BINARIES))
TESTS := util_test.go
BUILDS := $(patsubst %,build/%,$(BINARIES))
INSTALLS := $(patsubst %,$(PREFIX)/bin/%,$(BINARIES))
UNINSTALLS := $(INSTALLS)
TAREXTRAS := README.md LICENSE Makefile 
TARFILE := envcrypt.tar.gz
TARDIRECTORY := envcrypt

build: $(BUILDS) ## Build binaries

$(BUILDS): $(SOURCES)
	mkdir -p build
	go build -o $@ $(patsubst build/%,./cmd/%,$@)

test: $(SOURCES) $(TESTS) ## Run tests
	go test

install: $(INSTALLS) ## Install binaries into $(PREFIX)/bin

$(INSTALLS): $(PREFIX)/bin/%: build/%
	install -d $(PREFIX)/bin
	install -m 755 $< $@

clean: ## Remove all local artifacts
	rm -rf build || true
	-rm $(TARFILE)
	go clean

uninstall: ## Uninstall the binaries
	rm $(INSTALLS) || true

tar: $(TARFILE) ## Tar up the source
	
$(TARFILE): $(SOURCES) $(TESTS) $(TAREXTRAS)
	mkdir -p $(TARDIRECTORY)
	tar cf - $^ | tar xf - -C $(TARDIRECTORY)
	tar cvfz $@ $(TARDIRECTORY)
	rm -rf $(TARDIRECTORY)

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build test install clean uninstall tar help
