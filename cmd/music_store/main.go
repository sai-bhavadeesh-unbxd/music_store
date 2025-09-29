package main

import (
	"fmt"
	"log"
	flags "music-store/cmd/flags"
	build "music-store/cmd/ld"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	err := (&cli.App{
		Name:     "Music Store",
		Usage:    "Music Store API",
		Version:  build.Version,
		Flags:    flags.Flags(),
		Commands: commands,
	}).Run(os.Args)
	if err != nil {
		fmt.Println("Something Went Wrong. Failed to start Music Store.: " + err.Error())
		log.Fatalf(
			"%s", fmt.Sprintf(
				"-- \nfailed to start music store. \n--\n Caused By:\n%s\n--",
				errorstack(err.Error()),
			),
		)
	}
}
