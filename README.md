# go-template

*SHORT_DESCRIPTION*

<img height="200" src="https://raw.githubusercontent.com/qdm12/go-template/main/title.svg?sanitize=true">

[![Build status](https://github.com/qdm12/go-template/workflows/CI/badge.svg)](https://github.com/qdm12/go-template/actions?query=workflow%3ACI)
[![Docker Pulls](https://img.shields.io/docker/pulls/qmcgaw/go-template-docker.svg)](https://hub.docker.com/r/qmcgaw/go-template-docker)
[![Docker Stars](https://img.shields.io/docker/stars/qmcgaw/go-template-docker.svg)](https://hub.docker.com/r/qmcgaw/go-template-docker)
[![Image size](https://images.microbadger.com/badges/image/qmcgaw/go-template-docker.svg)](https://microbadger.com/images/qmcgaw/go-template-docker)
[![Image version](https://images.microbadger.com/badges/version/qmcgaw/go-template-docker.svg)](https://microbadger.com/images/qmcgaw/go-template-docker)

[![Join Slack channel](https://img.shields.io/badge/slack-@qdm12-yellow.svg?logo=slack)](https://join.slack.com/t/qdm12/shared_invite/enQtOTE0NjcxNTM1ODc5LTYyZmVlOTM3MGI4ZWU0YmJkMjUxNmQ4ODQ2OTAwYzMxMTlhY2Q1MWQyOWUyNjc2ODliNjFjMDUxNWNmNzk5MDk)
[![GitHub last commit](https://img.shields.io/github/last-commit/qdm12/go-template.svg)](https://github.com/qdm12/go-template/commits/main)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/qdm12/go-template.svg)](https://github.com/qdm12/go-template/graphs/contributors)
[![GitHub issues](https://img.shields.io/github/issues/qdm12/go-template.svg)](https://github.com/qdm12/go-template/issues)

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

1. You can update the image with `docker pull qmcgaw/go-template-docker:latest` or use one of [tags available](https://hub.docker.com/r/qmcgaw/go-template-docker/tags)

### Environment variables

| Environment variable | Default | Possible values | Description |
| --- | --- | --- | --- |
| `HTTP_SERVER_ADDRESS` | `0.0.0.0:8000` | Valid address | HTTP server listening address |
| `HTTP_SERVER_ROOT_URL` | `/` | URL path | HTTP server root URL |
| `HTTP_SERVER_LOG_REQUESTS` | `on` | `on` or `off` | Log requests and responses information |
| `HTTP_SERVER_ALLOWED_ORIGINS` | | CSV of addresses | Comma separated list of addresses to allow for CORS |
| `HTTP_SERVER_ALLOWED_HEADERS` | | CSV of HTTP header keys | Comma separated list of header keys to allow for CORS |
| `METRICS_SERVER_ADDRESS` | `0.0.0.0:9090` | Valid address | Prometheus HTTP server listening address |
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

    1. Install [Docker](https://docs.docker.com/install/)
       - On Windows, share a drive with Docker Desktop and have the project on that partition
       - On OSX, share your project directory with Docker Desktop
    1. With [Visual Studio Code](https://code.visualstudio.com/download), install the [remote containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
    1. In Visual Studio Code, press on `F1` and select `Remote-Containers: Open Folder in Container...`
    1. Your dev environment is ready to go!... and it's running in a container :+1: So you can discard it and update it easily!

    </p></details>

    <details><summary>Locally</summary><p>

    1. Install [Go](https://golang.org/dl/), [Docker](https://www.docker.com/products/docker-desktop) and [Git](https://git-scm.com/downloads)
    1. Install Go dependencies with

        ```sh
        go mod download
        ```

    1. Install [golangci-lint](https://github.com/golangci/golangci-lint#install)
    1. You might want to use an editor such as [Visual Studio Code](https://code.visualstudio.com/download) with the [Go extension](https://code.visualstudio.com/docs/languages/go). Working settings are already in [.vscode/settings.json](https://github.com/qdm12/go-template/main/.vscode/settings.json).

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
