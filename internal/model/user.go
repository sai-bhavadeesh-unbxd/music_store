package model

type (
	User struct {
		ID         string    `json:"id"`
		Name       string    `json:"name"`
		LikedSongs []string  `json:"liked_songs,omitempty"`
		Embedding  []float64 `json:"embedding,omitempty"`
	}

	GetUserRequest struct {
		ID string `json:"id"`
	}

	GetUserResponse struct {
		User *User `json:"user,omitempty"`
		Err  error `json:"error,omitempty"`
	}

	GetUserListRequest struct {
		Page     int `json:"page,omitempty"`
		PageSize int `json:"page_size,omitempty"`
	}

	GetUserListResponse struct {
		Users []*User `json:"users,omitempty"`
		Err   error   `json:"error,omitempty"`
	}

	CreateUserRequest struct {
		User User `json:"user"`
	}

	CreateUserResponse struct {
		Msg string `json:"msg"`
		Err error  `json:"error,omitempty"`
	}

	UpdateUserRequest struct {
		ID   string `json:"id"`
		User User   `json:"user"`
	}

	UpdateUserResponse struct {
		Msg string `json:"msg"`
		Err error  `json:"error,omitempty"`
	}

	DeleteUserRequest struct {
		ID string `json:"id"`
	}

	DeleteUserResponse struct {
		Msg string `json:"msg"`
		Err error  `json:"error,omitempty"`
	}

	LikeSongRequest struct {
		UserID   string `json:"user_id"`
		SongName string `json:"song_name"`
	}

	LikeSongResponse struct {
		Msg string `json:"msg"`
		Err error  `json:"error,omitempty"`
	}

	UnlikeSongRequest struct {
		UserID   string `json:"user_id"`
		SongName string `json:"song_name"`
	}

	UnlikeSongResponse struct {
		Msg string `json:"msg"`
		Err error  `json:"error,omitempty"`
	}

	GetLikedSongsRequest struct {
		UserID string `json:"user_id"`
	}

	GetLikedSongsResponse struct {
		LikedSongs []string `json:"liked_songs,omitempty"`
		Err        error    `json:"error,omitempty"`
	}
)
