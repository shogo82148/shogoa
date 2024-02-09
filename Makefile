.PHONY: test
test:
	go test ./... -v -cover -coverprofile profile.cov
