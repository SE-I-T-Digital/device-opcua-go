#
# Copyright (c) 2018, 2019 Intel
# Copyright (c) 2021-2026 Schneider Electric
#
# SPDX-License-Identifier: Apache-2.0
#
FROM golang:1.26-alpine3.23@sha256:c2a1f7b2095d046ae14b286b18413a05bb82c9bca9b25fe7ff5efef0f0826166 AS builder
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
FROM alpine:3.23.4@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11

# dumb-init needed for injected secure bootstrapping entrypoint script when run in secure mode.
# upgrade required to patch for open CVEs.
RUN apk add --update --no-cache dumb-init=1.2.5-r3 && \
  apk upgrade busybox libcrypto3 libssl3 zlib && \
  rm -rf /var/cache/apk/*

# expose command data port
EXPOSE 59997

COPY --from=builder /device-opcua-go/cmd/device-opcua /
COPY --from=builder /device-opcua-go/cmd/res /res
COPY LICENSE Attribution.txt /

ENTRYPOINT ["/device-opcua"]
CMD ["--cp=keeper.http://edgex-core-keeper:59890", "--registry"]
