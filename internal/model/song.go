package model

type (
	Song struct {
		Name      string    `json:"name"`
		Embedding []float32 `json:"embedding,omitempty"`
	}

	GetSongRequest struct {
		Name string `json:"name"`
	}

	GetSongResponse struct {
		Song *Song `json:"song,omitempty"`
		Err  error `json:"error,omitempty"`
	}

	GetSongListRequest struct {
		Page     int `json:"page,omitempty"`
		PageSize int `json:"page_size,omitempty"`
	}

	GetSongListResponse struct {
		Songs []*Song `json:"songs,omitempty"`
		Err   error   `json:"error,omitempty"`
	}

	CreateSongRequest struct {
		Song Song `json:"song"`
	}

	CreateSongResponse struct {
		Msg string `json:"msg"`
		Err error  `json:"error,omitempty"`
	}

	UpdateSongRequest struct {
		Name string `json:"name"`
		Song Song   `json:"song"`
	}

	UpdateSongResponse struct {
		Msg string `json:"msg"`
		Err error  `json:"error,omitempty"`
	}

	DeleteSongRequest struct {
		Name string `json:"name"`
	}

	DeleteSongResponse struct {
		Msg string `json:"msg"`
		Err error  `json:"error,omitempty"`
	}

	SimilarSongsRequest struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	SimilarSongsResponse struct {
		Songs []string `json:"songs"`
		Err   error    `json:"error,omitempty"`
	}
)
