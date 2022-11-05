package kafka

import (
	"os"

	"github.com/google/uuid"
	"go.taskfleet.io/packages/jack"
	"go.uber.org/zap"
)

var logger = func() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.Level.SetLevel(zap.InfoLevel)
	// logger := jack.Must(zap.NewDevelopment())
	return jack.Must(config.Build())
}()

var client = func() *Client {
	clientID := uuid.NewString()
	client := jack.Must(NewClient(clientID, []string{os.Getenv("KAFKA_BOOTSTRAP_SERVER")}, logger))
	return client
}()

var adminClient = func() *AdminClient {
	return jack.Must(client.Admin())
}()
