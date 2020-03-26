# REPONAME_GITHUB

*SHORT_DESCRIPTION*

<img height="200" src="title.svg?sanitize=true">

[![Build status](https://github.com/qdm12/REPONAME_GITHUB/workflows/Buildx%20latest/badge.svg)](https://github.com/qdm12/REPONAME_GITHUB/actions?query=workflow%3A%22Buildx+latest%22)
[![Docker Pulls](https://img.shields.io/docker/pulls/qmcgaw/REPONAME_DOCKER.svg)](https://hub.docker.com/r/qmcgaw/REPONAME_DOCKER)
[![Docker Stars](https://img.shields.io/docker/stars/qmcgaw/REPONAME_DOCKER.svg)](https://hub.docker.com/r/qmcgaw/REPONAME_DOCKER)
[![Image size](https://images.microbadger.com/badges/image/qmcgaw/REPONAME_DOCKER.svg)](https://microbadger.com/images/qmcgaw/REPONAME_DOCKER)
[![Image version](https://images.microbadger.com/badges/version/qmcgaw/REPONAME_DOCKER.svg)](https://microbadger.com/images/qmcgaw/REPONAME_DOCKER)

[![Join Slack channel](https://img.shields.io/badge/slack-@qdm12-yellow.svg?logo=slack)](https://join.slack.com/t/qdm12/shared_invite/enQtOTE0NjcxNTM1ODc5LTYyZmVlOTM3MGI4ZWU0YmJkMjUxNmQ4ODQ2OTAwYzMxMTlhY2Q1MWQyOWUyNjc2ODliNjFjMDUxNWNmNzk5MDk)
[![GitHub last commit](https://img.shields.io/github/last-commit/qdm12/REPONAME_GITHUB.svg)](https://github.com/qdm12/REPONAME_GITHUB/issues)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/qdm12/REPONAME_GITHUB.svg)](https://github.com/qdm12/REPONAME_GITHUB/issues)
[![GitHub issues](https://img.shields.io/github/issues/qdm12/REPONAME_GITHUB.svg)](https://github.com/qdm12/REPONAME_GITHUB/issues)

## Features

- Compatible with `amd64`, `386`, `arm64`, `arm32v7`, `arm32v6`, `ppc64le` and `s390x` CPU architectures
- [Docker image tags and sizes](https://hub.docker.com/r/qmcgaw/REPONAME_DOCKER/tags)

## Setup

1. Use the following command:

    ```sh
    docker run -d qmcgaw/REPONAME_DOCKER
    ```

    You can also use [docker-compose.yml](https://github.com/qdm12/REPONAME_GITHUB/blob/master/docker-compose.yml) with:

    ```sh
    docker-compose up -d
    ```

1. You can update the image with `docker pull qmcgaw/REPONAME_DOCKER:latest` or use one of [tags available](https://hub.docker.com/r/qmcgaw/REPONAME_DOCKER/tags)

### Environment variables

| Environment variable | Default | Possible values | Description |
| --- | --- | --- | --- |
| `HTTP_TIMEOUT` | `3000` | Integer from 1 | Default HTTP client timeout in milliseconds |
| `LOG_ENCODING` | `json` | `json`, `console` | Logging format |
| `LOG_LEVEL` | `info` | `debug`, `info`, `warning`, `error` | Logging level |
| `NODE_ID` | `0` | Integer | Node ID for clusters |
| `LISTENING_PORT` | `8000` | Integer between `1` and `65535` | Internal listening TCP port |
| `ROOT_URL` | `/` | URL path *string* | URL path, used if behind a reverse proxy |
| `SQL_HOST` | `postgres` | *string* | Database hostname |
| `SQL_USER` | `postgres` | *string* | Database user |
| `SQL_PASSWORD` | `postgres` | *string* | Database password |
| `SQL_DBNAME` | `postgres` | *string* | Database name |
| `REDIS_HOST` | `redis` | *string* | Redis hostname |
| `REDIS_PORT` | `6379` | Integer between `1` and `65535` | Redis listening TCP port |
| `REDIS_PASSWORD` | | *string* | Redis password if needed |
| `GOTIFY_URL` | | URL *string* | URL to Gotify server |
| `GOTIFY_TOKEN` | | *string* | Token for Gotify server |
| `TZ` | `America/Montreal` | *string* | Timezone |

## Development

1. Setup your environment

    <details><summary>Using VSCode and Docker</summary><p>

    1. Install [Docker](https://docs.docker.com/install/)
       - On Windows, share a drive with Docker Desktop and have the project on that partition
       - On OSX, share your project directory with Docker Desktop
    1. With [Visual Studio Code](https://code.visualstudio.com/download), install the [remote containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
    1. In Visual Studio Code, press on `F1` and select `Remote-Containers: Open Folder in Container...`
    1. Your dev environment is ready to go!... and it's running in a container :+1:

    </p></details>

    <details><summary>Locally</summary><p>

    Install [Go](https://golang.org/dl/), [Docker](https://www.docker.com/products/docker-desktop) and [Git](https://git-scm.com/downloads); then:

    ```sh
    go mod download
    go get github.com/golang/mock/gomock
    go get github.com/golang/mock/mockgen
    ```

    And finally install [golangci-lint](https://github.com/golangci/golangci-lint#install)

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
    docker build -t qmcgaw/REPONAME_DOCKER .
    ```

1. See [Contributing](.github/CONTRIBUTING.md) for more information on how to contribute to this repository.

## TODOs

- Switch to Go's mock library
- Use [Gox](https://github.com/mitchellh/gox)
- Database
    - Add database type environment variable
    - Add bbolt as database option
    - Add Redis interface in a kv package (as data)
    - Build tag for database integration tests
    - Unit tests with mocked
    - Integration tests
- Unit tests the server with httprecorder

## License

This repository is under an [MIT license](https://github.com/qdm12/REPONAME_GITHUB/master/license) unless otherwise indicated
