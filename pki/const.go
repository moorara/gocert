package pki

const (
	// CertTypeRoot represents a root certificate authority
	CertTypeRoot = 1 + iota
	// CertTypeInterm represents an intermediate certificate authority
	CertTypeInterm
	// CertTypeServer represents a server certificate
	CertTypeServer
	// CertTypeClient represents a client certificate
	CertTypeClient
)

const (
	// DirRoot is the name of directory for root certificate authority
	DirRoot = "root"
	// DirInterm is the name of directory for intermediate certificate authorities
	DirInterm = "intermediate"
	// DirServer is the name of directory for server certificates
	DirServer = "server"
	// DirClient is the name of directory for client certificates
	DirClient = "client"
	// DirCSR is the name of directory for certificate signing requests
	DirCSR = "csr"

	// FileState is the name of state file
	FileState = "state.yaml"
	// FileSpec is the name of spec file
	FileSpec = "spec.toml"
)
