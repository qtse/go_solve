GC=gd
LIB_DIR=lib
SRC_DIR=src
BIN=./poly

run :
	gd -L $(LIB_DIR) -o $(BIN) $(SRC_DIR)
	$(BIN)

clean :
	rm -r $(BIN) $(LIB_DIR)/*
