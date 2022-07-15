# did:cheqd DID Resolver

## ℹ️ Overview

cheqd DID resovler offers multiple implementations for resolving cheqd DIDs, according to the [cheqd DID method](https://docs.cheqd.io/node/architecture/adr-list/adr-002-cheqd-did-method#:~:text=Summary,on%20the%20Cosmos%20blockchain%20framework.)

This resolver aims to make it easy for third parties to resolve cheqd DIDs, using either a full DID resolver, a proxy DID resolver or through the [Universal Resolver](https://dev.uniresolver.io/).

## Example DIDs

```commandline
    did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY
    did:cheqd:testnet:MTMxDQKMTMxDQKMT
```

## Quick Start

If you do not want to install anything, but just want to resolve a DID right now, then make a request in the browser https://resolver.cheqd.net/1.0/identifiers/did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY:

or through the terminal:

```bash
curl -X GET https://resolver.cheqd.net/1.0/identifiers/did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY
```

### Response example

```json
    {
       "didResolutionMetadata":{
          "contentType":"application/did+ld+json",
          "retrieved":"2022-07-15T14:55:16Z",
          "did":{
             "didString":"did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY",
             "methodSpecificId":"zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY",
             "method":"cheqd"
          }
       },
       "didDocument":{
          "@context":[
             "https://www.w3.org/ns/did/v1"
          ],
          "id":"did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY",
          "verificationMethod":[
             {
                "id":"did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY#key1",
                "type":"Ed25519VerificationKey2020",
                "controller":"did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY",
                "publicKeyMultibase":"zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkYsWCo7fztHtepn"
             }
          ],
          "authentication":[
             "did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY#key1"
          ],
          "service":[
             {
                "id":"did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY#website",
                "type":"LinkedDomains",
                "serviceEndpoint":"https://www.cheqd.io"
             },
             {
                "id":"did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY#non-fungible-image",
                "type":"LinkedDomains",
                "serviceEndpoint":"https://gateway.ipfs.io/ipfs/bafybeihetj2ng3d74k7t754atv2s5dk76pcqtvxls6dntef3xa6rax25xe"
             },
             {
                "id":"did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY#twitter",
                "type":"LinkedDomains",
                "serviceEndpoint":"https://twitter.com/cheqd_io"
             },
             {
                "id":"did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY#linkedin",
                "type":"LinkedDomains",
                "serviceEndpoint":"https://www.linkedin.com/company/cheqd-identity/"
             }
          ]
       },
       "didDocumentMetadata":{
          "created":"2022-04-05T11:49:19Z",
          "versionId":"EDEAD35C83E20A72872ACD3C36B7BA42300712FC8E3EEE1340E47E2F1B216B2D"
       }
    }
```

[Read more about cheqd DID resolver features](https://github.com/cheqd/identity-docs/blob/main/tutorials/resolver/using-cheqd-universal-resolver-driver.md)

## 🧑‍💻🛠 Developer Guide

## Full DID Resolver

For starting Full DID Resolver use

```bash
docker compose --profile full up --build
```

After you can check if it works

```bash
curl -X GET https://resolver.cheqd.net/1.0/identifiers/did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY
```

Response should be the same with [this example](#response-example)

[Read more about Full DID Resolver configuration](https://github.com/cheqd/identity-docs/blob/main/tutorials/resolver/using-full-cheqd-did-resolver.md)

## Proxy DID Resolver

You can use light profile:

```commandline
docker compose --profile light up --build
```

for having an opportunity to make a localhost requests

```commandline
curl -X GET http://localhost:8080/1.0/identifiers/did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY
```

with [this kind]((#response-example)) of responses redirected from https://resolver.cheqd.net

[More about light DID Resolver](https://github.com/cheqd/identity-docs/blob/main/tutorials/resolver/using-light-cheqd-did-resolver.md)

## Universal Resolver

The Universal Resolver wraps an API around a number of co-located Docker containers running DID-method-specific drivers.

Integration phase: in progress

## 🙋 Find us elsewhere

[![Telegram](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/cheqd) [![Discord](https://img.shields.io/badge/Discord-7289DA?style=for-the-badge&logo=discord&logoColor=white)](http://cheqd.link/discord-github) [![Twitter](https://img.shields.io/badge/Twitter-1DA1F2?style=for-the-badge&logo=twitter&logoColor=white)](https://twitter.com/intent/follow?screen_name=cheqd_io) [![LinkedIn](https://img.shields.io/badge/LinkedIn-0077B5?style=for-the-badge&logo=linkedin&logoColor=white)](http://cheqd.link/linkedin) [![Slack](https://img.shields.io/badge/Slack-4A154B?style=for-the-badge&logo=slack&logoColor=white)](http://cheqd.link/join-cheqd-slack) [![Medium](https://img.shields.io/badge/Medium-12100E?style=for-the-badge&logo=medium&logoColor=white)](https://blog.cheqd.io) [![YouTube](https://img.shields.io/badge/YouTube-FF0000?style=for-the-badge&logo=youtube&logoColor=white)](https://www.youtube.com/channel/UCBUGvvH6t3BAYo5u41hJPzw/)
