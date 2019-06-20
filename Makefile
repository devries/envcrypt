BINARIES := pgdecrypt pgencrypt
PREFIX := /usr/local
SOURCES := go.mod go.sum util.go $(patsubst %,cmd/%/*.go,$(BINARIES))
TESTS := util_test.go
BUILDS := $(patsubst %,build/%,$(BINARIES))
INSTALLS := $(patsubst %,$(PREFIX)/bin/%,$(BINARIES))
UNINSTALLS := $(INSTALLS)
TAREXTRAS := README.md LICENSE Makefile 
TARFILE := envcrypt.tar.gz

bin: $(BUILDS)

$(BUILDS): $(SOURCES)
	mkdir -p build
	go build -o $@ $(patsubst build/%,./cmd/%,$@)

test: $(SOURCES) $(TESTS)
	go test

install: $(INSTALLS)

$(INSTALLS): $(PREFIX)/bin/%: build/%
	install -d $(PREFIX)/bin
	install -m 755 $< $@

clean:
	rm -rf build || true
	rm $(TARFILE)
	go clean

uninstall:
	rm $(INSTALLS) || true

tar: $(TARFILE)
	
$(TARFILE): $(SOURCES) $(TESTS) $(TAREXTRAS)
	tar cvfz $@ $^
	
