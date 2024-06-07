# Syncer

Syncer is an enhancer service designed to provide advanced features for various database engines including MySQL, WeSQL, PostgreSQL, and MongoDB. It operates as a sidecar process alongside your database process, supercharging it with additional functionality.

## 1. Features

Syncer offers advanced features such as:
- Automatic failover - Promote new leader if the current leader fails
- Switchover - Manual controlled leader election
- Data replication - Keep standbys in sync with leader

These features are exposed and can be controlled via an HTTP API.

## 2. Getting Started

To get started with Syncer, you can choose one of the following methods:

- Download the appropriate release for your platform from [our GitHub repository](https://github.com/apecloud/syncer/releases).
- Use our pre-built Docker image: `docker run -it apecloud/syncer`.
- Build Syncer from source.

### 2.1 Prerequisites
- Go 1.20+ (for generics support)
- Database (MySQL, PostgreSQL, etc)

### 2.2 Building from Source

To build Syncer from source, you need Go 1.20 or later, which includes support for generics programming. Check the [Go Installation Guide](https://golang.org/doc/install) for instructions on how to install Go.

Use the `go build` command to build Syncer and produce the binary file. The executable will be in your current directory.

```shell
$ cd syncer/cmd/syncer
$ go build -o syncer main.go
```

### 2.3 Configuration
Syncer retrieves configuration settings from environment variables. Currently, these include:

- KB_SERVICE_PORT: The service port for DB service (e.g., 3306).
- KB_SERVICE_USER: The username used to connect to the service (e.g., root).
- KB_SERVICE_PASSWORD: The password used to connect to the service.
- KB_SERVICE_TYPE: The type of DB Service(e.g., mysql/postgresql).
- KB_CLUSTER_COMP_NAME: The cluster and component name for the db service.
- KB_CLUSTER_NAME: The cluster name of db service.
- KB_COMP_NAME: The namme of the database component.
- KB_POD_NAME: The Kubernetes pod name that the database is running in, if applicable. Used for Kubernetes deployments.
- KB_NAMESPACE: The Kubernetes namespace the pod belongs to, if running in Kubernetes.

### 2.4 Running Syncer
Once you've built Syncer, you can start it with the following command:

```shell
$ syncer --config-path config/components/ --port 3601
```

## 3. License
Syncer is licensed under the AGPL 3.0 license. See the LICENSE file in this repository for more details.
