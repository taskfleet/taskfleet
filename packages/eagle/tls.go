package eagle

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
)

//-------------------------------------------------------------------------------------------------
// SERVER
//-------------------------------------------------------------------------------------------------

// ServerTLS allows to read server TLS configuration. The certificate and certificate key must
// always be provided to establish a TLS connection. If the CA certificate is provided, it is
// assumed that clients must authenticate themselves (mTLS connection).
type ServerTLS struct {
	Certificate    *String `json:"certificate"`
	CertificateKey *String `json:"certificateKey"`
	CACertificate  *String `json:"caCertificate"`
}

// Config returns the TLS configuration that corresponds to the loaded configuration. If no
// certificate or certificate key is set, a nil configuration will be returned (but no error).
func (t ServerTLS) Config() (*tls.Config, error) {
	// Check if certificate and key are present
	if t.Certificate == nil || t.CertificateKey == nil {
		return nil, nil
	}

	// If so, load certificate
	certificate, err := tls.X509KeyPair(
		[]byte(t.Certificate.Value()), []byte(t.CertificateKey.Value()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS certificate: %s", err)
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}

	// If the CA is set, enable client authentication
	if t.CACertificate != nil {
		pool, err := caPool(t.CACertificate.Value())
		if err != nil {
			return nil, fmt.Errorf("failed to load CA certificate: %s", err)
		}
		config.ClientAuth = tls.RequireAndVerifyClientCert
		config.ClientCAs = pool
	} else {
		config.ClientAuth = tls.NoClientCert
	}

	return config, nil
}

//-------------------------------------------------------------------------------------------------
// CLIENT
//-------------------------------------------------------------------------------------------------

// ClientTLS allows to read client TLS configuration. The CA certificate should be set whenever the
// server does not use a known CA. The client certificate and key should be provided when the
// client should be configured for mTLS. The server name must be set if the server certificate does
// not provide a DNS alt name that corresponds to the target host of the connection.
type ClientTLS struct {
	CACertificate        *String `json:"caCertificate"`
	ClientCertificate    *String `json:"clientCertificate"`
	ClientCertificateKey *String `json:"clientCertificateKey"`
	ServerName           *string `json:"serverName"`
}

// Config returns the TLS configuration that corresponds to the loaded configuration. If no CA
// certificate is set, a `nil` configuration is returned, but no error.
func (t ClientTLS) Config() (*tls.Config, error) {
	// Check if CA is provided
	if t.CACertificate == nil {
		return nil, nil
	}

	// Create config with root CA
	pool, err := caPool(t.CACertificate.Value())
	if err != nil {
		return nil, fmt.Errorf("failed to load CA certificate: %s", err)
	}
	config := &tls.Config{
		RootCAs: pool,
	}

	// Optionally append client certificate
	if t.ClientCertificate != nil && t.ClientCertificateKey != nil {
		certificate, err := tls.X509KeyPair(
			[]byte(t.ClientCertificate.Value()),
			[]byte(t.ClientCertificateKey.Value()),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS certificate: %s", err)
		}
		config.Certificates = []tls.Certificate{certificate}
	}

	// Set additional values
	if t.ServerName != nil {
		config.ServerName = *t.ServerName
	}
	return config, nil
}

//-------------------------------------------------------------------------------------------------
// UTILS
//-------------------------------------------------------------------------------------------------

func caPool(pem string) (*x509.CertPool, error) {
	// Read certificate and construct pool
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM([]byte(pem)) {
		return nil, fmt.Errorf("failed to add CA certificate to cert pool")
	}
	return pool, nil
}
