GC = 6g
LD = 6l
TARG = ann

DEPS = file.6 network.6
O_FILES =  main.6

all: $(TARG)

.SUFFIXES: .go .6

$(TARG): $(DEPS) $(O_FILES)
	$(LD) -o $@ $(O_FILES)
@echo "Done. Executable is: $@"

#$(O_FILES): %.6: %.go
#	$(GC) -c $<

.go.6:
	$(GC) -c $<


clean:
	rm -rf *.[$(OS)o] *.a [$(OS)].out _obj $(TARG) *.6
