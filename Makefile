COMMIT = $$(git describe --always)

deps:
	@echo "====> Install dependencies..."

clean:
	@echo "====> Remove installed binary"
	rm -f bin/aclient

install: deps
	@echo "====> Build hget in ./bin "
	go build -ldflags "-X main.GitCommit=\"$(COMMIT)\"" -o bin/aclient
