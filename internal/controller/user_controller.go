package controller

import (
	"music-store/internal/handler"
	"music-store/internal/service"

	"github.com/unbxd/go-base/kit/transport/http"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) Bind(tr *http.Transport, opts []http.HandlerOption) {
	tr.POST(
		"/users",
		handler.CreateUserHandler(c.userService),
		handler.NewCreateUserHandlerOption(opts)...,
	)

	tr.GET(
		"/users/:id",
		handler.GetUserHandler(c.userService),
		handler.NewGetUserHandlerOption(opts)...,
	)

	tr.GET(
		"/users",
		handler.GetAllUsersHandler(c.userService),
		handler.NewGetAllUsersHandlerOption(opts)...,
	)

	tr.PUT(
		"/users/:id",
		handler.UpdateUserHandler(c.userService),
		handler.NewUpdateUserHandlerOption(opts)...,
	)

	tr.DELETE(
		"/users/:id",
		handler.DeleteUserHandler(c.userService),
		handler.NewDeleteUserHandlerOption(opts)...,
	)

	tr.POST(
		"/users/:id/like/:song_name",
		handler.LikeSongHandler(c.userService),
		handler.NewLikeSongHandlerOption(opts)...,
	)

	tr.DELETE(
		"/users/:id/unlike/:song_name",
		handler.UnlikeSongHandler(c.userService),
		handler.NewUnlikeSongHandlerOption(opts)...,
	)

	tr.GET(
		"/users/:id/liked_songs",
		handler.GetLikedSongsHandler(c.userService),
		handler.NewGetLikedSongsHandlerOption(opts)...,
	)
}
