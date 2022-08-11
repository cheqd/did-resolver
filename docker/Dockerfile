#####################################################################
###    STAGE 1: Build cheqd did-resolver binary pre-requisites    ###
#####################################################################

FROM golang:1.17-alpine AS builder

WORKDIR /resolver
RUN apk add --no-cache git
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./


# Build did-resolver binary
ENV CGO_ENABLED=0
RUN go build -o cheqd-did-resolver main.go

#####################################################################
###    STAGE 2: Build cheqd did-resolver container image          ###
#####################################################################

FROM gcr.io/distroless/static AS resolver

# Set working directory & bash defaults
WORKDIR /resolver

USER nonroot:nonroot

# Copy compiled cheqd-did-resolver binary from Stage 1
COPY --from=builder /resolver/cheqd-did-resolver /bin/cheqd-did-resolver

# Copy base config.yaml
COPY config.yaml .

EXPOSE 8080

ENTRYPOINT ["cheqd-did-resolver"]
CMD ["serve"]
