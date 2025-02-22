package main

import (
	_ "embed"
	"log"
	"os"
	"path/filepath"

	"github.com/kanthorlabs/runit/platform/dockerx"
	"github.com/kanthorlabs/runit/runtime/pythonx"
	"github.com/spf13/cobra"
)

//go:embed .version
var version string

func main() {
	cmd := &cobra.Command{
		Use:     "runit [flags] <script>",
		Short:   "Run arbitrary python script you have",
		Version: version,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			scriptPath := args[0]
			if scriptPath == "" {
				return ErrScriptPathRequired()
			}

			absPath, err := filepath.Abs(scriptPath)
			if err != nil {
				return ErrInvalidScriptPath()
			}
			if _, err := os.Stat(absPath); err != nil {
				return ErrScriptNotFound()
			}

			version, err := cmd.Flags().GetString("platform-version")
			if err != nil {
				return ErrArgsParams("platform-version")
			}
			ports, err := cmd.Flags().GetStringSlice("ports")
			if err != nil {
				return ErrArgsParams("ports")
			}
			arguments, err := cmd.Flags().GetString("arguments")
			if err != nil {
				return ErrArgsParams("arguments")
			}
			params, err := cmd.Flags().GetString("params")
			if err != nil {
				return ErrArgsParams("params")
			}
			vars := &pythonx.DockerfileVars{
				Version:   version,
				Ports:     ports,
				Arguments: arguments,
				Params:    params,
			}

			return dockerx.Exec(scriptPath, vars)
		},
	}

	// Add flags
	cmd.Flags().String("platform-version", pythonx.DefaultPythonVersion, "Python version to use (e.g. python:3.13-slim)")
	cmd.Flags().StringSlice("ports", []string{}, "Ports to expose (can be specified multiple times)")
	cmd.Flags().String("arguments", "", "Main script arguments (e.g. --arguments=\"kanthorlabs/runit\")")
	cmd.Flags().String("params", "", "Additional script parameters (e.g. --params=\"--token=xxx\")")

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
