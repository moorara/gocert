package cli

const (
	rootName = "root"

	mdRootSkip   = "rootSkip"
	mdIntermSkip = "intermSkip"
	mdServerSkip = "serverSkip"
	mdClientSkip = "clientSkip"

	promptTemplate = "%s (type: %s):"

	// ErrorEnterState is returned when entering state fails
	ErrorEnterState = 11
	// ErrorEnterSpec is returned when entering spec fails
	ErrorEnterSpec = 12
	// ErrorEnterConfig is returned when entering config fails
	ErrorEnterConfig = 13
	// ErrorEnterClaim is returned when entering claim fails
	ErrorEnterClaim = 14

	// ErrorMakeDir is returned when cannot make a directory
	ErrorMakeDir = 21
	// ErrorWriteState is returned when cannot write state file
	ErrorWriteState = 22
	// ErrorWriteSpec is returned when cannot write spec file
	ErrorWriteSpec = 23
	// ErrorReadState is returned when cannot read state
	ErrorReadState = 24
	// ErrorReadSpec is returned when cannot read spec file
	ErrorReadSpec = 25

	// ErrorInvalidFlag is returned when an invalid flag is provided
	ErrorInvalidFlag = 31
	// ErrorInvalidName is returned when no name is provided
	ErrorInvalidName = 32
	// ErrorInvalidCA is returned when an invalid ca name is set
	ErrorInvalidCA = 33
	// ErrorInvalidCSR is returned when an invalid csr name is set
	ErrorInvalidCSR = 34
	// ErrorInvalidCert is returned when an invalid cert is set
	ErrorInvalidCert = 35

	// ErrorCert is returned when generating root ca fails
	ErrorCert = 41
	// ErrorCSR is returned when generating csr fails
	ErrorCSR = 42
	// ErrorSign is returned when signing a csr fails
	ErrorSign = 43
	// ErrorVerify is returned when verifying a cert fails
	ErrorVerify = 44
)
