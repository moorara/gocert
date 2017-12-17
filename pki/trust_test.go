package pki

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPolicyTrustFunc(t *testing.T) {
	tests := []struct {
		title          string
		policy         Policy
		ca             *x509.Certificate
		csr            *x509.CertificateRequest
		expectedResult bool
	}{
		{
			"NoCA",
			Policy{},
			nil,
			nil,
			false,
		},
		{
			"NoCSR",
			Policy{},
			&x509.Certificate{},
			nil,
			false,
		},
		{
			"CommonNameNotSupplied",
			Policy{
				Supplied: []string{"common_name"},
			},
			&x509.Certificate{},
			&x509.CertificateRequest{
				Subject: pkix.Name{},
			},
			false,
		},
		{
			"CountryNotSupplied",
			Policy{
				Supplied: []string{"country"},
			},
			&x509.Certificate{},
			&x509.CertificateRequest{
				Subject: pkix.Name{},
			},
			false,
		},
		{
			"CountryNotSupplied",
			Policy{
				Supplied: []string{"country"},
			},
			&x509.Certificate{},
			&x509.CertificateRequest{
				Subject: pkix.Name{
					Country: []string{},
				},
			},
			false,
		},
		{
			"ProvinceNotSupplied",
			Policy{
				Supplied: []string{"province"},
			},
			&x509.Certificate{},
			&x509.CertificateRequest{
				Subject: pkix.Name{},
			},
			false,
		},
		{
			"ProvinceNotSupplied",
			Policy{
				Supplied: []string{"province"},
			},
			&x509.Certificate{},
			&x509.CertificateRequest{
				Subject: pkix.Name{
					Province: []string{},
				},
			},
			false,
		},
		{
			"LocalityNotSupplied",
			Policy{
				Supplied: []string{"locality"},
			},
			&x509.Certificate{},
			&x509.CertificateRequest{
				Subject: pkix.Name{},
			},
			false,
		},
		{
			"LocalityNotSupplied",
			Policy{
				Supplied: []string{"locality"},
			},
			&x509.Certificate{},
			&x509.CertificateRequest{
				Subject: pkix.Name{
					Locality: []string{},
				},
			},
			false,
		},
		{
			"OrganizationNotSupplied",
			Policy{
				Supplied: []string{"organization"},
			},
			&x509.Certificate{},
			&x509.CertificateRequest{
				Subject: pkix.Name{},
			},
			false,
		},
		{
			"OrganizationNotSupplied",
			Policy{
				Supplied: []string{"organization"},
			},
			&x509.Certificate{},
			&x509.CertificateRequest{
				Subject: pkix.Name{
					Organization: []string{},
				},
			},
			false,
		},
		{
			"OrganizationalUnitNotSupplied",
			Policy{
				Supplied: []string{"organizationalUnit"},
			},
			&x509.Certificate{},
			&x509.CertificateRequest{
				Subject: pkix.Name{},
			},
			false,
		},
		{
			"OrganizationalUnitNotSupplied",
			Policy{
				Supplied: []string{"organizational_unit"},
			},
			&x509.Certificate{},
			&x509.CertificateRequest{
				Subject: pkix.Name{
					OrganizationalUnit: []string{},
				},
			},
			false,
		},
		{
			"DNSNameNotSupplied",
			Policy{
				Supplied: []string{"DNSNames"},
			},
			&x509.Certificate{},
			&x509.CertificateRequest{},
			false,
		},
		{
			"DNSNameNotSupplied",
			Policy{
				Supplied: []string{"dns_name"},
			},
			&x509.Certificate{},
			&x509.CertificateRequest{
				EmailAddresses: []string{},
			},
			false,
		},
		{
			"IPAddressNotSupplied",
			Policy{
				Supplied: []string{"IPAddresses"},
			},
			&x509.Certificate{},
			&x509.CertificateRequest{},
			false,
		},
		{
			"IPAddressNotSupplied",
			Policy{
				Supplied: []string{"ip_address"},
			},
			&x509.Certificate{},
			&x509.CertificateRequest{
				EmailAddresses: []string{},
			},
			false,
		},
		{
			"EmailAddressNotSupplied",
			Policy{
				Supplied: []string{"EmailAddresses"},
			},
			&x509.Certificate{},
			&x509.CertificateRequest{},
			false,
		},
		{
			"EmailAddressNotSupplied",
			Policy{
				Supplied: []string{"email_address"},
			},
			&x509.Certificate{},
			&x509.CertificateRequest{
				EmailAddresses: []string{},
			},
			false,
		},
		{
			"StreetAddressNotSupplied",
			Policy{
				Supplied: []string{"streetAddress"},
			},
			&x509.Certificate{},
			&x509.CertificateRequest{
				Subject: pkix.Name{},
			},
			false,
		},
		{
			"StreetAddressNotSupplied",
			Policy{
				Supplied: []string{"street_address"},
			},
			&x509.Certificate{},
			&x509.CertificateRequest{
				Subject: pkix.Name{
					StreetAddress: []string{},
				},
			},
			false,
		},
		{
			"PostalCodeNotSupplied",
			Policy{
				Supplied: []string{"postalCode"},
			},
			&x509.Certificate{},
			&x509.CertificateRequest{
				Subject: pkix.Name{},
			},
			false,
		},
		{
			"PostalCodeNotSupplied",
			Policy{
				Supplied: []string{"postal_code"},
			},
			&x509.Certificate{},
			&x509.CertificateRequest{
				Subject: pkix.Name{
					PostalCode: []string{},
				},
			},
			false,
		},
		{
			"CommonNameNotMatched",
			Policy{
				Match: []string{"CommonName"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					CommonName: "SRE CA",
				},
			},
			&x509.CertificateRequest{
				Subject: pkix.Name{
					CommonName: "Ops CA",
				},
			},
			false,
		},
		{
			"CommonNameNotMatched",
			Policy{
				Match: []string{"common_name"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					CommonName: "SRE CA",
				},
			},
			&x509.CertificateRequest{
				Subject: pkix.Name{},
			},
			false,
		},
		{
			"CountryNotMatched",
			Policy{
				Match: []string{"Country"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					Country: []string{"CA"},
				},
			},
			&x509.CertificateRequest{
				Subject: pkix.Name{
					Country: []string{"IR"},
				},
			},
			false,
		},
		{
			"CountryNotMatched",
			Policy{
				Match: []string{"country"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					Country: []string{"CA"},
				},
			},
			&x509.CertificateRequest{
				Subject: pkix.Name{},
			},
			false,
		},
		{
			"CountryNotMatched",
			Policy{
				Match: []string{"country"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					Country: []string{"CA"},
				},
			},
			&x509.CertificateRequest{
				Subject: pkix.Name{
					Country: []string{},
				},
			},
			false,
		},
		{
			"ProvinceNotMatched",
			Policy{
				Match: []string{"Province"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					Province: []string{"Ontario"},
				},
			},
			&x509.CertificateRequest{
				Subject: pkix.Name{
					Province: []string{"California"},
				},
			},
			false,
		},
		{
			"ProvinceNotMatched",
			Policy{
				Match: []string{"province"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					Province: []string{"Ontario"},
				},
			},
			&x509.CertificateRequest{
				Subject: pkix.Name{},
			},
			false,
		},
		{
			"ProvinceNotMatched",
			Policy{
				Match: []string{"province"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					Province: []string{"Ontario"},
				},
			},
			&x509.CertificateRequest{
				Subject: pkix.Name{
					Province: []string{},
				},
			},
			false,
		},
		{
			"LocalityNotMatched",
			Policy{
				Match: []string{"Locality"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					Locality: []string{"Ottawa"},
				},
			},
			&x509.CertificateRequest{
				Subject: pkix.Name{
					Locality: []string{"San Francisco"},
				},
			},
			false,
		},
		{
			"LocalityNotMatched",
			Policy{
				Match: []string{"locality"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					Locality: []string{"Ottawa"},
				},
			},
			&x509.CertificateRequest{
				Subject: pkix.Name{},
			},
			false,
		},
		{
			"LocalityNotMatched",
			Policy{
				Match: []string{"locality"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					Locality: []string{"Ottawa"},
				},
			},
			&x509.CertificateRequest{
				Subject: pkix.Name{
					Locality: []string{},
				},
			},
			false,
		},
		{
			"OrganizationNotMatched",
			Policy{
				Match: []string{"Organization"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					Organization: []string{"Milad"},
				},
			},
			&x509.CertificateRequest{
				Subject: pkix.Name{
					Organization: []string{"Moorara"},
				},
			},
			false,
		},
		{
			"OrganizationNotMatched",
			Policy{
				Match: []string{"organization"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					Organization: []string{"Milad"},
				},
			},
			&x509.CertificateRequest{
				Subject: pkix.Name{},
			},
			false,
		},
		{
			"OrganizationNotMatched",
			Policy{
				Match: []string{"organization"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					Organization: []string{"Milad"},
				},
			},
			&x509.CertificateRequest{
				Subject: pkix.Name{
					Organization: []string{},
				},
			},
			false,
		},
		{
			"OrganizationalUnitNotMatched",
			Policy{
				Match: []string{"OrganizationalUnit"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					OrganizationalUnit: []string{"R&D"},
				},
			},
			&x509.CertificateRequest{
				Subject: pkix.Name{
					OrganizationalUnit: []string{"QE"},
				},
			},
			false,
		},
		{
			"OrganizationalUnitNotMatched",
			Policy{
				Match: []string{"organizationalUnit"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					OrganizationalUnit: []string{"R&D"},
				},
			},
			&x509.CertificateRequest{
				Subject: pkix.Name{},
			},
			false,
		},
		{
			"OrganizationalUnitNotMatched",
			Policy{
				Match: []string{"organizational_unit"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					OrganizationalUnit: []string{"R&D"},
				},
			},
			&x509.CertificateRequest{
				Subject: pkix.Name{
					OrganizationalUnit: []string{},
				},
			},
			false,
		},
		{
			"DNSNameNotMatched",
			Policy{
				Match: []string{"DNSName"},
			},
			&x509.Certificate{
				DNSNames: []string{"example.com"},
			},
			&x509.CertificateRequest{
				DNSNames: []string{"example.org"},
			},
			false,
		},
		{
			"DNSNameNotMatched",
			Policy{
				Match: []string{"DNSNames"},
			},
			&x509.Certificate{
				DNSNames: []string{"example.com"},
			},
			&x509.CertificateRequest{},
			false,
		},
		{
			"DNSNameNotMatched",
			Policy{
				Match: []string{"dns_name"},
			},
			&x509.Certificate{
				DNSNames: []string{"example.com"},
			},
			&x509.CertificateRequest{
				DNSNames: []string{},
			},
			false,
		},
		{
			"IPAddressNotMatched",
			Policy{
				Match: []string{"IPAddress"},
			},
			&x509.Certificate{
				IPAddresses: []net.IP{net.IPv4(8, 8, 8, 8)},
			},
			&x509.CertificateRequest{
				IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
			},
			false,
		},
		{
			"IPAddressNotMatched",
			Policy{
				Match: []string{"IPAddresses"},
			},
			&x509.Certificate{
				IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
			},
			&x509.CertificateRequest{},
			false,
		},
		{
			"IPAddressNotMatched",
			Policy{
				Match: []string{"ip_address"},
			},
			&x509.Certificate{
				IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
			},
			&x509.CertificateRequest{
				IPAddresses: []net.IP{},
			},
			false,
		},
		{
			"EmailAddressNotMatched",
			Policy{
				Match: []string{"EmailAddress"},
			},
			&x509.Certificate{
				EmailAddresses: []string{"milad@example.com"},
			},
			&x509.CertificateRequest{
				EmailAddresses: []string{"moorara@example.com"},
			},
			false,
		},
		{
			"EmailAddressNotMatched",
			Policy{
				Match: []string{"EmailAddresses"},
			},
			&x509.Certificate{
				EmailAddresses: []string{"milad@example.com"},
			},
			&x509.CertificateRequest{},
			false,
		},
		{
			"EmailAddressNotMatched",
			Policy{
				Match: []string{"email_address"},
			},
			&x509.Certificate{
				EmailAddresses: []string{"milad@example.com"},
			},
			&x509.CertificateRequest{
				EmailAddresses: []string{},
			},
			false,
		},
		{
			"StreetAddressNotMatched",
			Policy{
				Match: []string{"StreetAddress"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					StreetAddress: []string{"Apadana"},
				},
			},
			&x509.CertificateRequest{
				Subject: pkix.Name{
					StreetAddress: []string{"Ekbatan"},
				},
			},
			false,
		},
		{
			"StreetAddressNotMatched",
			Policy{
				Match: []string{"streetAddress"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					StreetAddress: []string{"Apadana"},
				},
			},
			&x509.CertificateRequest{
				Subject: pkix.Name{},
			},
			false,
		},
		{
			"StreetAddressNotMatched",
			Policy{
				Match: []string{"street_address"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					StreetAddress: []string{"Apadana"},
				},
			},
			&x509.CertificateRequest{
				Subject: pkix.Name{
					StreetAddress: []string{},
				},
			},
			false,
		},
		{
			"PostalCodeNotMatched",
			Policy{
				Match: []string{"PostalCode"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					PostalCode: []string{"K1Z"},
				},
			},
			&x509.CertificateRequest{
				Subject: pkix.Name{
					PostalCode: []string{"N2L"},
				},
			},
			false,
		},
		{
			"PostalCodeNotMatched",
			Policy{
				Match: []string{"postalCode"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					PostalCode: []string{"K1Z"},
				},
			},
			&x509.CertificateRequest{
				Subject: pkix.Name{},
			},
			false,
		},
		{
			"PostalCodeNotMatched",
			Policy{
				Match: []string{"postal_code"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					PostalCode: []string{"K1Z"},
				},
			},
			&x509.CertificateRequest{
				Subject: pkix.Name{
					PostalCode: []string{},
				},
			},
			false,
		},

		{
			"PolicyCACSREmpty",
			Policy{},
			&x509.Certificate{},
			&x509.CertificateRequest{},
			true,
		},
		{
			"PolicyCACSRProvided",
			Policy{
				Match:    []string{"Country", "Province", "Locality", "Organization", "DNSName"},
				Supplied: []string{"CommonName", "OrganizationalUnit", "IPAddress", "EmailAddress"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					CommonName:   "Intermediate CA",
					Country:      []string{"CA", "US"},
					Province:     []string{"Ontario", "California"},
					Locality:     []string{"Ottawa", "San Francisco"},
					Organization: []string{"Milad"},
				},
				DNSNames: []string{"example.com"},
			},
			&x509.CertificateRequest{
				Subject: pkix.Name{
					CommonName:         "example.com",
					Country:            []string{"CA", "US"},
					Province:           []string{"Ontario", "California"},
					Locality:           []string{"Ottawa", "San Francisco"},
					Organization:       []string{"Milad"},
					OrganizationalUnit: []string{"R&D", "SRE", "IT"},
				},
				DNSNames:       []string{"example.com"},
				IPAddresses:    []net.IP{net.ParseIP("127.0.0.1")},
				EmailAddresses: []string{"milad@example.com"},
			},
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			trust := PolicyTrustFunc(test.policy)
			result := trust(test.ca, test.csr)

			assert.Equal(t, test.expectedResult, result)
		})
	}
}
