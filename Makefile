MAIN = day2
BIN = bin/
CMD = cmd/adventofcode2023/$(MAIN)/
PIPE = # cat internal/orderbook/input2.stream |
ARGS = # 5


all: build run

build:
	go build -o $(BIN)$(MAIN) $(CMD)$(MAIN).go

run:
	$(PIPE) ./$(BIN)$(MAIN) $(ARGS)

clean:
	rm -f $(BIN)$(MAIN) output.log
