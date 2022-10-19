package mercury

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.taskfleet.io/packages/eagle"
	"go.taskfleet.io/packages/jack"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/test/bufconn"
)

func TestNoTLS(t *testing.T) {
	err := runTLSTest(t, eagle.ServerTLS{}, insecure.NewCredentials())
	assert.Nil(t, err)
}

func TestServerTLS(t *testing.T) {
	serverTLS := eagle.ServerTLS{
		Certificate:    jack.Ptr(jack.Must(eagle.NewStringFromFile("testdata/server.crt"))),
		CertificateKey: jack.Ptr(jack.Must(eagle.NewStringFromFile("testdata/server.key"))),
	}
	clientTLS := eagle.ClientTLS{
		CACertificate: jack.Ptr(jack.Must(eagle.NewStringFromFile("testdata/ca.crt"))),
	}

	// Insecure connection should fail
	err := runTLSTest(t, serverTLS, insecure.NewCredentials())
	assert.NotNil(t, err)

	// Secure connection should fail with incorrect server name
	clientConfig, err := clientTLS.Config()
	require.Nil(t, err)
	err = runTLSTest(t, serverTLS, credentials.NewTLS(clientConfig))
	assert.NotNil(t, err)

	// Secure connection should succeed with correct server name
	clientTLS.ServerName = jack.Ptr("dymant-server")
	clientConfig, err = clientTLS.Config()
	require.Nil(t, err)
	err = runTLSTest(t, serverTLS, credentials.NewTLS(clientConfig))
	assert.Nil(t, err)
}

func TestMutualTLS(t *testing.T) {
	serverTLS := eagle.ServerTLS{
		Certificate:    jack.Ptr(jack.Must(eagle.NewStringFromFile("testdata/server.crt"))),
		CertificateKey: jack.Ptr(jack.Must(eagle.NewStringFromFile("testdata/server.key"))),
		CACertificate:  jack.Ptr(jack.Must(eagle.NewStringFromFile("testdata/ca.crt"))),
	}
	clientTLS := eagle.ClientTLS{
		CACertificate: jack.Ptr(jack.Must(eagle.NewStringFromFile("testdata/ca.crt"))),
		ServerName:    jack.Ptr("dymant-server"),
	}

	// Insecure connection should fail
	err := runTLSTest(t, serverTLS, insecure.NewCredentials())
	assert.NotNil(t, err)

	// Secure connection without client certificate should fail
	clientConfig, err := clientTLS.Config()
	require.Nil(t, err)
	err = runTLSTest(t, serverTLS, credentials.NewTLS(clientConfig))
	assert.NotNil(t, err)

	// Mutually secure connection with incorrect certificates should fail
	clientTLS.ClientCertificate = jack.Ptr(
		jack.Must(eagle.NewStringFromFile("testdata/server.crt")),
	)
	clientTLS.ClientCertificateKey = jack.Ptr(
		jack.Must(eagle.NewStringFromFile("testdata/server.key")),
	)
	clientConfig, err = clientTLS.Config()
	require.Nil(t, err)
	err = runTLSTest(t, serverTLS, credentials.NewTLS(clientConfig))
	assert.NotNil(t, err)

	// Mutually secure connection with incorrect server name should fail
	clientTLS.ServerName = nil
	clientTLS.ClientCertificate = jack.Ptr(
		jack.Must(eagle.NewStringFromFile("testdata/client.crt")),
	)
	clientTLS.ClientCertificateKey = jack.Ptr(
		jack.Must(eagle.NewStringFromFile("testdata/client.key")),
	)
	clientConfig, err = clientTLS.Config()
	require.Nil(t, err)
	err = runTLSTest(t, serverTLS, credentials.NewTLS(clientConfig))
	assert.NotNil(t, err)

	// Mutually secure connection with correct server name should succeed
	clientTLS.ServerName = jack.Ptr("dymant-server")
	clientConfig, err = clientTLS.Config()
	require.Nil(t, err)
	err = runTLSTest(t, serverTLS, credentials.NewTLS(clientConfig))
	assert.Nil(t, err)
}

//-------------------------------------------------------------------------------------------------

func runTLSTest(
	t *testing.T, tls eagle.ServerTLS, credentials credentials.TransportCredentials,
) error {
	serverTLS, err := tls.Config()
	require.Nil(t, err)

	server, err := NewGrpc(5404,
		WithHealthService(),
		WithTLS(serverTLS),
	)
	require.Nil(t, err)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 250*time.Millisecond)
	defer cancel()

	listener := bufconn.Listen(1024 * 1024)

	dialer := func(ctx context.Context, s string) (net.Conn, error) {
		return listener.DialContext(ctx)
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		defer cancel()
		conn, err := grpc.DialContext(
			ctx,
			"bufnet",
			grpc.WithContextDialer(dialer),
			grpc.WithTransportCredentials(credentials),
		)
		if err != nil {
			return err
		}
		client := grpc_health_v1.NewHealthClient(conn)
		_, err = client.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
		return err
	})
	eg.Go(func() error {
		return server.runBuf(ctx, listener)
	})
	err = eg.Wait()
	if err == context.Canceled {
		return nil
	}
	return err
}

//-------------------------------------------------------------------------------------------------

func (i *Grpc) runBuf(ctx context.Context, listener *bufconn.Listener) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return i.Server.Serve(listener)
	})
	eg.Go(func() error {
		<-ctx.Done()
		i.Server.Stop()
		return ctx.Err()
	})
	err := eg.Wait()
	return err
}
