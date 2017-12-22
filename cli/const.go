package cli

const (
	rootName = "root"

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
	// ErrorInvalidName is returned when no name is provided
	ErrorInvalidName = 22
	// ErrorInvalidCA is returned when an invalid ca name is set
	ErrorInvalidCA = 23
	// ErrorInvalidCSR is returned when an invalid csr name is set
	ErrorInvalidCSR = 24
	// ErrorInvalidCert is returned when an invalid cert is set
	ErrorInvalidCert = 25

	// ErrorCert is returned when generating root ca fails
	ErrorCert = 31
	// ErrorCSR is returned when generating csr fails
	ErrorCSR = 32
	// ErrorSign is returned when signing a csr fails
	ErrorSign = 33
	// ErrorVerify is returned when verifying a cert fails
	ErrorVerify = 34
)
