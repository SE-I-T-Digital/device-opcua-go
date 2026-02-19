#
# Copyright (c) 2018, 2019 Intel
# Copyright (c) 2021-2026 Schneider Electric
#
# SPDX-License-Identifier: Apache-2.0
#
FROM golang:1.26-alpine3.22@sha256:169d3991a4f795124a88b33c73549955a3d856e26e8504b5530c30bd245f9f1b AS builder
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
FROM alpine:3.23.3@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659

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
