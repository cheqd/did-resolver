###############################################################
###    STAGE 1: Build did-resolver binary pre-requisites    ###
###############################################################

FROM golang:1.17.8-buster as builder

WORKDIR /root

COPY types ./types
COPY services ./services
COPY go.mod .
COPY go.sum .
COPY main.go .

# Make did-resolver binary
RUN go build -o did-resolver main.go

###############################################################
###           STAGE 2: Build did-resolver runner            ###
###############################################################

FROM ubuntu:focal AS runner
LABEL org.opencontainers.image.description "Cheqd DID-Resolver runner"
LABEL org.opencontainers.image.source "https://github.com/cheqd/cheqd-did-resolver"

# Copy compiled did-resolver binary from Stage 1
COPY --from=builder /root/did-resolver /bin

# Copy base config.yml
WORKDIR /root

EXPOSE 1313
ENTRYPOINT ["did-resolver"]