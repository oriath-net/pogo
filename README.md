Path of Go - Golang tools for reading PoE data files
====================================================

Installation
------------

    go get -u github.com/oriath.net/pogo

You'll need to have [Go 1.14](https://golang.org/dl/) or later installed.


Usage
-----

    pogo ggpk list Content.ggpk

List the contents of a GGPK file.


    pogo ggpk extract --no-recurse --into root/ Content.ggpk Data/

Extract the `Data` directory from `Content.ggpk` into the `root` directory,
skipping any subdirectories.

(The `ggpk` tool supports several other useful options; run `ggpk --help` for
details.)


    pogo data2json -f formats/demo.go Content.ggpk:Data/ActiveSkills.dat

Dump a data file (directly from the GGPK!) to JSON. This currently only works
for ActiveSkills (since it's the only format specified in `formats/demo.go`),
but adding other formats to that or another file will make them supported as
well.
