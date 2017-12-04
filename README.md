# gocert
A lightweight command-line tool for generating self-signed SSL/TLS certificates using pure go!

## Installing

## Type of Certificates
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

**Server** certificates can be used for securing servers and establishing SSL/TLS server.

**Client** certificates can be used for client authentication and MTLS communications.

## Generating Certificates
For generating new sets of certificates, create a new directory and inside that run the following command:

```
gocert new
```

You will be first asked for entering **common specs** which all of your certificates share. So, you enter them once.
Next, you will be asked for entering more-specific specs for **Root CA**, **Intermediate CA**, **Server**, and **Client** certificates.
You can enter a list by comma-separating values. If you don't want to use any of the specs, leave it empty.
You can later change these specs by editing `spec.toml` file.

Here is the default settings for certificates:

| Type         | Key Length | Serial Number | Expiry Days     |
| ------------ | ---------- | ------------- | --------------- |
| Root         | 4098       | 10            | 7300 (20 years) |
| Intermediate | 4098       | 100           | 3650 (10 years) |
| Server       | 2048       | 1000          | 375 (1 year)    |
| Client       | 2048       | 10000         | 40 (1 month)    |

You can change these settings by editing `state.yaml` file.
