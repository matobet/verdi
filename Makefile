BIN           = $(GOPATH)/bin
ON            = $(BIN)/on
GO_BINDATA    = $(BIN)/go-bindata
DIST          = frontend/dist
BINDATA       = frontend/bindata.go
BINDATA_FLAGS = -pkg frontend -prefix $(DIST)

all: $(BINDATA)
	@go build

$(ON):
	go install github.com/olebedev/on

$(GO_BINDATA):
	go install github.com/jteeuwen/go-bindata/...

$(BINDATA):
	@npm run build
	$(GO_BINDATA) -o=$@ $(BINDATA_FLAGS) $(DIST)/...

check:
	go test -v $(shell glide nv)

clean:
	-@rm verdi
	-@rm -rf $(DIST)/*
	-@rm $(BINDATA)
