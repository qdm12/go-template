# go-template

*SHORT_DESCRIPTION*

![Title](https://raw.githubusercontent.com/qdm12/go-template/main/title.svg)

[![Build status](https://github.com/qdm12/go-template/workflows/CI/badge.svg)](https://github.com/qdm12/go-template/actions?query=workflow%3ACI)

[![dockeri.co](https://dockeri.co/image/qmcgaw/go-template-docker)](https://hub.docker.com/r/qmcgaw/go-template-docker)

![Last release](https://img.shields.io/github/release/qdm12/go-template?label=Last%20release)
![Last Docker tag](https://img.shields.io/docker/v/qmcgaw/go-template-docker?sort=semver&label=Last%20Docker%20tag)
[![Last release size](https://img.shields.io/docker/image-size/qmcgaw/go-template-docker?sort=semver&label=Last%20released%20image)](https://hub.docker.com/r/qmcgaw/go-template-docker/tags?page=1&ordering=last_updated)
![GitHub last release date](https://img.shields.io/github/release-date/qdm12/go-template?label=Last%20release%20date)
![Commits since release](https://img.shields.io/github/commits-since/qdm12/go-template/latest?sort=semver)

[![Latest size](https://img.shields.io/docker/image-size/qmcgaw/go-template-docker/latest?label=Latest%20image)](https://hub.docker.com/r/qmcgaw/go-template-docker/tags)

[![GitHub last commit](https://img.shields.io/github/last-commit/qdm12/go-template.svg)](https://github.com/qdm12/go-template/commits/main)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/qdm12/go-template.svg)](https://github.com/qdm12/go-template/graphs/contributors)
[![GitHub closed PRs](https://img.shields.io/github/issues-pr-closed/qdm12/go-template.svg)](https://github.com/qdm12/go-template/pulls?q=is%3Apr+is%3Aclosed)
[![GitHub issues](https://img.shields.io/github/issues/qdm12/go-template.svg)](https://github.com/qdm12/go-template/issues)
[![GitHub closed issues](https://img.shields.io/github/issues-closed/qdm12/go-template.svg)](https://github.com/qdm12/go-template/issues?q=is%3Aissue+is%3Aclosed)

[![Lines of code](https://img.shields.io/tokei/lines/github/qdm12/go-template)](https://github.com/qdm12/go-template)
![Code size](https://img.shields.io/github/languages/code-size/qdm12/go-template)
![GitHub repo size](https://img.shields.io/github/repo-size/qdm12/go-template)
![Go version](https://img.shields.io/github/go-mod/go-version/qdm12/go-template)

![Visitors count](https://visitor-badge.laobi.icu/badge?page_id=go-template.readme)

## Features

- Compatible with `amd64`, `386`, `arm64`, `arm32v7`, `arm32v6`, `ppc64le`, `s390x` and `riscv64` CPU architectures
- [Docker image tags and sizes](https://hub.docker.com/r/qmcgaw/go-template-docker/tags)

## Setup

1. Use the following command:

    ```sh
    docker run -d qmcgaw/go-template-docker
    ```

    You can also use [docker-compose.yml](https://github.com/qdm12/go-template/blob/main/docker-compose.yml) with:

    ```sh
    docker-compose up -d
    ```

1. You can update the image with `docker pull qmcgaw/go-template-docker:latest` or use one of the [tags available](https://hub.docker.com/r/qmcgaw/go-template-docker/tags)

### Environment variables

| Environment variable | Default | Possible values | Description |
| --- | --- | --- | --- |
| `HTTP_SERVER_ADDRESS` | `:8000` | Valid address | HTTP server listening address |
| `HTTP_SERVER_ROOT_URL` | `/` | URL path | HTTP server root URL |
| `HTTP_SERVER_LOG_REQUESTS` | `on` | `on` or `off` | Log requests and responses information |
| `HTTP_SERVER_ALLOWED_ORIGINS` | | CSV of addresses | Comma separated list of addresses to allow for CORS |
| `HTTP_SERVER_ALLOWED_HEADERS` | | CSV of HTTP header keys | Comma separated list of header keys to allow for CORS |
| `METRICS_SERVER_ADDRESS` | `:9090` | Valid address | Prometheus HTTP server listening address |
| `LOG_LEVEL` | `info` | `debug`, `info`, `warning`, `error` | Logging level |
| `STORE_TYPE` | `memory` | `memory`, `json` or `postgres` | Data store type |
| `STORE_JSON_FILEPATH` | `data.json` | Valid filepath | JSON file to use if `STORE_TYPE=json` |
| `STORE_POSTGRES_ADDRESS` | `psql:5432` | Valid address | Postgres database address if `STORE_TYPE=postgres` |
| `STORE_POSTGRES_USER` | `postgres` | | Postgres database user if `STORE_TYPE=postgres` |
| `STORE_POSTGRES_PASSWORD` | `postgres` | | Postgres database password if `STORE_TYPE=postgres` |
| `STORE_POSTGRES_DATABASE` | `database` | | Postgres database name if `STORE_TYPE=postgres` |
| `HEALTH_SERVER_ADDRESS` | `127.0.0.1:9999` | Valid address | Health server listening address |
| `TZ` | `America/Montreal` | *string* | Timezone |

## Development

1. Setup your environment

    <details><summary>Using VSCode and Docker (easier)</summary><p>

    Please refer to the corresponding [readme](.devcontainer).

    </p></details>

    <details><summary>Locally</summary><p>

    1. Install [Go](https://golang.org/dl/), [Docker](https://www.docker.com/products/docker-desktop) and [Git](https://git-scm.com/downloads)
    1. Install Go dependencies with

        ```sh
        go mod download
        ```

    1. Install [golangci-lint](https://github.com/golangci/golangci-lint#install)
    1. You might want to use an editor such as [Visual Studio Code](https://code.visualstudio.com/download) with the [Go extension](https://code.visualstudio.com/docs/languages/go).

    </p></details>

1. Commands available:

    ```sh
    # Build the binary
    go build cmd/app/main.go
    # Test the code
    go test ./...
    # Lint the code
    golangci-lint run
    # Build the Docker image
    docker build -t qmcgaw/go-template-docker .
    ```

1. See [Contributing](https://github.com/qdm12/go-template/main/.github/CONTRIBUTING.md) for more information on how to contribute to this repository.

## TODOs

## License

This repository is under an [MIT license](https://github.com/qdm12/go-template/main/license) unless otherwise indicated
