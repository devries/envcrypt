BINARIES := pgdecrypt pgencrypt
PREFIX := /usr/local
SOURCES := go.mod go.sum util.go $(patsubst %,cmd/%/*.go,$(BINARIES))
BUILDS := $(patsubst %,build/%,$(BINARIES))
INSTALLS := $(patsubst %,$(PREFIX)/bin/%,$(BINARIES))

bin: $(BUILDS)

$(BUILDS): $(SOURCES)
	mkdir -p build
	go build -o $@ $(patsubst build/%,./cmd/%,$@)

test:
	go test

install: $(INSTALLS)

$(INSTALLS): $(PREFIX)/bin/%: build/%
	install -d $(PREFIX)/bin
	install -m 755 $< $@

clean:
	rm -rf build || true
	go clean
