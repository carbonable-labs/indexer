# syntax=docker/dockerfile:1

ARG GO_VERSION=1.22.2-alpine3.19
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS builder
WORKDIR /src

RUN --mount=type=cache,target=/var/cache/apk apk update && apk upgrade && apk add --no-cache build-base

RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,source=go.sum,target=go.sum \
  --mount=type=bind,source=go.mod,target=go.mod \
  go mod download -x

ARG TARGETARCH

# RUN --mount=type=cache,target=/go/pkg/mod/ \
#   --mount=type=bind,target=. \
#   GOARCH=$TARGETARCH go build -ldflags="-linkmode external -extldflags -static" -o indexer main.go
COPY . .
RUN set -eux; \
  GOARCH=$TARGETARCH go build -ldflags="-linkmode external -extldflags -static" -o indexer main.go


FROM alpine:3.19 AS production

RUN --mount=type=cache,target=/var/cache/apk \
  apk --update add \
  ca-certificates \
  tzdata \
  && \
  update-ca-certificates


WORKDIR /srv/app
COPY --from=builder /src/indexer ./indexer

ARG UID=10001
RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid "${UID}" \
  appuser
USER appuser

EXPOSE 8080

ENTRYPOINT [ "/srv/app/indexer" ]
