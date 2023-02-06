MKDIR_P=mkdir -p
RM=rm
INSTALL=install
INSTALL_PROGRAM=$(INSTALL) -m 0755
PREFIX=/usr
BINDIR=$(PREFIX)/bin

PROGRAM=jwt-agent-core

all: $(PROGRAM)

$(PROGRAM): $(PROGRAM).go
	go mod download golang.org/x/crypto
	go build $(PROGRAM).go

install:
	@$(MKDIR_P) $(DESTDIR)$(BINDIR)
	$(INSTALL_PROGRAM) $(PROGRAM) $(DESTDIR)$(BINDIR)/$(PROGRAM)
	$(INSTALL_PROGRAM) jwt-agent $(DESTDIR)$(BINDIR)/jwt-agent

clean:
	$(RM) -f $(PROGRAM) go.sum
