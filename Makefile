COMMIT = $$(git describe --always)

deps:
	@echo "====> Install dependencies..."

clean:
	@echo "====> Remove installed binary"
	rm -f bin/aclient

multiget:
	@echo "====> Building multiget"
	go build -o bin/cmd/multiget github.com/tortuoise/aclient/cmd/multiget

install: deps
	@echo "====> Build aclient in ./bin "
	go build -ldflags "-X main.GitCommit=\"$(COMMIT)\"" -o bin/aclient
