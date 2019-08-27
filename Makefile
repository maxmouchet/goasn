.PHONY: all
all: goasn

.PHONY: clean
clean:
	rm -f goasn

goasn: $(shell find . -name '*.go')
	go build -o goasn main.go