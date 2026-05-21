# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.25

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG VERSION=dev

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
	go build -trimpath -ldflags="-s -w -X github.com/mschuchard/vault-raft-backup/util.Version=${VERSION}" \
	-o /out/vault-raft-backup . \
	&& mkdir -p /tmp-root/tmp \
	&& chmod 1777 /tmp-root/tmp

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build --chown=nonroot:nonroot --chmod=0555 /out/vault-raft-backup /vault-raft-backup
COPY --from=build --chown=nonroot:nonroot --chmod=1777 /tmp-root/tmp /tmp

USER nonroot:nonroot
WORKDIR /
ENTRYPOINT ["/vault-raft-backup"]
