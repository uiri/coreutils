%: %.go
	go build -o bin/$@ $<
	./bin/$@ --create-manpage > doc/$@.1
	gzip doc/$@.1

all:
	for SOURCEFILE in *go; do go build -o bin/$${SOURCEFILE:0:-3} $$SOURCEFILE; done

setup:
	if [ ! -d "bin" ]; then mkdir bin; fi
	if [ ! -d "doc" ]; then mkdir doc; fi

doc:	bin
	for BINFILE in bin/*; do ./$$BINFILE --create-manpage > doc/$${BINFILE:4}.1; gzip doc/$${BINFILE:4}.1; done
