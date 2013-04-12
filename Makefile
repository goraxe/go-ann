GC = go
TARG = anntest

all: $(TARG)

myann: src/myann/myann.go
	$(GC) install $@

main: src/main/main.go
	$(GC) install $@

$(TARG): myann main
	
run:
	./bin/main

clean:
	$(GC) clean
	rm -rf pkg/* bin/*

