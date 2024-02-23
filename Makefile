NAME			:= 	server.tmp
PROJECT_PATH	:= 	./app/

all:
	go build -o $(NAME) $(PROJECT_PATH)*.go


t: unit cov bech

unit:
	go test ./...

cov:
	cd app && go test ./... -race -coverprofile=../coverage.txt -covermode=atomic

bech:
	go test ./... -bench=.

.PHONY: $(TEST_PATH)