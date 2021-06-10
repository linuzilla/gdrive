SRCS = $(wildcard *.go)
OBJS = $(patsubst %.go,%.so,$(SRCS))

all: $(OBJS)

clean:
	@rm -f $(OBJS)

%.so: %.go
	go build -buildmode=plugin -o $@ $<
