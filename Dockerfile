#
# Copyright (c) 2018, 2019 Intel
# Copyright (c) 2021-2025 Schneider Electric
#
# SPDX-License-Identifier: Apache-2.0
#
FROM golang:1.24-alpine3.21@sha256:c8a283e2b10c79e3283910fe742953e05e9e4f8c58500c6ed04f1ee3f9ca7732 AS builder
WORKDIR /device-opcua-go

# Install our build time packages.
RUN apk update && apk add --no-cache make git gcc pkgconfig musl-dev

COPY go.* ./

RUN [ ! -d "vendor" ] && go mod download all || echo "skipping..."

ADD cmd cmd
ADD internal internal 
ADD pkg pkg
COPY version.go Makefile ./

ARG ADD_BUILD_TAGS=""
RUN make -e ADD_BUILD_TAGS="$ADD_BUILD_TAGS" build

# Next image - Copy built Go binary into new workspace
FROM alpine:3.21.5@sha256:5405e8f36ce1878720f71217d664aa3dea32e5e5df11acbf07fc78ef5661465b

# dumb-init needed for injected secure bootstrapping entrypoint script when run in secure mode.
# upgrade required to patch for open CVEs.
RUN apk add --update --no-cache dumb-init=1.2.5-r3 && \
  apk upgrade busybox libcrypto3 libssl3 && \
  rm -rf /var/cache/apk/*

# expose command data port
EXPOSE 59997

COPY --from=builder /device-opcua-go/cmd/device-opcua /
COPY --from=builder /device-opcua-go/cmd/res /res
COPY LICENSE Attribution.txt /

ENTRYPOINT ["/device-opcua"]
CMD ["--cp=keeper.http://edgex-core-keeper:59890", "--registry"]
