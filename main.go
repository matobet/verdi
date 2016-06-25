package main

import (
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/matobet/verdi/config"
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

	// backend, err := backend.Init()
	// if err != nil {
	// 	log.Fatal("Failed to initialize backend: ", err)
	// }

	// reply, err := backend.Run("RemoveVM", &cmd.IDParams{
	// 	ID: "4392ce9e-5a9b-442e-9037-91764e14b129",
	// })

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(reply)

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339}: ${status} ${response_time} ${method} ${uri} - ${tx_bytes} bytes\n",
	}))
	e.Use(middleware.Recover())
	// e.Static("/", "frontend/static")
	e.Use(middleware.Static("frontend/static"))
	api := e.Group("/api")
	{
		api.GET("/echo", func(c echo.Context) error {
			return c.JSON(200, map[string]string{"hello": "world"})
		})
	}
	e.Run(standard.New(config.Conf.HTTPPort))
}
