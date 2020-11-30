[![Build Status][workflow-image]][workflow-url]
[![Go Report Card][goreport-image]][goreport-url]
[![Test Coverage][coverage-image]][coverage-url]
[![Maintainability][maintainability-image]][maintainability-url]

# gocert

If you are having a hard time every time using `openssl` for generating self-signed certificates, this tool is for you!
A lightweight library and also command-line interface for generating self-signed SSL/TLS certificates using pure go.

[![asciicast](https://asciinema.org/a/vGNpB4ClRhBBoR3KOH6EVRzpH.svg)](https://asciinema.org/a/vGNpB4ClRhBBoR3KOH6EVRzpH)

## Install

```
brew install moorara/brew/gocert
```

For other platforms, you can download the binary from the [latest release](https://github.com/moorara/gocert/releases/latest).

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

  - Root Certificate Authority
  - Intermediate Certificate Authority
  - Server Certificate
  - Client Certificate

**Root CA** is only used for signing intermediate CA.
There is only one root CA called `root` by default.
Root CA never signs user certificates (server or client) directly.
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


[workflow-url]: https://github.com/moorara/gocert/actions
[workflow-image]: https://github.com/moorara/gocert/workflows/Main/badge.svg
[goreport-url]: https://goreportcard.com/report/github.com/moorara/gocert
[goreport-image]: https://goreportcard.com/badge/github.com/moorara/gocert
[coverage-url]: https://codeclimate.com/github/moorara/gocert/test_coverage
[coverage-image]: https://api.codeclimate.com/v1/badges/c42cb8902ef865a053eb/test_coverage
[maintainability-url]: https://codeclimate.com/github/moorara/gocert/maintainability
[maintainability-image]: https://api.codeclimate.com/v1/badges/c42cb8902ef865a053eb/maintainability
