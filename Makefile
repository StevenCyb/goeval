test:
	@go test ./... -cover

test_local:
	@go test ./... -coverprofile="/tmp/go-cover.tmp" $@
	@go tool cover -html="/tmp/go-cover.tmp"
	@unlink "/tmp/go-cover.tmp"