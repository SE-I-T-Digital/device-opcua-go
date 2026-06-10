#
# Copyright (c) 2018, 2019 Intel
# Copyright (c) 2021-2026 Schneider Electric
#
# SPDX-License-Identifier: Apache-2.0
#
FROM golang:1.26.4-alpine3.23@sha256:f23e8b227fb4493eabe03bede4d5a32d04092da71962f1fb79b5f7d1e6c2a17f AS builder
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
FROM alpine:3.24.0@sha256:a2d49ea686c2adfe3c992e47dc3b5e7fa6e6b5055609400dc2acaeb241c829f4

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
