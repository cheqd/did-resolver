#####################################################################
###    STAGE 1: Build cheqd did-resolver binary pre-requisites    ###
#####################################################################

FROM golang:1.17-alpine AS builder

WORKDIR /builder
COPY . ./

# Build did-resolver binary
RUN go mod download && go build -o did-resolver main.go

#####################################################################
###    STAGE 2: Build cheqd did-resolver container image          ###
#####################################################################

FROM alpine:3.16 AS resolver

# Install pre-requisites
RUN apk update && apk add --no-cache bash ca-certificates

# Copy did-resolver binary from Stage 1
COPY --from=builder /builder/did-resolver /usr/local/bin/did-resolver

# Set user directory and details
ARG HOME_DIR="/resolver"
ARG USER="resolver"
SHELL ["/bin/bash", "-euo", "pipefail", "-c"]

# Add non-root user to use in the container
RUN addgroup --system $USER \
    && adduser $USER --system --home $HOME_DIR --shell /bin/bash

# Set working directory & bash defaults
WORKDIR $HOME_DIR
USER $USER

# Copy base config.yaml
COPY --chown=$USER:$USER config.yaml .

EXPOSE 8080

ENTRYPOINT ["did-resolver"]
CMD ["serve"]
