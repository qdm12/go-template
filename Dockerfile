ARG ALPINE_VERSION=3.13
ARG GO_VERSION=1.16
# Sets linux/amd64 in case it's not injected by older Docker versions
ARG BUILDPLATFORM=linux/amd64

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS base
RUN apk --update add git
ENV CGO_ENABLED=0
WORKDIR /tmp/gobuild
# Copy repository code and install Go dependencies
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/

FROM base AS test
# Note on the go race detector:
# - we set CGO_ENABLED=1 to have it enabled
# - we install g++ to support the race detector
ENV CGO_ENABLED=1
RUN apk -q --update --no-cache add g++

FROM base AS lint
ARG GOLANGCI_LINT_VERSION=v1.35.2
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
    sh -s -- -b /usr/local/bin ${GOLANGCI_LINT_VERSION}
COPY .golangci.yml ./
RUN golangci-lint run --timeout=10m

FROM base AS tidy
RUN git init && \
    git config user.email ci@localhost && \
    git config user.name ci && \
    git add -A && git commit -m ci && \
    sed -i '/\/\/ indirect/d' go.mod && \
    go mod tidy && \
    git diff --exit-code -- go.mod

FROM base AS build
COPY --from=qmcgaw/xcputranslate:v0.4.0 /xcputranslate /usr/local/bin/xcputranslate
ARG TARGETPLATFORM
ARG VERSION=unknown
ARG BUILD_DATE="an unknown date"
ARG COMMIT=unknown
RUN GOARCH="$(xcputranslate -targetplatform=${TARGETPLATFORM} -field arch)" \
    GOARM="$(xcputranslate -targetplatform=${TARGETPLATFORM} -field arm)" \
    go build -trimpath -ldflags="-s -w \
    -X 'main.version=$VERSION' \
    -X 'main.buildDate=$BUILD_DATE' \
    -X 'main.commit=$COMMIT' \
    " -o app cmd/app/main.go

FROM --platform=$BUILDPLATFORM alpine:${ALPINE_VERSION} AS alpine
RUN apk --update add ca-certificates tzdata

FROM scratch
USER 1000
COPY --from=alpine --chown=1000 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=alpine --chown=1000 /usr/share/zoneinfo /usr/share/zoneinfo
ENTRYPOINT ["/app"]
EXPOSE 8000/tcp
HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=2 CMD ["/app","healthcheck"]
ENV HTTP_SERVER_ADDRESS=0.0.0.0:8000 \
    HTTP_SERVER_ROOT_URL=/ \
    HTTP_SERVER_LOG_REQUESTS=on \
    HTTP_SERVER_ALLOWED_ORIGINS= \
    HTTP_SERVER_ALLOWED_HEADERS= \
    METRICS_SERVER_ADDRESS=0.0.0.0:9090 \
    LOG_LEVEL=info \
    STORE_TYPE=memory \
    STORE_JSON_FILEPATH=data.json \
    STORE_POSTGRES_ADDRESS=psql:5432 \
    STORE_POSTGRES_USER=postgres \
    STORE_POSTGRES_PASSWORD=postgres \
    STORE_POSTGRES_DATABASE=database \
    HEALTH_SERVER_ADDRESS=127.0.0.1:9999 \
    TZ=America/Montreal
COPY --chown=1000 postgres/schema.sql /schema.sql
ARG VERSION=unknown
ARG BUILD_DATE="an unknown date"
ARG COMMIT=unknown
LABEL \
    org.opencontainers.image.authors="quentin.mcgaw@gmail.com" \
    org.opencontainers.image.version=$VERSION \
    org.opencontainers.image.created=$BUILD_DATE \
    org.opencontainers.image.revision=$COMMIT \
    org.opencontainers.image.url="https://github.com/qdm12/go-template" \
    org.opencontainers.image.documentation="https://github.com/qdm12/go-template/blob/main/README.md" \
    org.opencontainers.image.source="https://github.com/qdm12/go-template" \
    org.opencontainers.image.title="go-template" \
    org.opencontainers.image.description="SHORT_DESCRIPTION"
COPY --from=build --chown=1000 /tmp/gobuild/app /app