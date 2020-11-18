package main

import (
	"github.com/keithzetterstrom/db_forum/cmd/handlers"
	"github.com/labstack/echo"
)

func main()  {
	e := echo.New()

	handlers.Router(e)

	e.Logger.Fatal(e.Start(":8080"))
}