#####################################################################
###    STAGE 1: Build cheqd-did-resolver binary pre-requisites    ###
#####################################################################

FROM golang:1.17.8-buster as builder

WORKDIR /root

COPY types ./types
COPY services ./services
COPY cmd ./cmd
COPY go.mod .
COPY go.sum .
COPY main.go .

# Make cheqd-did-resolver binary
RUN go build -o cheqd-did-resolver main.go

#####################################################################
###           STAGE 2: Build cheqd-did-resolver runner            ###
#####################################################################

FROM ubuntu:focal AS runner

LABEL org.opencontainers.image.description="Cheqd DID-Resolver runner"
LABEL org.opencontainers.image.source="https://github.com/cheqd/did-resolver"


# Set user details
ARG UID=1000
ARG GID=1000

ARG USER="resolver"
ARG GROUP="resolver"

ARG HOME="/home/resolver"

# Create user and group
RUN groupadd --system --gid $GID $USER &&\
    useradd --system --create-home --home-dir $HOME --shell /bin/bash --gid $GROUP --uid $UID $USER

# Copy compiled cheqd-did-resolver binary from Stage 1
COPY --from=builder /root /bin

# Copy base config.yaml
COPY config.yaml $HOME

WORKDIR $HOME
USER $USER

EXPOSE 1313

ENTRYPOINT ["cheqd-did-resolver", "serve"]
