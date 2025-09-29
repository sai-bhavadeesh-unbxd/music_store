package main

import "github.com/urfave/cli/v2"

var (
	start = &cli.Command{
		Name:    "start",
		Aliases: []string{"s"},
		Usage:   "Start the music store",
		Before:  beforeStart,
		Action:  actionStart,
	}

	version = &cli.Command{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "prints version details of the binary",
	}

	debug = &cli.Command{
		Name:    "debug",
		Aliases: []string{"d"},
		Usage:   "prints debug information about the binary",
	}

	commands = []*cli.Command{start, version, debug}
)
