package handler

import (
	"context"
	"encoding/json"
	"music-store/internal/model"
	"music-store/internal/service"
	net_http "net/http"

	"github.com/pkg/errors"
	"github.com/unbxd/go-base/kit/endpoint"
	"github.com/unbxd/go-base/kit/transport/http"
)

var (
	errBadRequest = errors.New("bad request")
)

func MakeCreateSongEndpoint(s service.SongService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(model.CreateSongRequest)
		if !ok {
			return nil, errors.Wrap(
				errBadRequest, "failed to cast object to CreateSongRequest",
			)
		}
		msg, err := s.CreateSong(ctx, &req)
		return model.CreateSongResponse{Msg: msg, Err: err}, nil
	}
}

func MakeGetSongEndpoint(s service.SongService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(model.GetSongRequest)
		if !ok {
			return nil, errors.Wrap(
				errBadRequest, "failed to cast object to GetSongRequest",
			)
		}
		song, err := s.GetSong(ctx, req.Name)
		return model.GetSongResponse{Song: song.Song, Err: err}, nil
	}
}

func MakeGetAllSongsEndpoint(s service.SongService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		songs, err := s.GetAllSongs(ctx)
		return model.GetSongListResponse{Songs: songs.Songs, Err: err}, nil
	}
}

func MakeUpdateSongEndpoint(s service.SongService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(model.UpdateSongRequest)
		if !ok {
			return nil, errors.Wrap(
				errBadRequest, "failed to cast object to UpdateSongRequest",
			)
		}
		msg, err := s.UpdateSong(ctx, &req)
		return model.UpdateSongResponse{Msg: msg, Err: err}, nil
	}
}

func MakeDeleteSongEndpoint(s service.SongService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(model.DeleteSongRequest)
		if !ok {
			return nil, errors.Wrap(
				errBadRequest, "failed to cast object to DeleteSongRequest",
			)
		}
		msg, err := s.DeleteSong(ctx, req.Name)
		return model.DeleteSongResponse{Msg: msg, Err: err}, nil
	}
}

func CreateSongHandler(service service.SongService) http.Handler {
	return http.Handler(MakeCreateSongEndpoint(service))
}

func GetSongHandler(service service.SongService) http.Handler {
	return http.Handler(MakeGetSongEndpoint(service))
}

func GetAllSongsHandler(service service.SongService) http.Handler {
	return http.Handler(MakeGetAllSongsEndpoint(service))
}

func UpdateSongHandler(service service.SongService) http.Handler {
	return http.Handler(MakeUpdateSongEndpoint(service))
}

func DeleteSongHandler(service service.SongService) http.Handler {
	return http.Handler(MakeDeleteSongEndpoint(service))
}

func NewCreateSongHandlerOption(opts []http.HandlerOption) []http.HandlerOption {
	return append([]http.HandlerOption{
		http.HandlerWithDecoder(CreateSongDecoderFunc),
		http.HandlerWithEncoder(CreateSongEncoderFunc),
		http.HandlerWithErrorEncoder(errorEncoder),
	}, opts...)
}

func NewGetSongHandlerOption(opts []http.HandlerOption) []http.HandlerOption {
	return append([]http.HandlerOption{
		http.HandlerWithDecoder(GetSongDecoderFunc),
		http.HandlerWithEncoder(GetSongEncoderFunc),
		http.HandlerWithErrorEncoder(errorEncoder),
	}, opts...)
}

func NewGetAllSongsHandlerOption(opts []http.HandlerOption) []http.HandlerOption {
	return append([]http.HandlerOption{
		http.HandlerWithDecoder(GetAllSongsDecoderFunc),
		http.HandlerWithEncoder(GetAllSongsEncoderFunc),
		http.HandlerWithErrorEncoder(errorEncoder),
	}, opts...)
}

func NewUpdateSongHandlerOption(opts []http.HandlerOption) []http.HandlerOption {
	return append([]http.HandlerOption{
		http.HandlerWithDecoder(UpdateSongDecoderFunc),
		http.HandlerWithEncoder(UpdateSongEncoderFunc),
		http.HandlerWithErrorEncoder(errorEncoder),
	}, opts...)
}

func NewDeleteSongHandlerOption(opts []http.HandlerOption) []http.HandlerOption {
	return append([]http.HandlerOption{
		http.HandlerWithDecoder(DeleteSongDecoderFunc),
		http.HandlerWithEncoder(DeleteSongEncoderFunc),
		http.HandlerWithErrorEncoder(errorEncoder),
	}, opts...)
}

func CreateSongDecoderFunc(ctx context.Context, r *net_http.Request) (interface{}, error) {
	var req model.CreateSongRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func GetSongDecoderFunc(ctx context.Context, r *net_http.Request) (interface{}, error) {
	return model.GetSongRequest{Name: http.Parameters(r).ByName("name")}, nil
}

func GetAllSongsDecoderFunc(ctx context.Context, r *net_http.Request) (interface{}, error) {
	return model.GetSongListRequest{}, nil
}

func UpdateSongDecoderFunc(ctx context.Context, r *net_http.Request) (interface{}, error) {
	var req model.UpdateSongRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	// Extract the name from path parameter
	name := http.Parameters(r).ByName("name")
	// Set the name in the request
	req.Name = name
	// Set the name in the song if not already set
	if req.Song.Name == "" {
		req.Song.Name = name
	}
	return req, nil
}

func DeleteSongDecoderFunc(ctx context.Context, r *net_http.Request) (interface{}, error) {
	return model.DeleteSongRequest{Name: http.Parameters(r).ByName("name")}, nil
}

func CreateSongEncoderFunc(ctx context.Context, w net_http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func GetSongEncoderFunc(ctx context.Context, w net_http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func GetAllSongsEncoderFunc(ctx context.Context, w net_http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func UpdateSongEncoderFunc(ctx context.Context, w net_http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func DeleteSongEncoderFunc(ctx context.Context, w net_http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

// Error encoder
func errorEncoder(ctx context.Context, err error, w net_http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(net_http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
