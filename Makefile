BIN           = $(GOPATH)/bin
ON            = $(BIN)/on
GO_BINDATA    = $(BIN)/go-bindata
DIST          = frontend/dist
BINDATA       = frontend/bindata.go
SCRIPTS_DIR   = backend/redis/scripts
REDIS_SCRIPTS = backend/redis/bindata.go

all: $(BINDATA) $(REDIS_SCRIPTS)
	@go build

$(ON):
	go install github.com/olebedev/on

$(GO_BINDATA):
	go install github.com/jteeuwen/go-bindata/...

$(BINDATA):
	@npm run build
	$(GO_BINDATA) -o=$@ -pkg frontend -prefix $(DIST) $(DIST)/...

$(REDIS_SCRIPTS): $(wildcard $(SCRIPTS_DIR)/*.lua)
	$(GO_BINDATA) -o=$@ -pkg redis -prefix $(SCRIPTS_DIR) $(SCRIPTS_DIR)/...

check:
	go test -v $(shell glide nv)

clean:
	-@rm verdi
	-@rm -rf $(DIST)/*
	-@rm $(BINDATA)
	-@rm $(REDIS_SCRIPTS)
