
.PHONY: test
test:
	go test ./...

.PHONY: test-c
test-c:
	go test -cover -coverprofile=cover.out ./... && go tool cover -html=cover.out