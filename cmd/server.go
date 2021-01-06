package main

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/keithzetterstrom/db_forum/cmd/handlers"
	forumHandler "github.com/keithzetterstrom/db_forum/cmd/handlers/forum"
	profileHandler "github.com/keithzetterstrom/db_forum/cmd/handlers/profile"
	serviceHandler "github.com/keithzetterstrom/db_forum/cmd/handlers/service"
	"github.com/keithzetterstrom/db_forum/internal/storages"
	"github.com/labstack/echo"
)

func main()  {
	e := echo.New()

	connectionString := "postgres://forum_user:1221@localhost/tp_forum?sslmode=disable"
	config, err := pgx.ParseURI(connectionString)
	if err != nil {
		fmt.Println(err)
		return
	}

	db, err := pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig:     config,
			MaxConnections: 2000,
		})

	if err != nil {
		fmt.Println(err)
		return
	}

	forStorage := storages.NewForumStorage(db)
	thrStorage := storages.NewThreadStorage(db)
	poStorage := storages.NewPostStorage(db)
	uStorage := storages.NewUserStorage(db)

	forHandler := forumHandler.NewHandler(thrStorage, poStorage, forStorage, uStorage)
	profHandler := profileHandler.NewHandler(uStorage)
	servHandler := serviceHandler.NewHandler(forStorage)

	handlers.Router(e, forHandler, profHandler, servHandler)

	e.Logger.Fatal(e.Start(":5000"))
}