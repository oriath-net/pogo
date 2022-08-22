# Current status of data file parsing

* **Releases** is the range of major releases that data files were included
  in. "all" means that the file has been observed in all available releases.

* **Structure** is how much work has been done to determine the structure of
  this data file:

  * ✅  means that the file can be parsed without warnings.

  * ⚠️  means that the file parses with warnings, usually indicating that some
    dynamic data is not being processed correctly.

  * ❌  means that the file has a format defined, but it fails to parse the
    current version of the data file.

  * ❓  means that no format is defined for this file.

* **dat64** indicates whether the file's structure is compatible with dat64
  decoding.

  * ✅  means that the current dat64 file can be parsed without warnings, and
    yields the same results as the legacy dat file.

  * ⚠️  means that the dat64 file parses with warnings, or that parsing yields
    different results from the dat file.

  * ❌  means that the dat64 file fails to parse at all.

  * "n/a" means that no dat64 version of this file exists, probably because it
    was removed before the 2.5.0 release.

* **History** indicates whether the structure is consistent with all release
  versions of the data file.

  * ✅  means that all historical versions can be parsed.

  * ❌  means that some or all historical versions fail to parse.
