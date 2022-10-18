# Go Package Dymant

Dymant is a high-level library for interacting with message queues of PubSub systems. Dymant
requires that messages are serialized with Protobuf and provides a common interface that makes it
simple to swap out the actual underlying PubSub system.

Currently, Dymant supports the following message queue implementations:

- Apache Kafka (based on [confluent-kafka-go](https://github.com/confluentinc/confluent-kafka-go))
- Native Go Channels (should only be used for testing)

## Installation

```bash
go get go.taskfleet.io/packages/dymant
```
