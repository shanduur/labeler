package main

import (
	"context"
	"os"

	"github.com/doomshrine/labeler/cmd/labeler/download"
	"github.com/doomshrine/labeler/cmd/labeler/upload"
	"github.com/urfave/cli/v3"
	"golang.org/x/exp/slog"
)

var app = &cli.App{
	Commands: []*cli.Command{
		upload.New(),
		download.New(),
	},
}

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
}

func main() {
	err := app.RunContext(context.TODO(), os.Args)

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
