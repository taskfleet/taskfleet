package kafka

import (
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

const (
	logKeyComponent = "component"
	logKeyTopic     = "topic"
	logKeyPartition = "partition"
	logKeyOffset    = "offset"
)

func logProduced(l *zap.Logger, msg *kafka.Message) {
	if msg.TopicPartition.Error != nil {
		l.Error("failed to publish message",
			zap.Error(msg.TopicPartition.Error),
		)
	} else {
		l.Debug("successfully published message",
			zap.Int32(logKeyPartition, msg.TopicPartition.Partition),
			zap.Int64(logKeyOffset, int64(msg.TopicPartition.Offset)),
		)
	}
}

func logConsumed(l *zap.Logger, msg *kafka.Message) {
	if msg.TopicPartition.Error != nil {
		l.Error("failed to consume message",
			zap.Error(msg.TopicPartition.Error),
		)
	} else {
		l.Debug("successfully consumed message",
			zap.Int32(logKeyPartition, msg.TopicPartition.Partition),
			zap.Int64(logKeyOffset, int64(msg.TopicPartition.Offset)),
		)
	}
}

func logFieldPartitions(partitions []kafka.TopicPartition) zap.Field {
	p := []string{}
	for _, partition := range partitions {
		p = append(p, fmt.Sprintf("%d", partition.Partition))
	}
	return zap.String("partitions", strings.Join(p, ","))
}

func logFieldsOffsets(partitions []kafka.TopicPartition) []zap.Field {
	fields := []zap.Field{}
	for _, partition := range partitions {
		fields = append(fields, zap.Int64(
			fmt.Sprintf("partition-%d", partition.Partition), int64(partition.Offset),
		))
	}
	return fields
}
