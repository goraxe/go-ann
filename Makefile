GC = go
TARG = ann

DEPS = file.go network.go
O_FILES =  main.go

all: $(TARG)

.SUFFIXES: .go

$(TARG): $(DEPS) $(O_FILES)
	$(GC) build $@ $(O_FILES)
@echo "Done. Executable is: $@"

#$(O_FILES): %.6: %.go
#	$(GC) -c $<

.go:
	$(GC) build $<


clean:
	rm -rf *.[$(OS)o] *.a [$(OS)].out _obj $(TARG)

