package main

import (
	"github.com/keithzetterstrom/db_forum/cmd/handlers"
	forumHandler "github.com/keithzetterstrom/db_forum/cmd/handlers/forum"
	profileHandler "github.com/keithzetterstrom/db_forum/cmd/handlers/profile"
	serviceHandler "github.com/keithzetterstrom/db_forum/cmd/handlers/service"
	"github.com/labstack/echo"
)

func main()  {

	e := echo.New()

	forHandler := forumHandler.NewHandler()
	profHandler := profileHandler.NewHandler()
	servHandler := serviceHandler.NewHandler()

	handlers.Router(e, forHandler, profHandler, servHandler)

	e.Logger.Fatal(e.Start(":8080"))
}