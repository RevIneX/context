APP_NAME    := context
CMD_DIR     := ./cmd/$(APP_NAME)
LDFLAGS     := -ldflags="-s -w"

.PHONY: build build-win run clean

build:
	go build $(LDFLAGS) -o $(APP_NAME) $(CMD_DIR)

build-win:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(APP_NAME).exe $(CMD_DIR)

run:
	go run $(CMD_DIR)

clean:
	go clean
	rm -f $(APP_NAME) $(APP_NAME).exe