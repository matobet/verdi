package main

import (
	"log"
	"net/http"

	"strings"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/matobet/verdi/backend"
	"github.com/matobet/verdi/config"
	"github.com/matobet/verdi/frontend"
	"github.com/matobet/verdi/ws"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration: ", err)
	}

	backend, err := backend.Init()
	if err != nil {
		log.Fatal("Failed to initialize backend: ", err)
	}

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339}: ${status} ${latency_human} ${method} ${uri} - ${tx_bytes} bytes\n",
	}))
	e.Use(middleware.Recover())

	websocketHandler := ws.NewHandler(backend)

	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:     frontend.Asset,
		AssetDir:  frontend.AssetDir,
		AssetInfo: frontend.AssetInfo,
		Prefix:    "/",
	})

	e.Use(standard.WrapMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.RequestURI, "/ws") {
				websocketHandler.ServeHTTP(w, r)
			} else if r.RequestURI == "/" {
				fileServer.ServeHTTP(w, r)
			} else if _, err := frontend.Asset(r.RequestURI[1:]); err == nil {
				fileServer.ServeHTTP(w, r)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}))
	e.Run(standard.New(config.Conf.HTTPPort))
}
