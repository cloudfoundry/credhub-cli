package util

import "github.com/fatih/color"

var Warning = color.New(color.Bold, color.FgYellow).PrintlnFunc()

func Header(modelType, name string) string {
	return BuildLineOfFixedLength("Type:", modelType) + "\n" +
		BuildLineOfFixedLength("Name:", name) + "\n"
}

func Footer(updatedAt string) string {
	return BuildLineOfFixedLength("Updated:", updatedAt)
}