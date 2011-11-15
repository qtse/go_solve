GC=gd
LIB_DIR=lib
SRC_DIR=src
BIN=./poly

run :
	gd -L $(LIB_DIR) -o $(BIN) $(SRC_DIR)
	$(BIN) 10 15

clean :
	rm -r $(BIN) $(LIB_DIR)/*
