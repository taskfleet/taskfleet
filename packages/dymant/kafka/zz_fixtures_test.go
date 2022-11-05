package kafka

import (
	"os"

	"github.com/google/uuid"
	"go.taskfleet.io/packages/jack"
	"go.uber.org/zap"
)

var logger = func() *zap.Logger {
	logger := jack.Must(zap.NewDevelopment())
	return logger
}()

var client = func() *Client {
	clientID := uuid.NewString()
	client := jack.Must(NewClient(clientID, []string{os.Getenv("KAFKA_BOOTSTRAP_SERVER")}, logger))
	return client
}()

var adminClient = func() *AdminClient {
	admin := jack.Must(client.Admin())
	return admin
}()
