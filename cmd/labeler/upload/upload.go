package upload

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/v53/github"
	"github.com/shanduur/labeler/labels"
	"github.com/urfave/cli/v3"
	"golang.org/x/exp/slog"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
)

var (
	owner string
	repo  string

	flagOwner = &cli.StringFlag{
		Name:        "owner",
		Aliases:     []string{"o"},
		Destination: &owner,
	}
	flagRepo = &cli.StringFlag{
		Name:        "repo",
		Aliases:     []string{"r"},
		Destination: &repo,
	}
)

func New() *cli.Command {
	return &cli.Command{
		Name:      "upload",
		Aliases:   []string{"u", "up"},
		ArgsUsage: "<file>",
		Flags: []cli.Flag{
			flagOwner,
			flagRepo,
		},
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() != 1 {
				return errors.New("wrong number of arguments, expected single <file>")
			}

			token, ok := os.LookupEnv("LABELER_TOKEN")
			if !ok {
				return errors.New("GitHub Token (env = LABELER_TOKEN) not found")
			}

			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			)

			tc := oauth2.NewClient(ctx.Context, ts)

			client := github.NewClient(tc)

			f, err := os.Open(ctx.Args().First())
			if err != nil {
				return fmt.Errorf("unable to open file: %w", err)
			}

			lbls := make(labels.Labels)

			err = yaml.NewDecoder(f).Decode(&lbls)
			if err != nil {
				return fmt.Errorf("unable to decode file: %w", err)
			}

			for name, label := range lbls {
				slog.Info("processing label", "label_name", name, "label.color", label.Color)

				err = uploadLabel(ctx.Context, client, owner, repo, label)
				if err != nil {
					return fmt.Errorf("unable to update label: %w", err)
				}
			}

			return nil
		},
	}
}

func uploadLabel(ctx context.Context, client *github.Client, owner, repo string, label labels.Label) error {
	slog.Debug("listing", "label_name", label.Name, "label", label)

	if client == nil {
		return errors.New("client is nil")
	}

	err := label.Validate()
	if err != nil {
		return fmt.Errorf("lablel validation failed: %w", err)
	}

	// first, check if label exist
	ghLabel, err := getLabel(ctx, client, owner, repo, label.Name)
	if err != nil {
		return fmt.Errorf("unable to get label: %w", err)
	}

	if ghLabel != nil {
		// then, compare, if it has valid fields
		if labelsEqual(label.ToGitHub(), ghLabel) {
			return nil
		}

		err := updateLabel(ctx, client, owner, repo, label.ToGitHub())
		if err != nil {
			return fmt.Errorf("unable to update label: %w", err)
		}

		return nil
	} else {
		err := createLabel(ctx, client, owner, repo, label.ToGitHub())
		if err != nil {
			return fmt.Errorf("unable to create label: %w", err)
		}

		return nil
	}
}

func labelsEqual(a, b *github.Label) bool {
	if a == nil && b == nil {
		return true
	} else if (a != nil && b == nil) || (a == nil && b != nil) {
		return false
	}

	if (a.Name != nil && b.Name == nil) || (a.Name == nil && b.Name != nil) {
		return false
	} else if a.Name == nil && b.Name == nil {
		// noop
	} else {
		if *a.Name != *b.Name {
			return false
		}
	}

	if (a.Color != nil && b.Color == nil) || (a.Color == nil && b.Color != nil) {
		return false
	} else if a.Color == nil && b.Color == nil {
		// noop
	} else {
		if *a.Color != *b.Color {
			return false
		}
	}

	if (a.Description != nil && b.Description == nil) || (a.Description == nil && b.Description != nil) {
		return false
	} else if a.Description == nil && b.Description == nil {
		// noop
	} else {
		if *a.Description != *b.Description {
			return false
		}
	}

	return true
}

func getLabel(ctx context.Context, client *github.Client, owner, repo, name string) (*github.Label, error) {
	slog.Debug("checking if label exists", "owner", owner, "repo", repo, "label_name", name)

	if client == nil {
		return nil, errors.New("client is nil")
	}

	label, res, err := client.Issues.GetLabel(ctx, owner, repo, name)
	if err != nil && res.StatusCode == http.StatusNotFound {
		slog.Debug("label not found", "label_name", name)
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("unable to get label: %w", err)
	}

	return label, nil
}

func updateLabel(ctx context.Context, client *github.Client, owner, repo string, label *github.Label) error {
	slog.Debug("updating label", "owner", owner, "repo", repo, "label_name", *label.Name)

	if client == nil {
		return errors.New("client is nil")
	}

	_, _, err := client.Issues.EditLabel(ctx, owner, repo, *label.Name, label)
	if err != nil {
		return fmt.Errorf("unable to edit the label: %w", err)
	}

	return nil
}

func createLabel(ctx context.Context, client *github.Client, owner, repo string, label *github.Label) error {
	slog.Debug("creating label", "owner", owner, "repo", repo, "label_name", *label.Name)

	if client == nil {
		return errors.New("client is nil")
	}

	_, _, err := client.Issues.CreateLabel(ctx, owner, repo, label)
	if err != nil {
		return fmt.Errorf("unable to create the label: %w", err)
	}

	return nil
}
