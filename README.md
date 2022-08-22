Path of Go - Golang tools for reading PoE data files
====================================================

Installation
------------

    go get -u github.com/oriath-net/pogo

You'll need to have [Go 1.18](https://golang.org/dl/) or later installed.


Usage
-----

    pogo ls Content.ggpk:Data

    pogo ls 'Path of Exile':

    pogo ls some_files.zip:subdirectory/

List the contents of a GGPK file, Steam install, or PKZIP archive, including
files within bundles if appropriate.

The colon in this (and other) commands is significant; it separates files in
the real filesystem on the left from files within the virtual filesystem on
the right.

(Support for PKZIP archives is provided as a utility for developers. The game
never naturally provides data in this format.)


     pogo cat --utf16 Content.ggpk:Metadata/StatDescriptions/stat_descriptions.txt

Output the contents of a single file in the GGPK to standard output,
converting UTF-16 to UTF-8 for output.


    pogo extract Content.ggpk:Audio/

Extract all audio files from the GGPK.


    pogo data2json Content.ggpk:Data/WorldAreas.dat

Dump a data file to JSON. Data file formats are specified in `dat/formats`;
the contents of that directory are embedded in the pogo executable.

The [`jq`](https://stedolan.github.io/jq/) tool is incredibly useful for
displaying and manipulating JSON data, and is highly recommended.


    pogo analyze Content.ggpk:Data/WorldAreas.dat

Analyze the contents of ActiveSkills.dat, providing information which may be
useful in interpreting it.


Most of these commands have a bunch of additional options; they can be viewed
from within the application using e.g. `pogo help cat`.
