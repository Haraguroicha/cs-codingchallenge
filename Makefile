GO_BUILD_ENV := CGO_ENABLED=0 GOARCH=amd64
DOCKER_BUILD=$(shell pwd)/.docker_build
DOCKER_CMD=$(DOCKER_BUILD)/cs-codingchallenge

$(DOCKER_CMD): clean
	mkdir -p $(DOCKER_BUILD)
	GOOS=linux $(GO_BUILD_ENV) go build -v -o $(DOCKER_CMD) .

clean:
	rm -rf $(DOCKER_BUILD)

heroku: $(DOCKER_CMD)
	heroku container:push web

local: clean
	mkdir -p $(DOCKER_BUILD)
	$(GO_BUILD_ENV) go build -v -o $(DOCKER_CMD) .

debug: local
	heroku local --procfile Procfile_debug

test:
	go test -v -bench=. -benchmem -cover ./...
