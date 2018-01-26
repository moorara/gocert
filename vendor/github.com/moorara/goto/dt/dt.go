package dt

type (
	// Generic represents a generic data type
	Generic interface{}

	// Compare is used for comparing
	Compare func(a Generic, b Generic) int

	// BitString is used for bit-string representation
	BitString func(a Generic) []byte
)
