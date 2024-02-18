MAIN = cyoa
BIN = bin/
CMD = cmd/
PIPE = 
ARGS = 


all: build run

build:
	go build -o $(BIN)$(MAIN) $(CMD)$(MAIN)/$(MAIN).go

run:
	$(PIPE) ./$(BIN)$(MAIN) $(ARGS)

clean:
	rm -f $(BIN)$(MAIN) output.log
