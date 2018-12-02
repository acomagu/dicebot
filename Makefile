GENERATED_SOURCES = bot/joiner_mock_test.go \
                    bot/session_mock_test.go \
                    bot/vc_mock_test.go \
                    soundplayer/vc_mock_test.go \
                    soundplayer/voice_channel_joiner_mock_test.go
SOURCES = statik/statik.go $(GENERATED_SOURCES) $(wildcard *.go) $(wildcard **/*.go)

dicebot: $(SOURCES)
	go build -o dicebot

.PHONY: dockerimage
dockerimage: $(SOURCES)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dicebot
	docker build . -t dicebot

.PHONY: test
test: $(SOURCES)
	GOMAXPROCS=4 go test -v -race ./...

statik/statik.go: public/dicesound.dca bin/statik
	bin/statik -f

public/dicesound.dca: dicesound.* bin/dca
	mkdir -p public
	bin/dca -i $< > public/dicesound.dca

$(GENERATED_SOURCES): bin/moq
	PATH=$(PWD)/bin:$(PATH) go generate ./...

bin/statik bin/dca bin/moq:
	command -v dept >/dev/null 2>&1 || go get -v github.com/ktr0731/dept
	GOBIN=$(CURDIR)/bin dept build
