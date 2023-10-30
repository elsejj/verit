package projectid

import "strings"

const (
	Python = ProjectID(iota + 10)
	Go
	Node
)

type ProjectID int

func (p ProjectID) String() string {
	switch p {
	case Python:
		return "Python"
	case Go:
		return "Go"
	case Node:
		return "Node"
	default:
		return "Unknown"
	}
}

func ParseProjectID(s string) ProjectID {
	switch strings.ToLower(s) {
	case "python":
		return Python
	case "go":
		return Go
	case "node":
		return Node
	default:
		return 0
	}
}

func (p ProjectID) Project(workdir string) Project {
	switch p {
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
	default:
		return nil
	}
}
