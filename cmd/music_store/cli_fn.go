package main

import (
	"bytes"
	"fmt"
	musicstore "music-store"
	"strings"

	"github.com/urfave/cli/v2"
)

func beforeStart(cx *cli.Context) (err error) {
	musicStore, err = musicstore.NewMusicStore()
	if err != nil {
		return cli.Exit(
			fmt.Sprintf(
				"-- \nfailed to initialize music store. \n--\n Caused By:\n%s\n--",
				errorstack(err.Error()),
			), 9,
		)
	}
	return
}

func actionStart(cx *cli.Context) (err error) {
	fmt.Println("Starting music store...")
	fmt.Println("---------------------------------------------------------------------")
	fmt.Println("Startup Flags")
	fmt.Println("---------------------------------------------------------------------")
	fmt.Println("env:								", cx.String("env"))
	fmt.Println("http.host:							", cx.String("http.host"))
	fmt.Println("http.port:							", cx.String("http.port"))
	fmt.Println("http.monitor:							", cx.StringSlice("http.monitor"))
	fmt.Println("redis.conn:							", cx.String("redis.conn"))
	fmt.Println("redis.password:							", cx.String("redis.password"))
	fmt.Println("redis.database:							", cx.Int("redis.database"))
	fmt.Println("embedding.model:							", cx.String("embedding.model"))
	fmt.Println("embedding.dim:							", cx.Int("embedding.dim"))
	fmt.Println("---------------------------------------------------------------------")
	return musicStore.Open(cx.Context)
}

func errorstack(errorstr string) string {
	parts := strings.Split(errorstr, ": ")

	var buff bytes.Buffer

	for ix, p := range parts {
		buff.WriteRune('\n')

		for i := 0; i <= ix; i++ {
			buff.WriteRune(' ')
		}

		buff.WriteString("> ")
		buff.WriteString(p)
		if ix > 3 {
			break
		}
	}

	for i := 4; i < len(parts); i++ {
		buff.WriteString(parts[i])
		buff.WriteString(": ")
	}

	return buff.String()
}
