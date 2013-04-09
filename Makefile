GC = gccgo
LD = gccgo
TARG = ann

DEPS = network.go
O_FILES = main.go

all: $(TARG)

$(TARG): $(DEPS) $(O_FILES)
	$(GC) -c $(DEPS)
	$(GC) -c $(O_FILES)

clean:
	rm -rf *.out *.a *.o *.8

