package cmd

import _ "unsafe"

//go:linkname version main.version
var version string = "<unknown>"
