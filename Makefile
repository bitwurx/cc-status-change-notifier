.PHONY: build
build:
	@docker run \
		--rm \
		-e CGO_ENABLED=0 \
		-v $(PWD):/usr/src/concord-status-change-notifier \
		-w /usr/src/concord-status-change-notifier \
		golang /bin/sh -c "go get -v -d && go build -a -installsuffix cgo -o main"
	@docker build -t concord/status-change-notifier .
	@rm -f main

.PHONY: test
test:
	@docker run \
		-d \
		-v $(PWD):/go/src/concord-status-change-notifier \
		-v $(PWD)/.src:/go/src \
		-w /go/src/concord-status-change-notifier \
		--name concord-status-change-notifier_test \
		golang /bin/sh -c "go get -v -t -d && go test -v -coverprofile=.coverage.out"
	@docker logs -f concord-status-change-notifier_test
	@docker rm -f concord-status-change-notifier_test

.PHONY: test-short
test-short:
	@docker run \
		--rm \
		-v $(PWD):/go/src/concord-status-change-notifier \
		-v $(PWD)/.src:/go/src \
		-w /go/src/concord-status-change-notifier \
		golang /bin/sh -c "go get -v -t -d && go test -short -v -coverprofile=.coverage.out"
