NAME			:= 	server.tmp
PROJECT_PATH	:= 	./app/
TEST_PATH		:= 	$(addprefix $(PROJECT_PATH), \
					builtin parser server)

all:
	go build -o $(NAME) $(PROJECT_PATH)*.go

# Run tests
t: $(TEST_PATH)

# $(TEST_PATH):
# 	@echo "[================ RUN TEST $@ ================]"
# 	@cd $@ && go test
# 	@echo

$(TEST_PATH):
	@echo "[================ RUN TEST $@ ================]"
	@cd $@ && go test
	@cd $@ && go test -race -coverprofile=$(subst app/,,$@).txt -covermode=atomic
	@echo


# $(TEST_PATH):
# 	@echo "[================ RUN TEST $@ ================]"
# 	@cd $@ && go test
# 	@cd $@ && go test -cover
# 	@cd $@ && go test -bench=.
# 	@echo 

.PHONY: $(TEST_PATH)