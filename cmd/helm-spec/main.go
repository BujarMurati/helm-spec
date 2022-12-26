package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

const (
	errFailedToGetAbsolutePath = "failed to get absolute path of: %v"
	errNotADirectory           = "%v is not a directory"
	errNoSpecFilesFound        = "no %v files found in %v"
	specFileGlobPattern        = "*_spec.yaml"
)

type cliSettings struct {
	Reader         io.Reader
	Writer         io.Writer
	ErrWriter      io.Writer
	ExitErrHandler cli.ExitErrHandlerFunc
}

var defaultSettings = cliSettings{
	Reader:         os.Stdin,
	Writer:         os.Stdout,
	ErrWriter:      os.Stderr,
	ExitErrHandler: nil,
}

// validates that the path is an existing directory containing spec files
func validateTestSuitePath(cCtx *cli.Context, value string) (err error) {
	absPath, err := filepath.Abs(value)
	if err != nil {
		return fmt.Errorf(errFailedToGetAbsolutePath, value)
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf(errNotADirectory, absPath)
	}
	specFiles, err := filepath.Glob(filepath.Join(absPath, specFileGlobPattern))
	if err != nil {
		return err
	}
	if len(specFiles) == 0 {
		return fmt.Errorf(errNoSpecFilesFound, specFileGlobPattern, absPath)
	}
	return nil
}

func createApp(settings cliSettings) (app *cli.App, err error) {
	app = &cli.App{
		Name:           "helm-spec",
		DefaultCommand: "test",
		Commands: []*cli.Command{
			{
				Name: "test",
				Flags: []cli.Flag{
					&cli.PathFlag{
						Usage:       "path to a directory containing *_spec.yaml files",
						Name:        "testsuite",
						Aliases:     []string{"t"},
						Value:       "./specs",
						DefaultText: "./specs",
						Required:    false,
						TakesFile:   false,
						Action:      validateTestSuitePath,
					},
				},
			},
		},
		Reader:         settings.Reader,
		Writer:         settings.Writer,
		ErrWriter:      settings.ErrWriter,
		ExitErrHandler: settings.ExitErrHandler,
	}
	return app, nil
}

func main() {
	app, err := createApp(defaultSettings)
	if err != nil {
		panic("failed to initialize app")
	}
	app.Run(os.Args)
}
