package pki

import (
	"crypto/x509"
	"reflect"
	"regexp"
)

var (
	regex = map[string]*regexp.Regexp{
		"CommonName":         regexp.MustCompile("(?i)^Common[_-]?Name$"),
		"Country":            regexp.MustCompile("(?i)^Country$"),
		"Province":           regexp.MustCompile("(?i)^Province$"),
		"Locality":           regexp.MustCompile("(?i)^Locality$"),
		"Organization":       regexp.MustCompile("(?i)^Organization$"),
		"OrganizationalUnit": regexp.MustCompile("(?i)^Organizational[_-]?Unit$"),
		"DNSNames":           regexp.MustCompile("(?i)^DNS[_-]?Name(s)?$"),
		"IPAddresses":        regexp.MustCompile("(?i)^IP[_-]?Address(es)?$"),
		"EmailAddresses":     regexp.MustCompile("(?i)^Email[_-]?Address(es)?$"),
		"StreetAddress":      regexp.MustCompile("(?i)^Street[_-]?Address$"),
		"PostalCode":         regexp.MustCompile("(?i)^Postal[_-]?Code$"),
	}
)

type (
	// TrustFunc is the function for determing if a ca can sign a csr
	TrustFunc func(*x509.Certificate, *x509.CertificateRequest) bool
)

func matches(cert *x509.Certificate, req *x509.CertificateRequest, fieldName string) bool {
	zero := reflect.Value{}

	certV := reflect.ValueOf(cert).Elem()
	certFV := certV.FieldByName(fieldName)
	if certFV == zero {
		certFV = certV.FieldByName("Subject").FieldByName(fieldName)
	}

	reqV := reflect.ValueOf(req).Elem()
	reqFV := reqV.FieldByName(fieldName)
	if reqFV == zero {
		reqFV = reqV.FieldByName("Subject").FieldByName(fieldName)
	}

	if certFV == zero || reqFV == zero || !reflect.DeepEqual(certFV.Interface(), reqFV.Interface()) {
		return false
	}

	return true
}

func supplied(req *x509.CertificateRequest, fieldName string) bool {
	zero := reflect.Value{}

	reqV := reflect.ValueOf(req).Elem()
	reqFV := reqV.FieldByName(fieldName)
	if reqFV == zero {
		reqFV = reqV.FieldByName("Subject").FieldByName(fieldName)
	}

	if (reqFV == zero) ||
		(reqFV.Kind() == reflect.String && reqFV.String() == "") ||
		(reqFV.Kind() == reflect.Slice && reqFV.Len() == 0) {
		return false
	}

	return true
}

// PolicyTrustFunc returns a TrustFunc using Policy
func PolicyTrustFunc(policy Policy) TrustFunc {
	return func(ca *x509.Certificate, csr *x509.CertificateRequest) bool {
		if ca == nil || csr == nil {
			return false
		}

		// These fields should be matched
		for _, field := range policy.Match {
			for key, re := range regex {
				if re.MatchString(field) && !matches(ca, csr, key) {
					return false
				}
			}
		}

		// These fields should be present
		for _, field := range policy.Supplied {
			for key, re := range regex {
				if re.MatchString(field) && !supplied(csr, key) {
					return false
				}
			}
		}

		return true
	}
}
