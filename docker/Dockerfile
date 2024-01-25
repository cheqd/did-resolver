#####################################################################
###    STAGE 1: Build cheqd did-resolver binary pre-requisites    ###
#####################################################################

FROM golang:1.20-alpine AS builder

# Install minimum necessary dependencies
ENV PACKAGES make git bash linux-headers findutils
RUN apk update && apk add --no-cache $PACKAGES

# Set working directory for the build
WORKDIR /builder

# Add source files
COPY . .

# Build did-resolver binary
RUN make tidy && make build 

#####################################################################
###    STAGE 2: Build cheqd did-resolver container image          ###
#####################################################################

FROM alpine:3 AS resolver

# Install pre-requisites
RUN apk update && apk add --no-cache bash ca-certificates

# Copy did-resolver binary from Stage 1
COPY --from=builder /builder/build/did-resolver /usr/bin/did-resolver

# Set user directory and details
ARG HOME_DIR="/resolver"
ARG USER="resolver"
SHELL ["/bin/bash", "-euo", "pipefail", "-c"]

# Add non-root user to use in the container
RUN addgroup --system $USER \
    && adduser $USER --system --home $HOME_DIR --shell /bin/bash

# Set environment variables with default values
ENV MAINNET_ENDPOINT="grpc.cheqd.net:443,true,5s"
ENV TESTNET_ENDPOINT="grpc.cheqd.network:443,true,5s"
ENV LOG_LEVEL="warn"
ENV RESOLVER_LISTENER="0.0.0.0:8080"

# Set working directory & bash defaults
WORKDIR $HOME_DIR
USER $USER

EXPOSE 8080

ENTRYPOINT ["did-resolver"]
