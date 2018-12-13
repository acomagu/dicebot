GENERATED_SOURCES = bot/joiner_mock_test.go \
                    bot/session_mock_test.go \
                    bot/vc_mock_test.go \
                    soundplayer/vc_mock_test.go \
                    soundplayer/voice_channel_joiner_mock_test.go
SOURCES = statik/statik.go $GENERATED_SOURCES `{find -name '*.go'}

dicebot: $SOURCES
	go build -o dicebot

source:V: $SOURCES

dockerimage:V: $SOURCES
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dicebot
	docker build . -t dicebot

test:V: $SOURCES
	GOMAXPROCS=4 go test -v -race ./...

statik/statik.go: public/dicesound.dca bin/statik
	bin/statik -f

public/dicesound.dca: `{ls dicesound.*} bin/dca
	mkdir -p public
	bin/dca -i $prereq[1] > public/dicesound.dca

$GENERATED_SOURCES: bin/moq
	PATH=$PWD/bin:$PATH go generate ./...

bin/statik bin/dca bin/moq:
	command -v dept >/dev/null 2>&1 || go get -v github.com/ktr0731/dept
	GOBIN=$PWD/bin dept build
