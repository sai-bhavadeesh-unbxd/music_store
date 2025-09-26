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

func MakeCreateUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(model.CreateUserRequest)
		if !ok {
			return nil, errors.Wrap(
				errBadRequest, "failed to cast object to CreateUserRequest",
			)
		}
		msg, err := s.CreateUser(ctx, &req)
		return model.CreateUserResponse{Msg: msg, Err: err}, nil
	}
}

func MakeGetUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(model.GetUserRequest)
		if !ok {
			return nil, errors.Wrap(
				errBadRequest, "failed to cast object to GetUserRequest",
			)
		}
		user, err := s.GetUser(ctx, req.ID)
		if err != nil {
			return model.GetUserResponse{User: nil, Err: err}, nil
		}
		return model.GetUserResponse{User: user.User, Err: nil}, nil
	}
}

func MakeGetAllUsersEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		users, err := s.GetAllUsers(ctx)
		if err != nil {
			return model.GetUserListResponse{Users: nil, Err: err}, nil
		}
		return model.GetUserListResponse{Users: users.Users, Err: nil}, nil
	}
}

func MakeUpdateUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(model.UpdateUserRequest)
		if !ok {
			return nil, errors.Wrap(
				errBadRequest, "failed to cast object to UpdateUserRequest",
			)
		}
		msg, err := s.UpdateUser(ctx, &req)
		return model.UpdateUserResponse{Msg: msg, Err: err}, nil
	}
}

func MakeDeleteUserEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(model.DeleteUserRequest)
		if !ok {
			return nil, errors.Wrap(
				errBadRequest, "failed to cast object to DeleteUserRequest",
			)
		}
		msg, err := s.DeleteUser(ctx, req.ID)
		return model.DeleteUserResponse{Msg: msg, Err: err}, nil
	}
}

func MakeLikeSongEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(model.LikeSongRequest)
		if !ok {
			return nil, errors.Wrap(
				errBadRequest, "failed to cast object to LikeSongRequest",
			)
		}
		msg, err := s.LikeSong(ctx, req.UserID, req.SongName)
		return model.LikeSongResponse{Msg: msg, Err: err}, nil
	}
}

func MakeUnlikeSongEndpoint(s service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(model.UnlikeSongRequest)
		if !ok {
			return nil, errors.Wrap(
				errBadRequest, "failed to cast object to UnlikeSongRequest",
			)
		}
		msg, err := s.UnlikeSong(ctx, req.UserID, req.SongName)
		return model.UnlikeSongResponse{Msg: msg, Err: err}, nil
	}
}

func MakeGetLikedSongsEndpoint(s service.UserService) endpoint.Endpoint {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(model.GetLikedSongsRequest)
		if !ok {
			return nil, errors.Wrap(
				errBadRequest, "failed to cast object to GetLikedSongsRequest",
			)
		}
		likedSongs, err := s.GetLikedSongs(ctx, req.UserID)
		return model.GetLikedSongsResponse{LikedSongs: likedSongs, Err: err}, nil
	}
}

func CreateUserHandler(service service.UserService) http.Handler {
	return http.Handler(MakeCreateUserEndpoint(service))
}

func GetUserHandler(service service.UserService) http.Handler {
	return http.Handler(MakeGetUserEndpoint(service))
}

func GetAllUsersHandler(service service.UserService) http.Handler {
	return http.Handler(MakeGetAllUsersEndpoint(service))
}

func UpdateUserHandler(service service.UserService) http.Handler {
	return http.Handler(MakeUpdateUserEndpoint(service))
}

func DeleteUserHandler(service service.UserService) http.Handler {
	return http.Handler(MakeDeleteUserEndpoint(service))
}

func LikeSongHandler(service service.UserService) http.Handler {
	return http.Handler(MakeLikeSongEndpoint(service))
}

func UnlikeSongHandler(service service.UserService) http.Handler {
	return http.Handler(MakeUnlikeSongEndpoint(service))
}

func GetLikedSongsHandler(service service.UserService) http.Handler {
	return http.Handler(MakeGetLikedSongsEndpoint(service))
}

func NewCreateUserHandlerOption(opts []http.HandlerOption) []http.HandlerOption {
	return append([]http.HandlerOption{
		http.HandlerWithDecoder(CreateUserDecoderFunc),
		http.HandlerWithEncoder(CreateUserEncoderFunc),
		http.HandlerWithErrorEncoder(errorEncoder),
	}, opts...)
}

func NewGetUserHandlerOption(opts []http.HandlerOption) []http.HandlerOption {
	return append([]http.HandlerOption{
		http.HandlerWithDecoder(GetUserDecoderFunc),
		http.HandlerWithEncoder(GetUserEncoderFunc),
		http.HandlerWithErrorEncoder(errorEncoder),
	}, opts...)
}

