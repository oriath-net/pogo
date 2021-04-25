#!/bin/sh
xmllint --noout --dtdvalid typedef.dtd xml/*.xml && echo "All files OK"
