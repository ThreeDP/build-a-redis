NAME			:= 	server.tmp
PROJECT_PATH	:= 	./app/
DEFAUL_PORT		:= 	6379
HOST_MASTER		:= 	0.0.0.0
PORT			:= 	6377

all:
	go build -o $(NAME) $(PROJECT_PATH)*.go

run_m:
	./$(NAME)

run_s:
	./$(NAME) --port $(PORT) --replicaof $(HOST_MASTER) $(DEFAUL_PORT)

t: unit cov bench

unit:
	go test ./...

cov:
	cd app && go test ./... -race -coverprofile=../coverage.txt -covermode=atomic

bench:
	go test ./... -bench=.

.PHONY: $(TEST_PATH)