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
$ verit -b 1
# bump minor version
$ verit -b 2
# bump patch version
$ verit -b 3
```

## set project version

```bash
$ verit -v 1.2.3
```

# Implementation

for some language project, it has its own way to manage version, `verit` will use the way to manage version.

for some language project, it has no way to manage version, `verit` will do it for you.

# Go Project

for go project, `verit` will use a generated file to manage version, the generated file is `internal/version/version.go`

# Python Project

for python project, `verit` will use `pyproject.toml` to manage version,

# Node Project

for node project, `verit` will use `package.json` to manage version,

