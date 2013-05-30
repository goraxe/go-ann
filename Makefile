GC = go
TARG = anntest

PWD := $(shell pwd)

all: $(TARG)

warn:
	@ echo "please export GOPATH=$(PWD)"

myann: src/myann/myann.go
	$(GC) install $@

main: src/main/main.go
	$(GC) install $@

$(TARG): warn myann main
	
run:
	./bin/main

clean:
	$(GC) clean
	rm -rf pkg/* bin/*

