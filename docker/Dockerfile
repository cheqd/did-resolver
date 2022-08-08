#####################################################################
###    STAGE 1: Build cheqd did-resolver binary pre-requisites    ###
#####################################################################

FROM golang:1.17-alpine AS builder

WORKDIR /resolver
COPY . ./

# Build did-resolver binary
RUN go mod download && go build -o cheqd-did-resolver main.go

#####################################################################
###    STAGE 2: Build cheqd did-resolver container image          ###
#####################################################################

FROM alpine:3.16 AS resolver

LABEL org.opencontainers.image.description="cheqd DID Resolver"
LABEL org.opencontainers.image.source="https://github.com/cheqd/did-resolver"

# Install pre-requisites
RUN apk update && apk add --no-cache bash ca-certificates

# Set working directory & bash defaults
WORKDIR /resolver
SHELL ["/bin/bash", "-euo", "pipefail", "-c"]

# Copy compiled cheqd-did-resolver binary from Stage 1
COPY --from=builder /resolver/cheqd-did-resolver /bin/cheqd-did-resolver

# Copy base config.yaml
COPY config.yaml .

EXPOSE 8080

ENTRYPOINT ["cheqd-did-resolver serve"]
