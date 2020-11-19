package handlers

import (
	forumHandler "github.com/keithzetterstrom/db_forum/cmd/handlers/forum"
	profileHandler "github.com/keithzetterstrom/db_forum/cmd/handlers/profile"
	serviceHandler "github.com/keithzetterstrom/db_forum/cmd/handlers/service"
	"github.com/labstack/echo"
)


func Router(e *echo.Echo, forum forumHandler.Handler, profile profileHandler.Handler, service serviceHandler.Handler) {
	e.POST("/forum/create", forum.ForumCreate)
	e.POST("/forum/:slug/create", forum.ThreadCreate)
	e.GET("/forum/:slug/details", forum.ForumGet)
	e.GET("/forum/:slug/threads", forum.ForumThreadsGet)
	e.GET("/forum/:slug/users", forum.ForumUsersGet)

	e.GET("/post/:id/details", forum.PostGet)
	e.POST("/post/:id/details", forum.PostUpdate)

	e.POST("/service/clear", service.ServiceClear)
	e.GET("/service/status", service.ServiceStatus)

	e.POST("/thread/:slug_or_id/create", forum.PostCreate)
	e.GET("/thread/:slug_or_id/details", forum.ThreadGet)
	e.POST("/thread/:slug_or_id/details", forum.ThreadUpdate)
	e.GET("/thread/:slug_or_id/posts", forum.ThreadPostsGet)
	e.POST("/thread/:slug_or_id/vote", forum.ThreadVote)

	e.POST("/user/:nickname/create", profile.ProfileCreate)
	e.GET("/user/:nickname/profile", profile.ProfileGet)
	e.POST("/user/:nickname/profile", profile.ProfileUpdate)
}