# Taskfleet gRPC API

This module contains the gRPC/Protobuf definitions for the internal communication of all Taskfleet
services. The module is structured as follows:

- There exists one directory for each service which provides an API or publishes messages to Kafka
- For each service that provides an API (i.e. runs a gRPC server), contents are versioned
  - For each version (e.g. `v1`), no breaking changes can be introduced
  - If breaking changes are required, a new version will need to be added
- When a service publishes messages to Kafka, there exists a dedicated folder `messages`
  - Similarly to the gRPC API, the `messages` folder is versioned to prevent breaking changes in
    Kafka messages
  - Within each version folder, each message is defined in its own `.proto` file
