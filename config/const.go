package config

const (
	// FileNameState is the name of state file
	FileNameState = "state.yaml"
	// FileNameSpec is the name of spec file
	FileNameSpec = "spec.toml"

	// DirNameRoot is the name of directory for root certificate authority
	DirNameRoot = "root"
	// DirNameInterm is the name of directory for intermediate certificate authorities
	DirNameInterm = "intermediates"
	// DirNameServer is the name of directory for server certificates
	DirNameServer = "servers"
	// DirNameClient is the name of directory for client certificates
	DirNameClient = "clients"

	defaultRootCASerial = int64(10)
	defaultRootCALength = 4096
	defaultRootCADays   = 20 * 365

	defaultIntermCASerial = int64(100)
	defaultIntermCALength = 4096
	defaultIntermCADays   = 10 * 365

	defaultServerCertSerial = int64(1000)
	defaultServerCertLength = 2048
	defaultServerCertDays   = 10 + 365

	defaultClientCertSerial = int64(10000)
	defaultClientCertLength = 2048
	defaultClientCertDays   = 10 + 30
)
