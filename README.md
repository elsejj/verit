# manage project version

this tool is used to manage project version, it can read/write version from project file, it can also bump version in project file.

# Usage

## show project version

```bash
$ verit
```

## bump project version

```bash
# bump major version
$ verit -M
# bump minor version
$ verit -m
# bump patch version
$ verit -p
```

## set project version

```bash
$ verit -v 1.2.3
```

# Implementation

for some language project, it has its own way to manage version, `verit` will use the way to manage version.

for some language project, it has no way to manage version, `verit` will do it for you.

# Go Project

For go project, `verit` will lookup `version.txt` file to manage version, if not exist, will give up.
if you want to use `verit` to manage go project version, you need to create a `version.txt` and put version like `1.2.3` in it. then use `go:embed` directive to embed it to a go variable.

For example:

```go
package main
import (
  _ "embed"
  "fmt"
)
//go:embed version.txt
var version string
func main() {
  fmt.Println("version:", version)
}
```

# Python Project

for python project, `verit` will use `pyproject.toml` to manage version,

# Node Project

for node project, `verit` will use `package.json` to manage version,
