.PHONY: test
dep:
	@ go mod vendor

.PHONY: test
test:
	@ go test ./...

.PHONY: lint
lint:
	@ golangci-lint run

.PHONY: build
build:
	@ go build -o ./bin/notify-bot .

.PHONY: run
run:
	@ ./bin/notify-bot

.PHONY: docker
docker:
	@ docker build -t tetafro/feed-bot .
