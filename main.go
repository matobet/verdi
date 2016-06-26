package main

import (
	"log"
	"net/http"

	"fmt"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/matobet/verdi/backend"
	"github.com/matobet/verdi/backend/cmd"
	"github.com/matobet/verdi/config"
	"github.com/matobet/verdi/frontend"
)

func main() {
	//go func() {
	//	for {
	//		m, err := mem.VirtualMemory()
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//		fmt.Println(m)
	//		time.Sleep(3 * time.Second)
	//	}
	//}()

	err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration: ", err)
	}

	backend, err := backend.Init()
	if err != nil {
		log.Fatal("Failed to initialize backend: ", err)
	}

	reply, err := backend.Run("RemoveVM", &cmd.IDParams{
		ID: "687692de-e133-4f2b-b715-44cfcc1be81e",
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply)

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339}: ${status} ${latency_human} ${method} ${uri} - ${tx_bytes} bytes\n",
	}))
	e.Use(middleware.Recover())
	api := e.Group("/api")
	{
		api.GET("/echo", func(c echo.Context) error {
			return c.JSON(200, map[string]string{"hello": "world"})
		})
	}
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:     frontend.Asset,
		AssetDir:  frontend.AssetDir,
		AssetInfo: frontend.AssetInfo,
		Prefix:    "/",
	})
	e.GET("/*", standard.WrapHandler(fileServer))
	e.Run(standard.New(config.Conf.HTTPPort))
}
