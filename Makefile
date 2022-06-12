.PHONY: dep
dep:
	@ go mod tidy && go mod verify && go mod vendor

.PHONY: test
test:
	@ go test ./...

.PHONY: cover
cover:
	@ mkdir -p tmp
	@ go test -coverprofile=./tmp/cover.out ./...
	@ go tool cover -html=./tmp/cover.out

.PHONY: lint
lint:
	@ golangci-lint run --fix

.PHONY: build
build:
	@ go build -o ./bin/feed-bot ./cmd/feed-bot

.PHONY: run
run:
	@ ./bin/feed-bot

.PHONY: docker
docker:
	@ docker build -t ghcr.io/tetafro/feed-bot .

.PHONY: deploy
deploy:
	@ ansible-playbook
	--private-key ~/.ssh/id_ed25519
	--inventory "${SSH_SERVER},'
	--user ${SSH_USER}
	./playbook.yml
