# gocert
A lightweight library and also command-line interface for generating self-signed SSL/TLS certificates using pure go!

## Installing

## Quick Start

```
mkdir certs
cd certs

gocert init -short
gocert root

gocert intermediate -name=ops
gocert sign -ca=root -name=ops

gocert server -name=webapp
gocert client -name=myservice
gocert sign -ca=ops -name=webapp,myservice
```

## Certificates Explained
You can generate the following types of certificates:
  * Root Certificate Authority
  * Intermediate Certificate Authority
  * Server Certificate
  * Client Certificate

**Root CA** is only used for signing intermediate CA.
It never signs other server or client certificates directly.
It should be keep secured, offline, and unused as much as possible.

**Intermediate CA** is used for signing server and client certificates.
If the intermediate key is comprised, the root CA can revoke the intermediate CA and create a new one.

**Server** certificates can be used for securing servers and establishing SSL/TLS servers.

**Client** certificates can be used for client authentication and MTLS communications between services.

## Generating Certificates
For generating a new chain of certificates, create a new directory and inside that run the following command:

```
gocert init
```

You will be first asked for entering **common specs** which all of your certificates share. So, you enter them once.
Next, you will be asked for entering more-specific specs for **Root CA**, **Intermediate CA**, **Server**, and **Client** certificates separately.
You can enter a list by comma-separating values. If you don't want to use any of the specs, leave it empty.
You can later change these specs by editing `spec.toml` file.

Here is the default configurations for certificates:

| Type         | Key Length | Expiry Days     |
| ------------ | ---------- | --------------- |
| Root         | 4096       | 7300 (20 years) |
| Intermediate | 4096       | 3650 (10 years) |
| Server       | 2048       | 375 (1 year)    |
| Client       | 2048       | 40 (1 month)    |

You can change these configurations by editing `state.yaml` file.
