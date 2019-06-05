.PHONY: test
dep:
	@ go mod vendor

.PHONY: lint
lint:
	@ golangci-lint run

.PHONY: build
build:
	@ go build -o ./bin/notify-bot .

.PHONY: run
run:
	@ ./bin/notify-bot
