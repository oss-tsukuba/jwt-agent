MKDIR_P=mkdir -p
RM=rm
INSTALL=install
INSTALL_PROGRAM=$(INSTALL) -m 0755
PREFIX=/usr
BINDIR=$(PREFIX)/bin

PROGRAM=jwt-agent

all: $(PROGRAM)

$(PROGRAM): $(PROGRAM).go
	go mod download golang.org/x/crypto
	go get golang.org/x/crypto/ssh/terminal@v0.0.0-20220722155217-630584e8d5aa
	go get github.com/mattn/go-isatty
	go build $(PROGRAM).go

install:
	@$(MKDIR_P) $(DESTDIR)$(BINDIR)
	$(INSTALL_PROGRAM) $(PROGRAM) $(DESTDIR)$(BINDIR)/$(PROGRAM)

clean:
	$(RM) -f $(PROGRAM) go.sum
