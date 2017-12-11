package pki

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPolicyTrustFunc(t *testing.T) {
	tests := []struct {
		policy         Policy
		ca             *x509.Certificate
		csr            *x509.CertificateRequest
		expectedResult bool
	}{
		{
			Policy{},
			nil,
			nil,
			false,
		},
		{
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
			Policy{
				Supplied: []string{"EmailAddresses"},
			},
			&x509.Certificate{},
			&x509.CertificateRequest{},
			false,
		},
		{
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
			Policy{},
			&x509.Certificate{},
			&x509.CertificateRequest{},
			true,
		},
		{
			Policy{
				Match:    []string{"Country", "Province", "Locality", "Organization"},
				Supplied: []string{"CommonName", "OrganizationalUnit", "EmailAddress"},
			},
			&x509.Certificate{
				Subject: pkix.Name{
					CommonName:   "Intermediate CA",
					Country:      []string{"CA", "US"},
					Province:     []string{"Ontario", "California"},
					Locality:     []string{"Ottawa", "San Francisco"},
					Organization: []string{"Milad"},
				},
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
				EmailAddresses: []string{"milad@example.com"},
			},
			true,
		},
	}

	for _, test := range tests {
		trust := PolicyTrustFunc(test.policy)
		result := trust(test.ca, test.csr)

		assert.Equal(t, test.expectedResult, result)
	}
}
