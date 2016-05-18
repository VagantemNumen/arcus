package = github.com/VagantemNumen/arcus

.PHONY: release

release:
	mkdir -p release
	GOOS=linux GOARCH=amd64 go build -o release/arcus-linux-amd64 $(package)
	GOOS=freebsd GOARCH=amd64 go build -o release/arcus-freebsd-amd64 $(package)
