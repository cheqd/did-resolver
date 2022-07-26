# did:cheqd DID Resolver

## ℹ️ Overview

cheqd DID resovler offers multiple implementations for resolving cheqd DIDs, according to the [cheqd DID method](https://docs.cheqd.io/node/architecture/adr-list/adr-002-cheqd-did-method#:~:text=Summary,on%20the%20Cosmos%20blockchain%20framework.)

This resolver aims to make it easy for third parties to resolve cheqd DIDs, using either a Full DID resolver, a Light DID resolver or through the [Universal Resolver](https://dev.uniresolver.io/).

## Example DIDs

```commandline
    did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY
    did:cheqd:testnet:MTMxDQKMTMxDQKMT
```

## Quick Start

If you do not want to install anything, but just want to resolve a DID, then you can make a request in the browser:

<https://resolver.cheqd.net/1.0/identifiers/did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY>

or through the command terminal:

```bash
curl -X GET https://resolver.cheqd.net/1.0/identifiers/did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY
```

[Read more about cheqd DID resolver features](https://github.com/cheqd/identity-docs/blob/main/tutorials/resolver/using-cheqd-universal-resolver-driver.md)

## 🧑‍💻🛠 Developer Guide

## Full DID Resolver

For initiating the Full DID Resolver, use:

```bash
docker compose --profile full up --build
```

After, you can check if it works:

```bash
curl -X GET https://resolver.cheqd.net/1.0/identifiers/did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY
```

[Read more about Full DID Resolver configuration](https://github.com/cheqd/identity-docs/blob/main/tutorials/resolver/using-full-cheqd-did-resolver.md)

## Light DID Resolver

Status: In development

[More about light DID Resolver](https://github.com/cheqd/identity-docs/blob/main/tutorials/resolver/using-light-cheqd-did-resolver.md)

## Universal Resolver

The [Universal Resolver](https://github.com/decentralized-identity/universal-resolver) wraps an API around a number of co-located Docker containers running DID-method-specific drivers.

For a [quick start](https://github.com/decentralized-identity/universal-resolver#quick-start)

```bash
git clone https://github.com/decentralized-identity/universal-resolver
cd universal-resolver/
docker-compose -f docker-compose.yml pull
docker-compose -f docker-compose.yml up
```

You should then be able to resolve identifiers locally using simple `curl` requests as follow:

```bash
curl -X GET http://localhost:8080/1.0/identifiers/did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY
```

## 🙋 Find us elsewhere

[![Telegram](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/cheqd) [![Discord](https://img.shields.io/badge/Discord-7289DA?style=for-the-badge&logo=discord&logoColor=white)](http://cheqd.link/discord-github) [![Twitter](https://img.shields.io/badge/Twitter-1DA1F2?style=for-the-badge&logo=twitter&logoColor=white)](https://twitter.com/intent/follow?screen_name=cheqd_io) [![LinkedIn](https://img.shields.io/badge/LinkedIn-0077B5?style=for-the-badge&logo=linkedin&logoColor=white)](http://cheqd.link/linkedin) [![Slack](https://img.shields.io/badge/Slack-4A154B?style=for-the-badge&logo=slack&logoColor=white)](http://cheqd.link/join-cheqd-slack) [![Medium](https://img.shields.io/badge/Medium-12100E?style=for-the-badge&logo=medium&logoColor=white)](https://blog.cheqd.io) [![YouTube](https://img.shields.io/badge/YouTube-FF0000?style=for-the-badge&logo=youtube&logoColor=white)](https://www.youtube.com/channel/UCBUGvvH6t3BAYo5u41hJPzw/)
