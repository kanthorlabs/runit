package main

import (
	_ "embed"
	"log"
	"os"
	"path/filepath"

	"github.com/kanthorlabs/runit/platform/dockerx"
	"github.com/kanthorlabs/runit/runtime/pythonx"
	"github.com/urfave/cli/v2"
)

//go:embed .version
var version string

func main() {
	app := &cli.App{
		Name:    "runit",
		Usage:   "Run arbitrary python script you have, no setup, no configuration, just run it",
		Version: version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "platform-version",
				Value: pythonx.DefaultPythonVersion,
				Usage: "Python version to use (e.g. python:3.13-slim)",
			},
			&cli.StringSliceFlag{
				Name:  "ports",
				Usage: "Ports to expose (can be specified multiple times)",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return ErrScriptPathRequired
			}

			scriptPath := c.Args().First()
			absPath, err := filepath.Abs(scriptPath)
			if err != nil {
				return ErrInvalidScriptPath
			}
			if _, err := os.Stat(absPath); err != nil {
				return ErrScriptNotFound
			}

			vars := &pythonx.DockerfileVars{
				Version: c.String("platform-version"),
				Ports:   c.StringSlice("ports"),
			}

			return dockerx.Exec(scriptPath, vars)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
