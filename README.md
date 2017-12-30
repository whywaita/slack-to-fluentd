# slack-to-fluentd

slack log send fluentd

## Setup

### using Docker

- set `SLACK_TOKEN` and `FLUENTD_HOST`
  - edit docker-compose.yml
- `$ docker-compose up -d`

### without Docker

- build
  - `$ cd logger`
  - `$ dep ensure`
  - `$ go build .`

## Author

Tachibana waita a.k.a. [@whywaita](https://github.com/whywaita)
