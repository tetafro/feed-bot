.PHONY: dep
dep:
	@ go mod tidy && go mod verify

.PHONY: test
test:
	@ go test ./...

.PHONY: cover
cover:
	@ mkdir -p tmp
	@ go test -coverprofile=./tmp/cover.out ./...
	@ go tool cover -html=./tmp/cover.out

.PHONY: lint
lint: go-lint yamllint ansible-lint

.PHONY: go-lint
go-lint:
	@ echo '----------------'
	@ echo 'Running golangci-lint'
	@ echo '----------------'
	@ golangci-lint run --fix && echo OK

.PHONY: yamllint
yamllint:
	@ echo '----------------'
	@ echo 'Running yamllint'
	@ echo '----------------'
	@ yamllint --format colored --strict ./playbook.yml && echo OK

.PHONY: ansible-lint
ansible-lint:
	@ echo '--------------------'
	@ echo 'Running ansible-lint'
	@ echo '--------------------'
	@ ansible-lint -q ./playbook.yml && echo OK

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
	@ ansible-playbook \
	--vault-password-file .vault_pass.txt \
	--private-key ~/.ssh/id_ed25519 \
	--inventory '${SSH_SERVER},' \
	--user ${SSH_USER} \
	./playbook.yml
