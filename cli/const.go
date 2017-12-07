package cli

const (
	// ErrorMakeDir is returned when cannot make a directory
	ErrorMakeDir = 11
	// ErrorWriteState is returned when cannot write state file
	ErrorWriteState = 12
	// ErrorWriteSpec is returned when cannot write spec file
	ErrorWriteSpec = 13

	// ErrorInvalidFlag is returned when an invalid flag is provided
	ErrorInvalidFlag = 21
	// ErrorReadState is returned when cannot read state
	ErrorReadState = 22
	// ErrorReadSpec is returned when cannot read spec file
	ErrorReadSpec = 23

	// ErrorRootCA is returned when generating root ca failed
	ErrorRootCA = 31
	// ErrorIntermCA is returned when generating intermediate ca failed
	ErrorIntermCA = 32
)
