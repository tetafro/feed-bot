.PHONY: dep
dep:
	@ dep ensure -v

.PHONY: lint
lint:
	@ golangci-lint run

.PHONY: build
build:
	@ go build -o ./bin/notify-bot .

.PHONY: run
run:
	@ ./bin/notify-bot
