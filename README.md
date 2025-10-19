# Feed bot

[![Codecov](https://codecov.io/gh/tetafro/feed-bot/branch/master/graph/badge.svg)](https://codecov.io/gh/tetafro/feed-bot)
[![Go Report](https://goreportcard.com/badge/github.com/tetafro/feed-bot)](https://goreportcard.com/report/github.com/tetafro/feed-bot)
[![CI](https://github.com/tetafro/feed-bot/actions/workflows/push.yml/badge.svg)](https://github.com/tetafro/feed-bot/actions)

Telegram bot that reads RSS feeds and sends them to users.

## Build and run

Create a bot and get Telegram API token from the bot called `@botfather`.

Copy and populate config
```sh
cp config.example.yaml config.yaml
```

Start
```sh
make build run
```

## Deploy

Normally deploy is done by Github actions.

Manual deploy
```sh
SSH_SERVER=10.0.0.1:22 \
SSH_USER=user \
make deploy
```

## Encrypted config

Encrypt
```sh
echo "password" > .vault_pass.txt
make config-encrypt
```
