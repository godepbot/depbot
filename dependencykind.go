package depbot

// The kind of a dependency in a source code.
type DependencyKind string

const (
	// A library dependency. Most probably open source.
	DependencyKindLibrary DependencyKind = "library"

	// Language that the codebase depends on.
	DependencyKindLanguage DependencyKind = "language"

	// Tools like gofmt, yarn, npm, cargo etc.
	DependencyKindTool DependencyKind = "tool"
)
