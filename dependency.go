package depbot

// Dependency found within a source code, this struct is used to capture the data
// And be able to transport it.
type Dependency struct {
	// File where the dependency is defined.
	File string

	// Name of the dependency.
	Name string

	// Version of the dependency.
	Version string

	// License type of the dependency.
	License string

	// Kind of the dependency it could be language or library.
	Kind DependencyKind

	// Whether the dependency is direct or transitive.
	Direct bool

	Timestamp int64
}

type Dependencies []Dependency

func (d Dependencies) Languages() Dependencies {
	return Dependencies{}
}
