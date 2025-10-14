# Repository Guidelines

This guide outlines how to contribute effectively to verit, a Go CLI for managing project versions across language ecosystems.

## Design

`verit` is designed as a single binary CLI tool that bump a project's version in place, based on semantic versioning principles.
When a bump is requested, it also tags the current commit in git with the new version.

### How It Works

1. detect the project type (Go, Python, Node, etc) by looking for characteristic files in the current directory.
2. parse the current version from the appropriate file (e.g., `go.mod`, `pyproject.toml`, `package.json`).
3. increment the version according to the requested bump type (major, minor, patch).
4. write the new version back to the file.
5. create a git tag for the new version (optional).

## Project Structure & Module Organization

`main.go` houses the CLI entrypoint and wires sub packages together.
Core logic for detecting project types lives in `pkg/projectid`, with language-specific detectors in `p_go.go`, `p_python.go`, and `p_node.go`.
Each project type implements the `Project` interface defined in `pkg/projectid/projectid.go`.
Shared version helpers reside in `pkg/version`.
Go builds deposit binaries into `dist/`, while `internal/version/version.go` stores the generated version metadata that ships with each release.

## Build, Test, and Development Commands

- `make build`: runs `go build` with stripped symbols and writes `dist/verit`. Use after code changes to confirm a reproducible binary.
- `make tiny`: compiles with TinyGo and strips the artifact for minimal size; keep it green before publishing lightweight releases.
- `go test ./...`: executes all Go unit tests; add it to your local pre-push workflow even though the suite is currently sparse.

## Coding Style & Naming Conventions

Follow standard Go 1.20 practices: tabs for indentation, exported identifiers use CamelCase, and packages stay lower*snake for clarity. Run `gofmt -w` on touched files and `goimports` if available to keep imports tidy. When adding cross-language helpers, mirror the existing `p*<lang>.go` naming to maintain parity.

## Testing Guidelines

Place tests in the same package with `_test.go` suffixes and prefer table-driven cases for parsers and formatters. Stub file system access with temporary directories to avoid altering real project files. Aim for meaningful assertions around version parsing, bumping, and persistence; when coverage dips, document rationale in the PR.

## Commit & Pull Request Guidelines

History favors short, imperative commit messages (e.g., `add code`, `rm dist`). Continue that style and limit each commit to one logical change. PRs should describe the intent, list validation steps (`make build`, `go test ./...`), and reference related issues.
