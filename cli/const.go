package cli

const (
	// ErrorMakeDir is returned when cannot make a directory
	ErrorMakeDir = 11
	// ErrorWriteState is returned when cannot write state file
	ErrorWriteState = 12
	// ErrorWriteSpec is returned when cannot write spec file
	ErrorWriteSpec = 13
	// ErrorReadState is returned when cannot read state
	ErrorReadState = 14
	// ErrorReadSpec is returned when cannot read spec file
	ErrorReadSpec = 15

	// ErrorInvalidFlag is returned when an invalid flag is provided
	ErrorInvalidFlag = 21
	// ErrorNoName is returned when no name is provided
	ErrorNoName = 22

	// ErrorMetadata is returned when an invalid metadata is set
	ErrorMetadata = 31
	// ErrorRootCA is returned when generating root ca failed
	ErrorRootCA = 32
	// ErrorCSR is returned when generating csr failed
	ErrorCSR = 33
)
