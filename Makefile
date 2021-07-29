SHELL=/bin/bash

EXE = cdc-gateway

all:
	@echo "building $(EXE) ..."
	@$(MAKE) -s -f make.inc s=static

clean:
	rm -f $(EXE)
