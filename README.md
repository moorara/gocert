[![Build Status][travisci-image]][travisci-url]

# gocert
If you are having a hard time every time using `openssl` for generating self-signed certificates, this tool is for you!
A lightweight library and also command-line interface for generating self-signed SSL/TLS certificates using pure go.

## Installing
You can download the appropriate binary from [releases](https://github.com/moorara/gocert/releases) page.

## Quick Start

```
mkdir certs
cd certs

gocert init
gocert root

gocert intermediate -name=sre
gocert sign -ca=root -name=sre

gocert server -name=webapp
gocert client -name=myservice
gocert sign -ca=sre -name=webapp,myservice

gocert verify -ca=root -name=sre
gocert verify -ca=sre -name=webapp,myservice
```

## Certificates Explained
You can generate the following types of certificates:
  * Root Certificate Authority
  * Intermediate Certificate Authority
  * Server Certificate
  * Client Certificate

**Root CA** is only used for signing intermediate CA.
There is only one root CA called `root` by default.
Root CA never signs user certificates (server or client)                         directly.
It should be keep secured, offline, and unused as much as possible.

**Intermediate CA** is used for signing server and client certificates.
It must be signed by `root` CA.
If an intermediate key is comprised, the root CA can revoke the intermediate CA and create a new one.

**Server** certificates can be used for securing servers and establishing SSL/TLS servers.
They should be signed by an intermediate certificate.
The `CommonName` for server certificates must be a *Fully Qualified Domain Name* (FQDN).

**Client** certificates can be used for client authentication and MTLS communications between services.
They should be signed by an intermediate certificate.

### Default Configs

| Type         | Key Length | Expiry Days     |
| ------------ | ---------- | --------------- |
| Root         | 4096       | 7300 (20 years) |
| Intermediate | 4096       | 3650 (10 years) |
| Server       | 2048       | 375 (~1 year)   |
| Client       | 2048       | 40 (~1 month)   |

You can change these configs by editing `state.yaml` file.


[travisci-url]: https://travis-ci.com/moorara/gocert
[travisci-image]: https://travis-ci.com/moorara/gocert.svg?branch=master&token=HyJPFzY74fNDzrcekdXq
