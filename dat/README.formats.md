About formats
=============

The JSON format description files under `dat/formats` describe the contents of
Path of Exile data files across multiple game versions.

Each format description file must be named `FileName.json`, where `FileName` is
the base name of the data file it corresponds to. (For example, `Mods.dat` and
`Mods.dat64` are described by `Mods.json`.)


Top level
---------

The top level of a format description must be a JSON object with the keys
`file` and `fields`, and optionally with the fields `description`, `since`,
and `until`.

`file` must be set to the same base name as in the file name.

`fields` must be set to an array of fields, as defined below.

`description` may contain some freeform text describing the file and its
status.

`since` and `until` may be set to the first known version numbers in which
this file was present and removed, respectively. These fields are currently
treated as informational; they are not used during parsing.


Fields
------

Each field must be a JSON object with the keys `name` and `type`, and
optionally with a number of additional fields listed below.

`name` is the name of the field. It must be unique within the format
description, must begin with a capital letter, and must be a valid identifier
(no spaces or punctuation other than `_`).

`type` is the type of the field. Valid types are listed below.

Optional keys are:

* `description`: as above, freeform text describing the field and its status.

* `since` and `until`: first known version numbers in which this field was
  present and removed, respectively. Unlike the top-level fields described
  above, these **are** significant to the data file parser.

* `unique`: if set to the boolean `true`, this field is believed to contain a
  distinct value for every row. This is currently treated as informational.

* `ref`: may be set to the name of another data file which this field
  references.

* `ref-field`: if `ref` is set, may be set to the name of a field within that
  file. By default, `ref` is assumed to be the 0-based index of the row,
  rather than the value of any particular field.

* `path`: may be set to a string representing the extension of files which
  this field typically references.

  As a special case, `art` indicates that the field contains a symbolic path
  to artwork which can be resolved to a DDS path and bounding box through the
  files:

  * `Art/UIImages1.txt`
  * `Art/UIDivinationImages.txt`
  * `Art/UIShopImages.txt`
  * `Art/UIPS4.txt`
  * `Art/UIXbox.txt`


Types
-----

* `bool` - 1-byte value, always 1=true or 0=false
* `u32` - 4-byte unsigned integer value
* `u64` - 8-byte unsigned integer value
* `i32` - 4-byte signed integer value
* `i64` - 8-byte signed integer value
* `f32` - 4-byte floating-point value
* `string` - 4/8 byte dynamic offset to a UTF-16 string (or UTF-32 in `.datl`)
* `shortid` - 4/8 byte integer referencing a row in the current table
* `longid` - 8/16 byte integer referencing a row in another table

Each of these types can be suffixed by `[]` to indicate an array. An array
always exists within the row as a 4/8 byte element count followed by a 4/8
byte dynamic offset to the first value.

The following types are currently supported, but are theorized to not actually
exist (i.e. any extant fields with these types contain misparsed data):

* `u8`  - 1-byte unsigned integer value
* `u16` - 2-byte unsigned integer value
* `i16` - 2-byte signed integer value
* `f64` - 8-byte floating-point value


Schema validation
-----------------

You can validate the JSON schema of one or more format descriptions using the
`jv` utility, which can be installed using:

    go install github.com/santhosh-tekuri/jsonschema/cmd/jv@latest

With this installed, run:

    jv dat/formats.schema.json dat/formats/*.json

Any structural errors in format description files will be reported.
