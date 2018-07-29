CC := go build
BINARIES := hue-cli

all: $(BINARIES)

%: %.go
	$(CC) -o $@ $^

clean:
	$(RM) $(BINARIES)
