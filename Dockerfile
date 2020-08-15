ARG ALPINE_VERSION=3.12
ARG GO_VERSION=1.14

FROM alpine:${ALPINE_VERSION} AS alpine
RUN apk --update add ca-certificates tzdata

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder
ARG GOLANGCI_LINT_VERSION=v1.30.0
RUN apk --update add git
ENV CGO_ENABLED=0
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s ${GOLANGCI_LINT_VERSION}
WORKDIR /tmp/gobuild
COPY .golangci.yml .
COPY go.mod go.sum ./
RUN go mod download 2>&1
COPY cmd/app/main.go cmd/app/main.go
COPY internal ./internal
RUN go test ./...
RUN golangci-lint run --timeout=10m
ARG BUILD_DATE
ARG VCS_REF
ARG VERSION
RUN go build -o app -trimpath -ldflags="-s -w \
    -X 'main.BuildDate=$BUILD_DATE' \
    -X 'main.VcsRef=$VCS_REF' \
    -X 'main.Version=$VERSION'" \
    cmd/app/main.go

FROM scratch
ARG VERSION
ARG BUILD_DATE
ARG VCS_REF
LABEL \
    org.opencontainers.image.authors="quentin.mcgaw@gmail.com" \
    org.opencontainers.image.version=$VERSION \
    org.opencontainers.image.created=$BUILD_DATE \
    org.opencontainers.image.revision=$VCS_REF \
    org.opencontainers.image.url="https://github.com/qdm12/REPONAME_GITHUB" \
    org.opencontainers.image.documentation="https://github.com/qdm12/REPONAME_GITHUB/blob/master/README.md" \
    org.opencontainers.image.source="https://github.com/qdm12/REPONAME_GITHUB" \
    org.opencontainers.image.title="REPONAME_GITHUB" \
    org.opencontainers.image.description="SHORT_DESCRIPTION"
COPY --from=alpine --chown=1000 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=alpine --chown=1000 /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=America/Montreal \
    LOG_ENCODING=console \
    LOG_LEVEL=info \
    NODE_ID=-1 \
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
COPY --from=builder --chown=1000 /tmp/gobuild/app /app
