.PHONY: test
dep:
	@ go mod vendor

.PHONY: test
test:
	@ go test ./...

.PHONY: cover
cover:
	@ mkdir -p tmp
	@ go test -coverprofile ./tmp/cover.out ./...
	@ go tool cover -html=./tmp/cover.out -o ./tmp/cover.html
	@ rm -f ./tmp/cover.out

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
