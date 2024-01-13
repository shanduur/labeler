package download

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/google/go-github/v53/github"
	"github.com/shanduur/labeler/labels"
	"github.com/urfave/cli/v3"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
)

var (
	output string

	flagOutput = &cli.StringFlag{
		Name:        "output",
		Aliases:     []string{"o"},
		Destination: &output,
	}
)

func New() *cli.Command {
	return &cli.Command{
		Name:    "download",
		Aliases: []string{"d", "dl"},
		Flags: []cli.Flag{
			flagOutput,
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				return errors.New("nothing to do")
			}

			client := github.NewClient(nil)

			for i := 0; i < cmd.NArg(); i++ {
				ghPath := cmd.Args().Get(i)
				owner, repo := getOwnerRepo(ghPath)

				l, err := listAll(ctx, client, owner, repo)
				if err != nil {
					return fmt.Errorf("unable to list all: %w", err)
				}

				b, err := toYAML(l)
				if err != nil {
					return fmt.Errorf("converision failed: %w", err)
				}

				outputPath := ghPath
				if output != "" {
					outputPath = path.Join(output, ghPath)
				}

				err = save(outputPath, b)
				if err != nil {
					return fmt.Errorf("saving failed: %w", err)
				}
			}

			return nil
		},
	}
}

func getOwnerRepo(arg string) (string, string) {
	args := strings.SplitN(arg, "/", 2)
	if len(args) == 2 {
		return args[0], args[1]
	}
	return "", ""
}

func listAll(ctx context.Context, client *github.Client, owner, repo string) ([]*github.Label, error) {
	slog.Debug("listing", "owner", owner, "repo", repo)
	var (
		page      = 0
		allLabels []*github.Label
	)

	if client == nil {
		return nil, errors.New("client is nil")
	}

	for {
		labels, res, err := client.Issues.ListLabels(ctx, owner, repo, &github.ListOptions{
			Page: page,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to list labels: %w", err)
		}

		allLabels = append(allLabels, labels...)

		if res.NextPage == 0 {
			break
		}
		page = res.NextPage
	}

	slog.Debug("listing done", "labels_count", len(allLabels))

	return allLabels, nil
}

func toYAML(ghl []*github.Label) ([]byte, error) {
	slog.Debug("transforming to YAML")

	lbl := make(labels.LabelsMap)

	for _, l := range ghl {
		name := l.GetName()
		if _, ok := lbl[name]; ok {
			slog.Warn("duplicate found", "name", name)
			continue
		}

		slog.Debug("new label", "name", name)

		lbl[name] = labels.Label{
			Name:        l.GetName(),
			Color:       l.GetColor(),
			Description: l.Description,
		}
	}

	slog.Debug("transforming to YAML complete")

	out, err := yaml.Marshal(lbl.ToSlice())
	if err != nil {
		return nil, fmt.Errorf("unable to marshall YAML: %w", err)
	}

	return out, nil
}

func save(location string, data []byte) error {
	err := os.MkdirAll(location, 0o777)
	if err != nil {
		return fmt.Errorf("unable to create directory: %w", err)
	}

	output := path.Join(location, "labels.yaml")

	slog.Debug("saving to file", "location", output)

	err = os.WriteFile(output, data, 0o644)
	if err != nil {
		return fmt.Errorf("unable to save result: %w", err)
	}

	return nil
}
