[![Build Status][travisci-image]][travisci-url]

# gocert
A lightweight library and also command-line interface for generating self-signed SSL/TLS certificates using pure go!

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

Here is the default configs for certificates:

| Type         | Key Length | Expiry Days     |
| ------------ | ---------- | --------------- |
| Root         | 4096       | 7300 (20 years) |
| Intermediate | 4096       | 3650 (10 years) |
| Server       | 2048       | 375 (1 year)    |
| Client       | 2048       | 40 (1 month)    |

You can change these configs by editing `state.yaml` file.

## Root CA
For generating your root certificate authority, run the following command:

```
gocert root
```

You can only have one root certificate authority and it is called `root` by default.
Setting a password for root certificate authority key is **mandatory**.

## Intermediate CAs
For generating an intermediate certificate authority, run the following command:

```
gocert intermediate -name=<...>
```

You can have one or more intermediate certificate authorities.
Setting password for intermediate certificate authorities are **mandatory** too.

Then, you need to sign your intermediate ca by root ca as follows:

```
gocert sign -ca=root -name=<intermediate_name>
```

You will be asked for entering the password for root ca.

## Server Certificates
You can generate a server certificate by running:

```
gocert server -name=<server_name>
```

The `CommonName` for server certificates must be a **Fully Qualified Domain Name** (FQDN).

Your server certificate should be signed by an intermediate certificate.

```
gocert sign -ca=<intermediate_name> -name=<server_name>
```

## Client Certificates
You can generate a client certificate similarly:

```
gocert client -name=<client_name>
```

Likewise, your client certificate should be signed by an intermediate certificate.

```
gocert sign -ca=<intermediate_name> -name=<client_name>
```


[travisci-url]: https://travis-ci.com/moorara/gocert
[travisci-image]: https://travis-ci.com/moorara/gocert.svg?branch=master&token=HyJPFzY74fNDzrcekdXq
