package pki

import (
	"crypto/x509"
	"reflect"
	"regexp"
)

var (
	regexCommonName         = regexp.MustCompile("(?i)^Common[_-]?Name$")
	regexCountry            = regexp.MustCompile("(?i)^Country$")
	regexProvince           = regexp.MustCompile("(?i)^Province$")
	regexLocality           = regexp.MustCompile("(?i)^Locality$")
	regexOrganization       = regexp.MustCompile("(?i)^Organization$")
	regexOrganizationalUnit = regexp.MustCompile("(?i)^Organizational[_-]?Unit$")
	regexDNSName            = regexp.MustCompile("(?i)^DNS[_-]?Name(s)?$")
	regexIPAddress          = regexp.MustCompile("(?i)^IP[_-]?Address(es)?$")
	regexEmailAddress       = regexp.MustCompile("(?i)^Email[_-]?Address(es)?$")
	regexStreetAddress      = regexp.MustCompile("(?i)^Street[_-]?Address$")
	regexPostalCode         = regexp.MustCompile("(?i)^Postal[_-]?Code$")
)

type (
	// TrustFunc is the function for determing if a ca can sign a csr
	TrustFunc func(*x509.Certificate, *x509.CertificateRequest) bool
)

// PolicyTrustFunc returns a TrustFunc using Policy
func PolicyTrustFunc(policy Policy) TrustFunc {
	return func(ca *x509.Certificate, csr *x509.CertificateRequest) bool {
		if ca == nil || csr == nil {
			return false
		}

		// These fields should be matched
		for _, name := range policy.Match {
			switch {
			case regexCommonName.MatchString(name):
				if ca.Subject.CommonName != csr.Subject.CommonName {
					return false
				}

			case regexCountry.MatchString(name):
				if !reflect.DeepEqual(ca.Subject.Country, csr.Subject.Country) {
					return false
				}

			case regexProvince.MatchString(name):
				if !reflect.DeepEqual(ca.Subject.Province, csr.Subject.Province) {
					return false
				}

			case regexLocality.MatchString(name):
				if !reflect.DeepEqual(ca.Subject.Locality, csr.Subject.Locality) {
					return false
				}

			case regexOrganization.MatchString(name):
				if !reflect.DeepEqual(ca.Subject.Organization, csr.Subject.Organization) {
					return false
				}

			case regexOrganizationalUnit.MatchString(name):
				if !reflect.DeepEqual(ca.Subject.OrganizationalUnit, csr.Subject.OrganizationalUnit) {
					return false
				}

			case regexDNSName.MatchString(name):
				if !reflect.DeepEqual(ca.DNSNames, csr.DNSNames) {
					return false
				}

			case regexIPAddress.MatchString(name):
				if !reflect.DeepEqual(ca.IPAddresses, csr.IPAddresses) {
					return false
				}

			case regexEmailAddress.MatchString(name):
				if !reflect.DeepEqual(ca.EmailAddresses, csr.EmailAddresses) {
					return false
				}

			case regexStreetAddress.MatchString(name):
				if !reflect.DeepEqual(ca.Subject.StreetAddress, csr.Subject.StreetAddress) {
					return false
				}

			case regexPostalCode.MatchString(name):
				if !reflect.DeepEqual(ca.Subject.PostalCode, csr.Subject.PostalCode) {
					return false
				}
			}
		}

		// These fields should be present
		for _, name := range policy.Supplied {
			switch {
			case regexCommonName.MatchString(name):
				if csr.Subject.CommonName == "" {
					return false
				}

			case regexCountry.MatchString(name):
				if csr.Subject.Country == nil || len(csr.Subject.Country) == 0 {
					return false
				}

			case regexProvince.MatchString(name):
				if csr.Subject.Province == nil || len(csr.Subject.Province) == 0 {
					return false
				}

			case regexLocality.MatchString(name):
				if csr.Subject.Locality == nil || len(csr.Subject.Locality) == 0 {
					return false
				}

			case regexOrganization.MatchString(name):
				if csr.Subject.Organization == nil || len(csr.Subject.Organization) == 0 {
					return false
				}

			case regexOrganizationalUnit.MatchString(name):
				if csr.Subject.OrganizationalUnit == nil || len(csr.Subject.OrganizationalUnit) == 0 {
					return false
				}

			case regexDNSName.MatchString(name):
				if csr.DNSNames == nil || len(csr.DNSNames) == 0 {
					return false
				}

			case regexIPAddress.MatchString(name):
				if csr.IPAddresses == nil || len(csr.IPAddresses) == 0 {
					return false
				}

			case regexEmailAddress.MatchString(name):
				if csr.EmailAddresses == nil || len(csr.EmailAddresses) == 0 {
					return false
				}

			case regexStreetAddress.MatchString(name):
				if csr.Subject.StreetAddress == nil || len(csr.Subject.StreetAddress) == 0 {
					return false
				}

			case regexPostalCode.MatchString(name):
				if csr.Subject.PostalCode == nil || len(csr.Subject.PostalCode) == 0 {
					return false
				}
			}
		}

		return true
	}
}
