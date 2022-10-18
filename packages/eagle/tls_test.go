package eagle

import (
	"crypto/tls"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerTLSEmpty(t *testing.T) {
	data := `{}`
	var serverTLS ServerTLS
	err := json.Unmarshal([]byte(data), &serverTLS)
	require.Nil(t, err)

	config, err := serverTLS.Config()
	require.Nil(t, err)
	assert.Nil(t, config)
}

func TestServerTLSCertificate(t *testing.T) {
	data := `{
		"certificate": {"file": "testdata/server.crt"},
		"certificateKey": {"file": "testdata/server.key"}
	}`
	var serverTLS ServerTLS
	err := json.Unmarshal([]byte(data), &serverTLS)
	require.Nil(t, err)

	config, err := serverTLS.Config()
	require.Nil(t, err)
	assert.Len(t, config.Certificates, 1)
	assert.Equal(t, config.ClientAuth, tls.NoClientCert)
}

func TestServerMTLS(t *testing.T) {
	data := `{
		"certificate": {"file": "testdata/server.crt"},
		"certificateKey": {"file": "testdata/server.key"},
		"caCertificate": {"file": "testdata/ca.crt"}
	}`
	var serverTLS ServerTLS
	err := json.Unmarshal([]byte(data), &serverTLS)
	require.Nil(t, err)

	config, err := serverTLS.Config()
	require.Nil(t, err)
	assert.Len(t, config.Certificates, 1)
	assert.Len(t, config.ClientCAs.Subjects(), 1)
	assert.Equal(t, config.ClientAuth, tls.RequireAndVerifyClientCert)
}

func TestClientTLSEmpty(t *testing.T) {
	data := `{}`
	var clientTLS ClientTLS
	err := json.Unmarshal([]byte(data), &clientTLS)
	require.Nil(t, err)

	config, err := clientTLS.Config()
	require.Nil(t, err)
	assert.Nil(t, config)
}

func TestClientTLS(t *testing.T) {
	data := `{
		"caCertificate": {"file": "testdata/ca.crt"}
	}`
	var clientTLS ClientTLS
	err := json.Unmarshal([]byte(data), &clientTLS)
	require.Nil(t, err)

	config, err := clientTLS.Config()
	require.Nil(t, err)
	assert.Len(t, config.RootCAs.Subjects(), 1)
}

func TestClientMTLS(t *testing.T) {
	data := `{
		"caCertificate": {"file": "testdata/ca.crt"},
		"clientCertificate": {"file": "testdata/client.crt"},
		"clientCertificateKey": {"file": "testdata/client.key"},
		"serverName": "test"
	}`
	var clientTLS ClientTLS
	err := json.Unmarshal([]byte(data), &clientTLS)
	require.Nil(t, err)

	config, err := clientTLS.Config()
	require.Nil(t, err)
	assert.Len(t, config.RootCAs.Subjects(), 1)
	assert.Len(t, config.Certificates, 1)
	assert.Equal(t, "test", config.ServerName)
}
