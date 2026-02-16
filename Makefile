MKDIR_P=mkdir -p
RM=rm
INSTALL=install
INSTALL_PROGRAM=$(INSTALL) -m 0755
INSTALL_DOC=$(INSTALL) -m 0644
PREFIX=/usr
BINDIR=$(PREFIX)/bin
MANDIR=$(PREFIX)/man

PROGRAM=jwt-agent-core

MAN=jwt-agent.1

all: $(PROGRAM) $(MAN)

$(PROGRAM): $(PROGRAM).go
	go build $(PROGRAM).go

install:
	@$(MKDIR_P) $(DESTDIR)$(BINDIR)
	$(INSTALL_PROGRAM) $(PROGRAM) $(DESTDIR)$(BINDIR)/$(PROGRAM)
	$(INSTALL_PROGRAM) jwt-agent $(DESTDIR)$(BINDIR)/jwt-agent
	@$(MKDIR_P) $(DESTDIR)$(MANDIR)/man1
	$(INSTALL_DOC) $(MAN) $(DESTDIR)$(MANDIR)/man1

clean:
	$(RM) -f $(PROGRAM) go.sum

jwt-agent.1: jwt-agent.1.md
	-pandoc -s -t man $< -o $@
