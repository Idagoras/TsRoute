BIN_NAME  = tsroute
GO_FILES_PATH = ./cmd
build:
	go build -o ./cmd/bin/$(BIN_NAME) $(GO_FILES_PATH)

.PHONY:build