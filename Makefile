PROJECT_PATH	:= 	./app/
TEST_PATH		:= 	$(addprefix $(PROJECT_PATH), \
					builtin .)

all:

# Run tests
t: $(TEST_PATH)

$(TEST_PATH):
	cd $@ && go test

.PHONY: $(TEST_PATH)