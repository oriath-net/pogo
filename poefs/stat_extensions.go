package poefs

// These methods may be present on fs.FileInfo objects returned

type StatExtensions interface {
	// Get some information on the source of a file
	Provenance() string

	// Get the SHA256 signature of a file (or nil if unsupported)
	Signature() []byte
}
