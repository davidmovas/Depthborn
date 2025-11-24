package types

// Identity represents any identifiable game object
type Identity interface {
	// ID returns unique identifier
	ID() string

	// Type returns type name for categorization
	Type() string
}

// Named represents objects with display name
type Named interface {
	// Name returns display name
	Name() string

	// SetName updates display name
	SetName(name string)
}

// Leveled represents objects with level progression
type Leveled interface {
	// Level returns current level
	Level() int

	// SetLevel updates level
	SetLevel(level int)
}

// Tagged represents objects that can be tagged
type Tagged interface {
	// Tags returns tag set
	Tags() TagSet
}

// TagSet manages string tags
type TagSet interface {
	// Add adds tag
	Add(tag string)

	// Remove removes tag
	Remove(tag string)

	// Has checks if tag exists
	Has(tag string) bool

	// Contains checks if all tags exist
	Contains(tags ...string) bool

	// ContainsAny checks if any tag exists
	ContainsAny(tags ...string) bool

	// All returns all tags
	All() []string

	// Clear removes all tags
	Clear()
}
