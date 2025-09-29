package musicstore

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"runtime"

	"music-store/internal/controller"
	"music-store/internal/repository"
	"music-store/internal/service"
	"music-store/utils"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/unbxd/go-base/kit/transport/http"
	"github.com/unbxd/go-base/utils/log"
)

type (
	Option func(*MusicStore) error

	MusicStore struct {
		httpTransport  *http.Transport
		redisClient    *redis.Client
		logger         log.Logger
		userController *controller.UserController
		songController *controller.SongController
	}
)

func (m *MusicStore) listen(transport *http.Transport, errch chan error) {
	err := transport.Open()
	if err != nil {
		errch <- errors.Wrap(err, "failed to start transport")
	}
}

func (m *MusicStore) Close(cx context.Context) error {
	m.redisClient.Close()
	m.logger.Flush()
	return nil
}

func (m *MusicStore) Open(cx context.Context) error {
	m.logger.Info("-- Starting Music Store!")

	var (
		intchan = make(chan os.Signal, 1)
		errchan = make(chan error)
	)

	go m.listen(m.httpTransport, errchan)
	go signal.Notify(intchan, os.Interrupt)

	for {
		select {
		case <-intchan:
			m.logger.Info("Recieved os.Interrupt. Signal: ",
				log.String("signal", os.Interrupt.String()))
			m.logger.Info("Shutting down gracefully!!!")
			err := m.httpTransport.Close()
			if err != nil {
				panic(err)
			}
			m.logger.Info("Done!")
			return nil

		case err := <-errchan:
			m.logger.Error("Failed to start music store:", log.Error(err))
			return err
		}
	}
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func NewMusicStore(options ...Option) (*MusicStore, error) {
	utils.InitRedisWithDefaults()
	var (
		tr, _       = http.NewTransport("localhost", "8080")
		lg, _       = log.NewZapLogger()
		redisClient = utils.GetRedisClient()
		userRepo    = repository.NewUserRepository(redisClient)
		songRepo    = repository.NewSongRepository(redisClient)
		vectorRepo  = repository.NewVectorRepository(redisClient)
	)

	userService := service.NewUserService(userRepo, vectorRepo)
	songService := service.NewSongService(songRepo, vectorRepo)

	o := &MusicStore{
		httpTransport:  tr,
		logger:         lg,
		redisClient:    redisClient,
		userController: controller.NewUserController(userService),
		songController: controller.NewSongController(songService),
	}

	o.userController.Bind(tr, []http.HandlerOption{})
	o.songController.Bind(tr, []http.HandlerOption{})

	for _, ofn := range options {
		fmt.Println(">> ----- Initializing: ", getFunctionName(ofn))
		err := ofn(o)
		if err != nil {
			return nil, err
		}
	}

	return o, nil
}
