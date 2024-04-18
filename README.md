# cheqd DID Resolver

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/cheqd/did-resolver?color=green&label=stable%20release&style=flat-square)](https://github.com/cheqd/did-resolver/releases/latest) ![GitHub Release Date](https://img.shields.io/github/release-date/cheqd/did-resolver?color=green&style=flat-square) [![GitHub license](https://img.shields.io/github/license/cheqd/did-resolver?color=blue&style=flat-square)](https://github.com/cheqd/did-resolver/blob/main/LICENSE)

[![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/cheqd/did-resolver?include_prereleases&label=dev%20release&style=flat-square)](https://github.com/cheqd/did-resolver/releases/) ![GitHub commits since latest release (by date)](https://img.shields.io/github/commits-since/cheqd/did-resolver/latest?style=flat-square) [![GitHub contributors](https://img.shields.io/github/contributors/cheqd/did-resolver?label=contributors%20%E2%9D%A4%EF%B8%8F&style=flat-square)](https://github.com/cheqd/did-resolver/graphs/contributors)

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/cheqd/did-resolver/dispatch.yml?label=workflows&style=flat-square)](https://github.com/cheqd/did-resolver/actions/workflows/dispatch.yml) [![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/cheqd/did-resolver/codeql.yml?label=CodeQL&style=flat-square)](https://github.com/cheqd/did-resolver/actions/workflows/codeql.yml) ![GitHub repo size](https://img.shields.io/github/repo-size/cheqd/did-resolver?style=flat-square)

## ℹ️ Overview

DID methods are expected to provide [standards-compliant methods of DID and DID Document ("DIDDoc") production](https://w3c.github.io/did-core/#production-and-consumption). The **cheqd DID Resolver** is designed to implement the [W3C DID *Resolution* specification](https://w3c-ccg.github.io/did-resolution/) for [`did:cheqd`](https://docs.cheqd.io/identity/architecture/adr-list/adr-001-cheqd-did-method) method.

### 📝 Architecture

The [Architecture Decision Record for the cheqd DID Resolver](https://docs.cheqd.io/identity/architecture/adr-list/adr-003-did-resolver) describes the architecture & design decisions for this software package.

## ✅ Quick Start

If you do not want to install anything and just want to resolve a `did:cheqd` entry from the ledger, you can load the REST API endpoint for [resolver.cheqd.net](https://resolver.cheqd.net/) in your browser.

Or, make a request from terminal to this hosted REST API:

```bash
curl -X GET https://resolver.cheqd.net/1.0/identifiers/did:cheqd:testnet:55dbc8bf-fba3-4117-855c-1e0dc1d3bb47
```

## 🛠️ Running your own cheqd DID Resolver using Docker

Spinning up a Docker container from the [pre-built `did-resolver` Docker image on Github](https://github.com/cheqd/did-resolver/pkgs/container/did-resolver) is as simple as the command below:

```bash
docker compose -f docker/docker-compose.yml up --detach
```

### Configure resolver settings

To configure the resolver, modify the values under the `environment` section of the [Docker Compose file](https://github.com/cheqd/did-resolver/blob/main/docker/docker-compose.yml). The values that can be edited are as follows:

1. **`MAINNET_ENDPOINT`** : Mainnet Network endpoint as string with the following format" `<networks>,<useTls>,<timeout>`. Example: `grpc.cheqd.net:443,true,5s`
   1. `networks`: A string specifying the Cosmos SDK gRPC endpoint from which the Resolver pulls data. Format: `<resource_url>:<resource_port>`
   2. `useTls`: Specify whether gRPC connection to ledger should use secure or insecure pulls. Default is `true` since gRPC uses HTTP/2 with TLS as the transport mechanism.
   3. `timeout`: Timeout (in seconds) to wait for before any ledger requests are considered to have time out.
2. **`TESTNET_ENDPOINT`** : Testnet Network endpoint as string with the following format" `<networks>,<useTls>,<timeout>`. Example: `grpc.cheqd.network:443,true,5s`
3. **`RESOLVER_LISTENER`**`: A string with address and port where the resolver listens for requests from clients.
4. **`LOG_LEVEL`**: `debug`/`warn`/`info`/`error` - to define the application log level.

#### gRPC Endpoints used by DID Resolver

Our DID Resolver uses the [Cosmos gRPC endpoint](https://docs.cosmos.network/main/core/grpc_rest) from `cheqd-node` to fetch data. Typically, this would be running on port `9090` on a `cheqd-node` instance.

You can either use [public gRPC endpoints for the cheqd network](https://cosmos.directory/cheqd/nodes) (such as the default ones mentioned above), or point it to your own `cheqd-node` instance by enabling gRPC in the `app.toml` configuration file on a node:

```toml
[grpc]

# Enable defines if the gRPC server should be enabled.
enable = true

# Address defines the gRPC server address to bind to.
address = "0.0.0.0:9090"
```

**Note**: If you're pointing a DID Resolver to your own node instance, by default `cheqd-node` instance gRPC endpoints are *not* served up with a TLS certificate. This means the `useTls` property would need to be set to `false`, unless you're otherwise using a load balancer that provides TLS connections to the gRPC port.

## 🧑‍💻 Building your own Docker image

### Using Docker Build

You can build your own image using `docker build`

```bash
docker build --file docker/Dockerfile --target resolver . --tag did-resolver:local
```

### Using Docker Compose Build

Uncomment the `build` section in the `docker/docker-compose.yml` file. This relies on the `Dockerfile` above but uses Docker Compose syntax to customise the build:

```yaml
build:
  context: ../
  dockerfile: docker/Dockerfile
  target: resolver
image: did-resolver:local
# image: ghcr.io/cheqd/did-resolver:${IMAGE_VERSION}
```

Make sure you comment out the pre-existing `image` property that pulls in a container image from Github Container Registry, as shown above.

You can also do *just* a build with:

```bash
docker-compose -f docker/docker-compose.yml --env-file docker/docker-compose.env build --no-cache
```

### Running a custom built image

The instructions to configure and run the resolver are the same as when using the pre-built image.

## 🌐 Resolving `did:cheqd` via Universal Resolver

The [resolver.cheqd.net](https://resolver.cheqd.net/) API endpoint is run by the cheqd team and only handles `did:cheqd` credentials.

If you want to resolve DIDs from multiple DID methods, the [Universal Resolver](https://github.com/decentralized-identity/universal-resolver) project provides a multi DID method resolver.

### Using a pre-existing Universal Resolver endpoint

You can make resolution requests to a pre-existing Universal Resolver endpoint, such as [dev.uniresolver.io](https://dev.uniresolver.io), to their REST API endpoint:

```bash
curl -X GET https://resolver.cheqd.net/1.0/identifiers/did:cheqd:testnet:55dbc8bf-fba3-4117-855c-1e0dc1d3bb47
```

### Running your own Universal Resolver instance

You can also run your own instance of Universal Resolver, using the Docker Compose file of the project.

The [Universal Resolver quick start guide](https://github.com/decentralized-identity/universal-resolver#quick-start) provides instructions on how to do this:

```bash
git clone https://github.com/decentralized-identity/universal-resolver
cd universal-resolver/
docker-compose -f docker-compose.yml pull
docker-compose -f docker-compose.yml up
```

## 🧪 Running tests

This repository has two kinds of tests:

1. **Unit tests**: Test coverage for internal functions and utilities
2. **Integration tests**: Test coverage for a built version of this software with the ledger.

### Prerequisites

Tests are written in [Ginkgo](https://onsi.github.io/ginkgo/), a behaviour-driven test framework for Golang.

To install Ginkgo, run the following command:

```bash
go install github.com/onsi/ginkgo/v2/ginkgo@latest
```

### Execute unit tests with Ginkgo

Unit tests can be executed without running a Docker container. To execute all unit tests, use the following command:

```bash
ginkgo -r --tags unit --race --randomize-all --randomize-suites --keep-going --trace
```

The `unit` tag is specified to only run unit tests.

### Execute integration tests with Ginkgo

To run integration tests, it's necessary to have a instance of this DID Resolver running. The easiest way to do this is to use the [Docker Compose file under the `tests` directory](./tests/docker-compose-testing.yml), with a few modifications.

#### If you **don't have** an already-built Docker image

To build a local image for testing, uncomment the `build` section in the Docker Compose file. You can modify the `build` section according to your needs:

```yaml
build:
  context: ../
  dockerfile: docker/Dockerfile
  target: resolver
```

#### If you **do have** an already-built Docker image

Alternatively, if you want to run tests using an already-built image, just change the image name and tag in the `image` line in the Docker Compose file:

```yaml
# image: cheqd/did-resolver:staging-latest
image: cheqd/did-resolver:local
```

#### Run the test Docker container

Use Docker Compose to bring the container up (for both scenarios above):

```bash
docker compose -f tests/docker-compose-testing.yml up --detach
```

#### Execute the integration tests

You can execute the tests as long as you have Ginkgo CLI installed, which targets the tests towards the running Docker container:

```bash
ginkgo -r --tags integration --race --randomize-suites --keep-going --trace
```

**Note**: By default, the tests target `localhost:8080` as the port where it expects the running DID Resolver instance for testing. If your running instance is at a different address, you can override this by setting a value for the `TEST_HOST_ADDRESS` environment variable *before* executing the Ginkgo test suite.

```bash
export TEST_HOST_ADDRESS="where.is.did.resolver.running:port"
```

You can also do this inline without exporting the value:

```bash
TEST_HOST_ADDRESS="where.is.did.resolver.running:port" ginkgo -r --tags integration --race --randomize-suites --keep-going --trace
```

### (Alternative) Executing tests with Github Actions

All tests are run automatically by the [Github Actions `test` workflow](.github/workflows/test.yml) in this repository. Use that workflow for inspiration if you want to customise the test execution; or [Ginkgo documentation on how to customise test execution](https://onsi.github.io/ginkgo/#reporting-and-profiling-suites).

You can also directly [execute the Github Actions locally using `nektos/act`](https://github.com/nektos/act).

## 📖 Documentation

Further documentation on [cheqd DID Resolver](https://docs.cheqd.io/identity/advanced/did-resolver) is available on the [cheqd Identity Documentation site](https://docs.cheqd.io/identity/). This includes instructions on how to do custom builds using `Dockerfile` / Docker Compose.

## 🐞 Bug reports & 🤔 feature requests

If you notice anything not behaving how you expected, or would like to make a suggestion / request for a new feature, please create a [**new issue**](https://github.com/cheqd/did-resolver/issues/new/choose) and let us know.

## 💬 Community

Our [**Discord server**](http://cheqd.link/discord-github) is the primary chat channel for the open-source community, software developers, and node operators.

Please reach out to us there for discussions, help, and feedback on the project.

## 🙋 Find us elsewhere

[![Telegram](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/cheqd) [![Discord](https://img.shields.io/badge/Discord-7289DA?style=for-the-badge&logo=discord&logoColor=white)](http://cheqd.link/discord-github) [![Twitter](https://img.shields.io/badge/Twitter-1DA1F2?style=for-the-badge&logo=twitter&logoColor=white)](https://twitter.com/intent/follow?screen_name=cheqd_io) [![LinkedIn](https://img.shields.io/badge/LinkedIn-0077B5?style=for-the-badge&logo=linkedin&logoColor=white)](http://cheqd.link/linkedin) [![Medium](https://img.shields.io/badge/Medium-12100E?style=for-the-badge&logo=medium&logoColor=white)](https://blog.cheqd.io) [![YouTube](https://img.shields.io/badge/YouTube-FF0000?style=for-the-badge&logo=youtube&logoColor=white)](https://www.youtube.com/channel/UCBUGvvH6t3BAYo5u41hJPzw/)
