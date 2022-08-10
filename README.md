# cheqd DID Resolver

[![GitHub Workflow Status (branch)](https://img.shields.io/github/workflow/status/cheqd/did-resolver/Workflow%20Dispatch/main?label=workflows&style=flat-square)](https://github.com/cheqd/did-resolver/actions/workflows/dispatch.yml)

[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/cheqd/did-resolver?color=green&label=stable&sort=semver&style=flat-square)](https://github.com/cheqd/did-resolver/releases/latest) ![GitHub Release Date](https://img.shields.io/github/release-date/cheqd/did-resolver?style=flat-square)

[![GitHub release (latest SemVer including pre-releases)](https://img.shields.io/github/v/release/cheqd/did-resolver?include_prereleases&label=latest%20%28incl.%20pre-release%29&sort=semver&style=flat-square)](https://github.com/cheqd/did-resolver/releases/) ![GitHub commits since latest release (by date)](https://img.shields.io/github/commits-since/cheqd/did-resolver/latest?style=flat-square)

[![GitHub contributors](https://img.shields.io/github/contributors/cheqd/did-resolver?style=flat-square)](https://github.com/cheqd/did-resolver/graphs/contributors) ![GitHub repo size](https://img.shields.io/github/repo-size/cheqd/did-resolver?style=flat-square)

## ℹ️ Overview

DID methods are expected to provide [standards-compliant methods of DID and DID Document ("DIDDoc") production](https://w3c.github.io/did-core/#production-and-consumption).

The **cheqd DID Resolver** is designed to implement the [W3C DID *Resolution* specification](https://w3c-ccg.github.io/did-resolution/) for [`did:cheqd`](https://docs.cheqd.io/node/architecture/adr-list/adr-002-cheqd-did-method) method.

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
3. `RESOLVER_HOME_DIR` (default: `/resolver`): Default config directory inside the Docker container

#### Configure resolver settings

To configure the resolver, edit the `config.yaml` file in the root of the `@cheqd/did-resolver` repository. The values that can be edited are as follows:

1. **`ledger`**
   1. `networks`: A string specifying the Cosmos SDK gRPC endpoint from which the Resolver pulls data. Format: `<did-namespace>=<resource_url>:<resource_port>;...`
   2. `useTls`: Specify whether gRPC connection to ledger should use secure or insecure pulls. Default is `true` since gRPC uses HTTP/2 with TLS as the transport mechanism.
   3. `timeout`: Timeout (in seconds) to wait for before any ledger requests are considered to have time out.
2. **`resolver`**
   1. `method`: A string describing DID Method that the resolver uses. Set to `cheqd`.
3. **`api`**
   1. `listener`: A string with address and port where the resolver listens for requests from clients.
   2. `resolverPath`: A string with path for DID Doc resolution. Example: `/1.0/identifiers/:did` for requests like `/1.0/identifiers/did:cheqd:testnet:MTMxDQKMTMxDQKMT`
4. **`logLevel`**: `debug`/`warn`/`info`/`error` - to define the application log level

#### Example `config.yaml` file

```yaml
ledger:
  # Provide gRPC endpoints for one of more networks "namespaces" to fetch DIDs/DIDDocs from
  # Must be in format "namespace=endpoint:port"
  networks: "mainnet=grpc.cheqd.net:443;testnet=grpc.cheqd.network:443"
  # Specify whether gRPC connection to ledger should use secure or insecure pulls
  # default: true
  useTls: true
  # Specify a timeout value
  timeout: "5s"

resolver:
  method: "cheqd"

api:
  listener: "0.0.0.0:1313"
  resolverPath: "/1.0/identifiers/:did"

logLevel: "warn"
```

## 📖 Documentation

Further documentation on [cheqd DID Resolver](https://docs.cheqd.io/identity/on-ledger-identity/did-resolver) is available on the [cheqd Identity Documentation site](https://docs.cheqd.io/identity/). This includes instructions on how to do custom builds using `Dockerfile` / Docker Compose.

## 💬 Community

The [**cheqd Community Slack**](http://cheqd.link/join-cheqd-slack) is our primary chat channel for the open-source community, software developers, and node operators.

Please reach out to us there for discussions, help, and feedback on the project.

## 🙋 Find us elsewhere

[![Telegram](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/cheqd) [![Discord](https://img.shields.io/badge/Discord-7289DA?style=for-the-badge&logo=discord&logoColor=white)](http://cheqd.link/discord-github) [![Twitter](https://img.shields.io/badge/Twitter-1DA1F2?style=for-the-badge&logo=twitter&logoColor=white)](https://twitter.com/intent/follow?screen_name=cheqd_io) [![LinkedIn](https://img.shields.io/badge/LinkedIn-0077B5?style=for-the-badge&logo=linkedin&logoColor=white)](http://cheqd.link/linkedin) [![Slack](https://img.shields.io/badge/Slack-4A154B?style=for-the-badge&logo=slack&logoColor=white)](http://cheqd.link/join-cheqd-slack) [![Medium](https://img.shields.io/badge/Medium-12100E?style=for-the-badge&logo=medium&logoColor=white)](https://blog.cheqd.io) [![YouTube](https://img.shields.io/badge/YouTube-FF0000?style=for-the-badge&logo=youtube&logoColor=white)](https://www.youtube.com/channel/UCBUGvvH6t3BAYo5u41hJPzw/)
