package kafka

import (
	"os"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var logger = func() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger
}()

var client = func() *Client {
	clientID := uuid.NewString()
	client, err := NewClient(clientID, []string{os.Getenv("KAFKA_BOOTSTRAP_SERVER")}, logger)
	if err != nil {
		panic(err)
	}
	return client
}()

var adminClient = func() *AdminClient {
	admin, err := client.Admin()
	if err != nil {
		panic(err)
	}
	return admin
}()
