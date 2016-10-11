package models

import (
	"strings"
)

type AllPathResponseBody struct {
	Paths []Path `json:"paths"`
}

type Path struct {
	Path string `json:"path,omitempty"`
}

func (allPathResponseBody AllPathResponseBody) Terminal() string {
	lines := []string{}
	lines = append(lines, "Path")
	for _, path := range allPathResponseBody.Paths {
		lines = append(lines, path.Path)
	}
	return strings.Join(lines, "\n")
}

func (allPathsResponseBody AllPathResponseBody) Json() string {
	return ""
}
