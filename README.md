Path of Go - Golang tools for reading PoE data files
====================================================

Installation
------------

    go get -u github.com/oriath.net/pogo

You'll need to have [Go 1.14](https://golang.org/dl/) or later installed.


Usage
-----

    pogo ggpk list Content.ggpk
    pogo ggpk ls Content.ggpk:/Data

List the contents of a GGPK file.


    pogo ggpk extract --no-recurse --into root/ Content.ggpk Data/

Extract the `Data` directory from `Content.ggpk` into the `root` directory,
skipping any subdirectories.

(The `ggpk` tool supports several other useful options; run `ggpk --help` for
details.)


     pogo cat Content.ggpk:Metadata/StatDescriptions/stat_descriptions.txt

Output the contents of a single file in the GGPK to standard output.

(Piping the output of this command to `iconv -f utf16 -t utf8` may be useful,
as most text files in the container are encoded as UTF-16.)


    pogo data2json -f formats/demo.go Content.ggpk:Data/ActiveSkills.dat

Dump a data file (directly from the GGPK!) to JSON. This currently only works
for ActiveSkills (since it's the only format specified in `formats/demo.go`),
but adding other formats to that or another file will make them supported as
well.


    pogo analyze Content.ggpk:Data/ActiveSkills.dat

Analyze the contents of ActiveSkills.dat, providing information which may be
useful in interpreting it.
