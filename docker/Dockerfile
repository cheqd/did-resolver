###############################################################
###    STAGE 1: Build cheqd-did-resolver binary pre-requisites    ###
###############################################################

FROM golang:1.17.8-buster as builder

# Set user directory and details
ENV CHEQD_RESOLVER_HOME_DIR="/home/cheqd-resolver"
ARG UID=1000
ARG GID=1000

# Add cheqd user to use in the container
RUN groupadd --system --gid $GID cheqd-resolver \
    && useradd --system --create-home --home-dir ${CHEQD_RESOLVER_HOME_DIR} --shell /bin/bash --gid cheqd-resolver --uid $UID cheqd-resolver

WORKDIR ${CHEQD_RESOLVER_HOME_DIR}
USER cheqd-resolver

COPY types ./types
COPY services ./services
COPY go.mod .
COPY go.sum .
COPY main.go .

# Make cheqd-did-resolver binary
RUN go build -o cheqd-did-resolver main.go

###############################################################
###           STAGE 2: Build cheqd-did-resolver runner            ###
###############################################################

FROM ubuntu:focal AS runner
LABEL org.opencontainers.image.description "Cheqd DID-Resolver runner"
LABEL org.opencontainers.image.source "https://github.com/cheqd/cheqd-did-resolver"
ENV CHEQD_RESOLVER_HOME_DIR="/home/cheqd-resolver"

# Copy compiled cheqd-did-resolver binary from Stage 1
COPY --from=builder ${CHEQD_RESOLVER_HOME_DIR} /bin

# Copy base config.yml
WORKDIR ${CHEQD_RESOLVER_HOME_DIR}

EXPOSE 1313
ENTRYPOINT ["cheqd-did-resolver"]