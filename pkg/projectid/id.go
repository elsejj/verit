package projectid

import "strings"

const (
	Multiple = ProjectID(iota + 10)
	Python
	Go
	Node
	Flutter
	Rust
	MaxProjectID
)

type ProjectID int

func (p ProjectID) String() string {
	switch p {
	case Multiple:
		return "Multiple"
	case Python:
		return "Python"
	case Go:
		return "Go"
	case Node:
		return "Node"
	case Flutter:
		return "Flutter"
	case Rust:
		return "Rust"
	default:
		return "Unknown"
	}
}

func ParseProjectID(s string) ProjectID {
	switch strings.ToLower(s) {
	case "multiple":
		return Multiple
	case "python":
		return Python
	case "go":
		return Go
	case "node":
		return Node
	case "flutter":
		return Flutter
	case "rust":
		return Rust
	default:
		return 0
	}
}

func (p ProjectID) Project(workdir string) Project {
	switch p {
	case Multiple:
		m := &MultipleProject{
			workdir: workdir,
		}
		m.scanProjects()
		return m
	case Python:
		return &PythonProject{
			workdir: workdir,
		}
	case Go:
		return &GoProject{
			workdir: workdir,
		}
	case Node:
		return &NodeProject{
			workdir: workdir,
		}
	case Flutter:
		return &FlutterProject{
			workdir: workdir,
		}
	case Rust:
		return &RustProject{
			workdir: workdir,
		}
	default:
		return nil
	}
}
