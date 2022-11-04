package cmd

import (
	"go.taskfleet.io/packages/dymant/kafka"
	"go.taskfleet.io/packages/eagle"
)

type KafkaConfig struct {
	kafka.Config `json:",inline"`
	Topic        string `json:"topic"`
}

type MinionConfig struct {
	TLS eagle.ClientTLS `json:"tls"`
}
