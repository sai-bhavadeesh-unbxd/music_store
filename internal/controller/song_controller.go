package controller

import (
	"music-store/internal/handler"
	"music-store/internal/service"

	"github.com/unbxd/go-base/kit/transport/http"
)

type SongController struct {
	songService service.SongService
}

func NewSongController(songService service.SongService) *SongController {
	return &SongController{songService: songService}
}

func (c *SongController) Bind(tr *http.Transport, opts []http.HandlerOption) {
	tr.POST(
		"/songs",
		handler.CreateSongHandler(c.songService),
		handler.NewCreateSongHandlerOption(opts)...,
	)

	tr.GET(
		"/songs/:name",
		handler.GetSongHandler(c.songService),
		handler.NewGetSongHandlerOption(opts)...,
	)
	tr.GET(
		"/songs",
		handler.GetAllSongsHandler(c.songService),
		handler.NewGetAllSongsHandlerOption(opts)...,
	)

	tr.PUT(
		"/songs/:name",
		handler.UpdateSongHandler(c.songService),
		handler.NewUpdateSongHandlerOption(opts)...,
	)

	tr.DELETE(
		"/songs/:name",
		handler.DeleteSongHandler(c.songService),
		handler.NewDeleteSongHandlerOption(opts)...,
	)
}
