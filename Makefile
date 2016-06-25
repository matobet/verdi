check:
	go test -v $(shell glide nv)

all:
	go build

clean:
	-@rm verdi
