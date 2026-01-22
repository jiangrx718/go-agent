package server

import "github.com/urfave/cli/v2"

func Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "config",
			Value:    "",
			Usage:    "specify the location of the configuration file",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "env",
			Value:    "",
			Usage:    "specify the location of the configuration file",
			Required: false,
		},
	}
}
