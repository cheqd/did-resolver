# cheqd DID Resolver

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/cheqd/did-resolver?color=green&label=stable%20release&style=flat-square)](https://github.com/cheqd/did-resolver/releases/latest) ![GitHub Release Date](https://img.shields.io/github/release-date/cheqd/did-resolver?color=green&style=flat-square) [![GitHub license](https://img.shields.io/github/license/cheqd/did-resolver?color=blue&style=flat-square)](https://github.com/cheqd/did-resolver/blob/main/LICENSE)

[![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/cheqd/did-resolver?include_prereleases&label=dev%20release&style=flat-square)](https://github.com/cheqd/did-resolver/releases/) ![GitHub commits since latest release (by date)](https://img.shields.io/github/commits-since/cheqd/did-resolver/latest?style=flat-square) [![GitHub contributors](https://img.shields.io/github/contributors/cheqd/did-resolver?label=contributors%20%E2%9D%A4%EF%B8%8F&style=flat-square)](https://github.com/cheqd/did-resolver/graphs/contributors)

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/cheqd/did-resolver/Workflow%20Dispatch?label=workflows&style=flat-square)](https://github.com/cheqd/did-resolver/actions/workflows/dispatch.yml) [![GitHub Workflow Status](https://img.shields.io/github/workflow/status/cheqd/did-resolver/CodeQL?label=CodeQL&style=flat-square)](https://github.com/cheqd/did-resolver/actions/workflows/codeql.yml) ![GitHub repo size](https://img.shields.io/github/repo-size/cheqd/did-resolver?style=flat-square)

## ℹ️ Overview

DID methods are expected to provide [standards-compliant methods of DID and DID Document ("DIDDoc") production](https://w3c.github.io/did-core/#production-and-consumption). The **cheqd DID Resolver** is designed to implement the [W3C DID *Resolution* specification](https://w3c-ccg.github.io/did-resolution/) for [`did:cheqd`](https://docs.cheqd.io/node/architecture/adr-list/adr-002-cheqd-did-method) method.

### 📝 Architecture

The [Architecture Decision Record for the cheqd DID Resolver](https://docs.cheqd.io/identity/architecture/adr-list/adr-001-did-resolver) describes the architecture & design decisions for this software package.

## ✅ Quick Start

If you do not want to install anything and just want to resolve a `did:cheqd` entry from the ledger, you can [load the REST API endpoint for resolver.cheqd.net in your browser](https://resolver.cheqd.net/1.0/identifiers/did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY).

Or, make a request from terminal to this hosted REST API:

```bash
curl -X GET https://resolver.cheqd.net/1.0/identifiers/did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY
```

### Running your own cheqd DID Resolver using Docker

#### Docker Compose command

Spinning up a Docker container from the [pre-built `did-resolver` Docker image on Github](https://github.com/cheqd/did-resolver/pkgs/container/did-resolver) is as simple as the command below:

```bash
docker-compose -f docker/docker-compose.yml --env-file docker/docker-compose.env up --no-build
```

#### Docker Compose environment variable configuration

Environment variables needed in Docker Compose are defined in the `docker/docker-compose.env` file. There are defaults already specified, but you can edit these.

1. `IMAGE_VERSION` (default: `latest`): Version number / tag of the Docker image to run. By default, this is set to pull images from Github Container Registry.
2. `RESOLVER_PORT` (default: `8080`): Port on which the Resolver service is published/exposed on the host machine. Internally mapped to port 8080 in the container.

#### Configure resolver settings

To configure the resolver, use environment variables or edit them in the `container.env` file in the root of the `@cheqd/did-resolver` repository. The values that can be edited are as follows:

1. **`MAINNET_ENDPOINT`** : Mainnet Network endpoint as string with the following format" `<networks>,<useTls>,<timeout>`. Example: `grpc.cheqd.net:443,true,5s`
   1. `networks`: A string specifying the Cosmos SDK gRPC endpoint from which the Resolver pulls data. Format: `<resource_url>:<resource_port>`
   2. `useTls`: Specify whether gRPC connection to ledger should use secure or insecure pulls. Default is `true` since gRPC uses HTTP/2 with TLS as the transport mechanism.
   3. `timeout`: Timeout (in seconds) to wait for before any ledger requests are considered to have time out.
2. **`TESTNET_ENDPOINT`** : Testnet Network endpoint as string with the following format" `<networks>,<useTls>,<timeout>`. Example: `grpc.cheqd.network:443,true,5s`
3. **`RESOLVER_LISTNER`**`: A string with address and port where the resolver listens for requests from clients.
4. **`LOG_LEVEL`**: `debug`/`warn`/`info`/`error` - to define the application log level

#### Example `container.env` file

```bash
# Syntax: <grpc-endpoint-url:port>,boolean,time
# 1st parameter is gRPC endpoint
# 2nd (Boolean) parameter is whether to use TLS or not
# 3nd connection timeout 
MAINNET_ENDPOINT=grpc.cheqd.net:443,true,5s
TESTNET_ENDPOINT=grpc.cheqd.network:443,true,5s

# Sets Echo logging level
LOG_LEVEL="warn"

RESOLVER_LISTNER="0.0.0.0:8080"
```

## 📖 Documentation

Further documentation on [cheqd DID Resolver](https://docs.cheqd.io/identity/on-ledger-identity/did-resolver) is available on the [cheqd Identity Documentation site](https://docs.cheqd.io/identity/). This includes instructions on how to do custom builds using `Dockerfile` / Docker Compose.

## 💬 Community

The [**cheqd Community Slack**](http://cheqd.link/join-cheqd-slack) is our primary chat channel for the open-source community, software developers, and node operators.

Please reach out to us there for discussions, help, and feedback on the project.

## 🙋 Find us elsewhere

[![Telegram](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/cheqd) [![Discord](https://img.shields.io/badge/Discord-7289DA?style=for-the-badge&logo=discord&logoColor=white)](http://cheqd.link/discord-github) [![Twitter](https://img.shields.io/badge/Twitter-1DA1F2?style=for-the-badge&logo=twitter&logoColor=white)](https://twitter.com/intent/follow?screen_name=cheqd_io) [![LinkedIn](https://img.shields.io/badge/LinkedIn-0077B5?style=for-the-badge&logo=linkedin&logoColor=white)](http://cheqd.link/linkedin) [![Slack](https://img.shields.io/badge/Slack-4A154B?style=for-the-badge&logo=slack&logoColor=white)](http://cheqd.link/join-cheqd-slack) [![Medium](https://img.shields.io/badge/Medium-12100E?style=for-the-badge&logo=medium&logoColor=white)](https://blog.cheqd.io) [![YouTube](https://img.shields.io/badge/YouTube-FF0000?style=for-the-badge&logo=youtube&logoColor=white)](https://www.youtube.com/channel/UCBUGvvH6t3BAYo5u41hJPzw/)