func NewGetAllUsersHandlerOption(opts []http.HandlerOption) []http.HandlerOption {
	return append([]http.HandlerOption{
		http.HandlerWithDecoder(GetAllUsersDecoderFunc),
		http.HandlerWithEncoder(GetAllUsersEncoderFunc),
		http.HandlerWithErrorEncoder(errorEncoder),
	}, opts...)
}

func NewLikeSongHandlerOption(opts []http.HandlerOption) []http.HandlerOption {
	return append([]http.HandlerOption{
		http.HandlerWithDecoder(LikeSongDecoderFunc),
		http.HandlerWithEncoder(LikeSongEncoderFunc),
		http.HandlerWithErrorEncoder(errorEncoder),
	}, opts...)
}

func NewUnlikeSongHandlerOption(opts []http.HandlerOption) []http.HandlerOption {
	return append([]http.HandlerOption{
		http.HandlerWithDecoder(UnlikeSongDecoderFunc),
		http.HandlerWithEncoder(UnlikeSongEncoderFunc),
		http.HandlerWithErrorEncoder(errorEncoder),
	}, opts...)
}

func NewGetLikedSongsHandlerOption(opts []http.HandlerOption) []http.HandlerOption {
	return append([]http.HandlerOption{
		http.HandlerWithDecoder(GetLikedSongsDecoderFunc),
		http.HandlerWithEncoder(GetLikedSongsEncoderFunc),
		http.HandlerWithErrorEncoder(errorEncoder),
	}, opts...)
}

func NewUpdateUserHandlerOption(opts []http.HandlerOption) []http.HandlerOption {
	return append([]http.HandlerOption{
		http.HandlerWithDecoder(UpdateUserDecoderFunc),
		http.HandlerWithEncoder(UpdateUserEncoderFunc),
		http.HandlerWithErrorEncoder(errorEncoder),
	}, opts...)
}

func NewDeleteUserHandlerOption(opts []http.HandlerOption) []http.HandlerOption {
	return append([]http.HandlerOption{
		http.HandlerWithDecoder(DeleteUserDecoderFunc),
		http.HandlerWithEncoder(DeleteUserEncoderFunc),
		http.HandlerWithErrorEncoder(errorEncoder),
	}, opts...)
}

// Decoder functions
func CreateUserDecoderFunc(ctx context.Context, r *net_http.Request) (interface{}, error) {
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func GetUserDecoderFunc(ctx context.Context, r *net_http.Request) (interface{}, error) {
	id := http.Parameters(r).ByName("id")
	return model.GetUserRequest{ID: id}, nil
}

func GetAllUsersDecoderFunc(ctx context.Context, r *net_http.Request) (interface{}, error) {
	return model.GetUserListRequest{}, nil
}

func UpdateUserDecoderFunc(ctx context.Context, r *net_http.Request) (interface{}, error) {
	var req model.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	// Get ID from path parameter
	req.ID = http.Parameters(r).ByName("id")
	return req, nil
}

func DeleteUserDecoderFunc(ctx context.Context, r *net_http.Request) (interface{}, error) {
	id := http.Parameters(r).ByName("id")
	return model.DeleteUserRequest{ID: id}, nil
}

func LikeSongDecoderFunc(ctx context.Context, r *net_http.Request) (interface{}, error) {
	return model.LikeSongRequest{UserID: http.Parameters(r).ByName("id"), SongName: http.Parameters(r).ByName("song_name")}, nil
}

func UnlikeSongDecoderFunc(ctx context.Context, r *net_http.Request) (interface{}, error) {
	return model.UnlikeSongRequest{UserID: http.Parameters(r).ByName("id"), SongName: http.Parameters(r).ByName("song_name")}, nil
}

func GetLikedSongsDecoderFunc(ctx context.Context, r *net_http.Request) (interface{}, error) {
	return model.GetLikedSongsRequest{UserID: http.Parameters(r).ByName("id")}, nil
}

// Encoder functions
func CreateUserEncoderFunc(ctx context.Context, w net_http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func GetUserEncoderFunc(ctx context.Context, w net_http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func GetAllUsersEncoderFunc(ctx context.Context, w net_http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func UpdateUserEncoderFunc(ctx context.Context, w net_http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func LikeSongEncoderFunc(ctx context.Context, w net_http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func UnlikeSongEncoderFunc(ctx context.Context, w net_http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func GetLikedSongsEncoderFunc(ctx context.Context, w net_http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func DeleteUserEncoderFunc(ctx context.Context, w net_http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
