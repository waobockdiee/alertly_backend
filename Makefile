build-AlertlyAPIFunction:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o $(ARTIFACTS_DIR)/bootstrap ./cmd/app

build-AlertlyCronjobFunction:
	GOOS=linux GOARCH=amd64 go build -o $(ARTIFACTS_DIR)/bootstrap ./cmd/cronjob
