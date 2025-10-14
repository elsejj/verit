package foo

import _ "embed"

//go:embed version.txt
var Version string
