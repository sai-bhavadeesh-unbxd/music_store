package flags

import "github.com/urfave/cli/v2"

func Flags() []cli.Flag {
	flags := []cli.Flag{}
	flags = append(flags, appflags...)
	flags = append(flags, redisflags...)
	flags = append(flags, embeddingflags...)
	flags = append(flags, httpflags...)
	return flags
}

var (
	appflags = []cli.Flag{
		&cli.StringFlag{
			Name:        "env",
			Aliases:     []string{"e"},
			Value:       "dev",
			Usage:       "set the environment in which server is running",
			DefaultText: "--env=dev",
			EnvVars:     []string{"MS_ENV"},
		},
	}
)

var (
	redisflags = []cli.Flag{
		&cli.StringFlag{
			Name:        "redis.conn",
			Value:       "localhost:6379",
			Usage:       "set the connection string for redis",
			DefaultText: "--redis.conn=localhost:6379",
			EnvVars:     []string{"MS_REDIS_CONN"},
		},
		&cli.StringFlag{
			Name:        "redis.password",
			Value:       "",
			Usage:       "set the password for redis",
			DefaultText: "--redis.password=",
			EnvVars:     []string{"MS_REDIS_PASSWORD"},
		},
		&cli.IntFlag{
			Name:        "redis.database",
			Value:       0,
			Usage:       "set the database for redis",
			DefaultText: "--redis.database=0",
			EnvVars:     []string{"MS_REDIS_DATABASE"},
		},
	}
)

var (
	embeddingflags = []cli.Flag{
		&cli.StringFlag{
			Name:        "embedding.model",
			Value:       "sentence-transformers/all-MiniLM-L6-v2",
			Usage:       "set the embedding model",
			DefaultText: "--embedding.model=sentence-transformers/all-MiniLM-L6-v2",
			EnvVars:     []string{"MS_EMBEDDING_MODEL"},
		},
		&cli.StringFlag{
			Name:        "embedding.model.path",
			Value:       "onnx/model.onnx",
			Usage:       "set the embedding model path",
			DefaultText: "--embedding.model.path=onnx/model.onnx",
			EnvVars:     []string{"MS_EMBEDDING_MODEL_PATH"},
		},
		&cli.StringFlag{
			Name:        "embedding.models.dir",
			Value:       "./models/",
			Usage:       "set the embedding models directory",
			DefaultText: "--embedding.models.dir=./models/",
			EnvVars:     []string{"MS_EMBEDDING_MODELS_DIR"},
		},
		&cli.IntFlag{
			Name:        "embedding.dim",
			Value:       512,
			Usage:       "set the embedding model dimension",
			DefaultText: "--embedding.dim=512",
			EnvVars:     []string{"MS_EMBEDDING_DIM"},
		},
	}
)

var (
	httpflags = []cli.Flag{
		&cli.StringFlag{
			Name:        "http.host",
			Value:       "localhost",
			Usage:       "set host listener location for http server.",
			DefaultText: "localhost",
			EnvVars:     []string{"MS_HTTP_HOST"},
		},
		&cli.StringFlag{
			Name:        "http.port",
			Value:       "8080",
			Usage:       "set port number for http listener",
			DefaultText: "8080",
			EnvVars:     []string{"MS_HTTP_PORT"},
		},
	}
)
