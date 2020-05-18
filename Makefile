PKG     = github.com/duskwuff/pogo
TOOLS   = data2json datasniffer ggpk

default: $(TOOLS:%=bin/%)

clean:
	-rm -rf bin

bin/%:
	go build -o $@ ./cmd/$*
