ARG ALPINE_VERSION=3.12
ARG GO_VERSION=1.15

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS base
# g++ is installed for the -race detector in go test
RUN apk --update add git g++
ENV CGO_ENABLED=0
WORKDIR /tmp/gobuild
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/

FROM --platform=$BUILDPLATFORM base AS lint
ARG GOLANGCI_LINT_VERSION=v1.34.1
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
    sh -s -- -b /usr/local/bin ${GOLANGCI_LINT_VERSION}
COPY .golangci.yml ./
RUN golangci-lint run --timeout=10m

FROM --platform=$BUILDPLATFORM base AS tidy
COPY .git/ ./.git/
RUN sed -i '/\/\/ indirect/d' go.mod && \
    go mod tidy && \
    git diff --exit-code -- go.mod go.sum

FROM --platform=$BUILDPLATFORM base AS build
COPY --from=qmcgaw/xcputranslate /xcputranslate /usr/local/bin/xcputranslate
ARG TARGETPLATFORM
ARG VERSION=unknown
ARG BUILD_DATE="an unknown date"
ARG COMMIT=unknown
COPY cmd/ ./cmd/
COPY internal/ ./internal/
RUN GOARCH="$(echo ${TARGETPLATFORM} | xcputranslate -field arch)" \
    GOARM="$(echo ${TARGETPLATFORM} | xcputranslate -field arm)" \
    go build -trimpath -ldflags="-s -w \
    -X 'main.version=$VERSION' \
    -X 'main.buildDate=$BUILD_DATE' \
    -X 'main.commit=$COMMIT' \
    " -o app cmd/app/main.go

FROM --platform=$BUILDPLATFORM alpine:${ALPINE_VERSION} AS alpine
RUN apk --update add ca-certificates tzdata

FROM scratch
ARG VERSION=unknown
ARG BUILD_DATE="an unknown date"
ARG COMMIT=unknown
LABEL \
    org.opencontainers.image.authors="quentin.mcgaw@gmail.com" \
    org.opencontainers.image.version=$VERSION \
    org.opencontainers.image.created=$BUILD_DATE \
    org.opencontainers.image.revision=$COMMIT \
    org.opencontainers.image.url="https://github.com/qdm12/REPONAME_GITHUB" \
    org.opencontainers.image.documentation="https://github.com/qdm12/REPONAME_GITHUB/blob/main/README.md" \
    org.opencontainers.image.source="https://github.com/qdm12/REPONAME_GITHUB" \
    org.opencontainers.image.title="REPONAME_GITHUB" \
    org.opencontainers.image.description="SHORT_DESCRIPTION"
COPY --from=alpine --chown=1000 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=alpine --chown=1000 /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=America/Montreal \
    LOG_ENCODING=console \
    LOG_LEVEL=info \
    LISTENING_PORT=8000 \
    ROOT_URL=/ \
    SQL_HOST=postgres \
    SQL_USER=postgres \
    SQL_PASSWORD=postgres \
    SQL_DBNAME=postgres
ENTRYPOINT ["/app"]
HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=2 CMD ["/app","healthcheck"]
USER 1000
COPY --chown=1000 postgres/schema.sql /schema.sql
COPY --from=build --chown=1000 /tmp/gobuild/app /app
