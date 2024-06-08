# BugsChannel

![bugs channel logo](./images/logo.png)

![workflow](https://github.com/williampsena/bugs-channel/actions/workflows/main.yml/badge.svg)

This repository contains information about handling issues with proxy.
I decided to begin this project with the goal of making error handling as simple as possible.
I use [Sentry](https://sentry.io) and [Honeybadger](https://www.honeybadger.io), and both tools are fantastic for quickly tracking down issues. However, the purpose of this project is not to replace them, but rather to provide a simple solution for you to run on premise that is easy and has significant features.

> I started the project with Elixir, but I'm switching to Go to keep things as simple and productive as possible 😅.

# Requirements

```shell
go install golang.org/x/pkgsite/cmd/pkgsite@latest
go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
```

# Challenges
## Done 👌

- Send events to NATs
- Create the BugsChannel logo
- Implement the rate-limit strategy
- In db-less mode, define yaml as an option
- Support Nats and Redis as Pub/Sub
- Create a docker deployment example
- Check for the presence of authentication keys
- In db-less mode, define yaml as an option
- Identify the project by the requested authentication keys

## TODO

- Scrub events to avoid exposing sensitive information
- Get consumers (sub) and producers (pub) on board with NATS
- Generate documentation with pkgsite
- Grpc support
- Support Graylog as a error target
- Support Kibana as a error target

# Running project

The command below starts a web application that listens on port 4000 by default.
By default, Sentry Integration is started and listens on port 4001.

```shell
cp .env.sample .env
make start
make dev
```

The project listens on port 4001 (sentry). At the moment, just Sentry had been set up and you could test the following steps.

- Create a config file named `config.yml` to run as **dbless** mode.

```shell
cp fixtures/settings/config.yml .config/config.yml
```

- Create a file named `main.py`.

```python
import sentry_sdk

sentry_sdk.init(
    "http://key@localhost:4001/1",
    traces_sample_rate=1.0,
)

raise ValueError("Error SDK")
```

- Install python packages

```shell
# using venv
python -m venv .env
. .env/bin/activate
pip install sentry-sdk

# without venv
pip install --user sentry-sdk
```

- Now you can run project

```shell
python main.py
```

# Tests

```shell
make test
```

# Vulnerabilities

```shell
make vulns-check
```

# Docs

```shell
make docs
```
