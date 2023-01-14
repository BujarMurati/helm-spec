.PHONY: test

test:
	go test ./... -cover

snapshot:
	VERSION=$$(git rev-parse HEAD) goreleaser build --snapshot --rm-dist
