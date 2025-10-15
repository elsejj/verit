# manage project version

this tool is used to manage project version, it can read/write version from project file, and do bump version, create git tag etc.

# Usage

## Show project version

```bash
verit
```

## Bump project version

```bash
# bump major version
verit -M
# bump minor version
verit -m
# bump patch version
verit -p
```

## Bump to specific version

```bash
# bump major version to 2
verit -M=2
# bump minor version to 3
verit -m=3
# bump patch version to 4
verit -p=4
# set prerelease version
verit -r alpha.1
# set build metadata
verit -b 001
```

## Set project version

```bash
verit -V 1.2.3
```

## Create git tag with current version

```bash
verit -t
```

## Create git tag with current version and push to remote

```bash
verit -T
```

# Implementation

for some language project, it has its own way to manage version, `verit` will use the way to manage version.

for some language project, it has no way to manage version, `verit` will do it for you.

## Go Project

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

## Python Project

for python project, `verit` will use `pyproject.toml` to manage version,

## Node Project

for node project, `verit` will use `package.json` to manage version,

## Flutter Project

for flutter project, `verit` will use the `version` field in `pubspec.yaml`. the value should follow `x.y.z` with an optional `+build` suffix. when bumping or setting the version, `verit` keeps the build number if provided and writes updates back to `pubspec.yaml`.

## Rust Project

for rust project, `verit` will use `Cargo.toml` to manage version,

rust project may have multiple crates, `verit` will only manage the crate in the current working directory.

if you want to unify version for all crates, you can set the version in the workspace `Cargo.toml`, like:

```toml
[workspace]
members = [
    "crate1",
    "crate2",
]

[workspace.package]
version = "1.2.3"
```

then in each crate `Cargo.toml`, you can inherit the version from workspace, like:

```toml
[package]
name = "crate1"
version.workspace = true
```

## Mix Projects

when mix supported project files are present (for example, `pyproject.toml` and `Cargo.toml` together, a typical `Maturin` project, use `pyo3` build python native extensions), `verit` treats the workspace as a single project and applies operations to each manifest. all detected manifests must share the same version before a bump is performed, ensuring the versions stay aligned across languages.
