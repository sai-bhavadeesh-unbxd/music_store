package model

type (
	Song struct {
		Name      string    `json:"name"`
		Embedding []float64 `json:"embedding"`
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
		Song Song `json:"song"` // Make consistent with UpdateSongRequest
	}

	CreateSongResponse struct {
		Msg string `json:"msg"`
		Err error  `json:"error,omitempty"`
	}

	UpdateSongRequest struct {
		Name string `json:"name"` // Add Name field like User model has ID
		Song Song   `json:"song"` // Make it non-pointer to avoid nil issues
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
)
