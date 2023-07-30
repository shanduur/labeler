package labels

import (
	"errors"
	"strings"

	"github.com/google/go-github/v53/github"
)

type Label struct {
	Name        string  `yaml:"name"`
	Color       string  `yaml:"color"`
	Description *string `yaml:"description,omitempty"`
}

func (l *Label) Validate() error {
	if l.Name == "" {
		return errors.New("name is empty")
	}

	if len(l.Color) < 3 || len(l.Color) > 7 {
		return errors.New("color is empty")
	}

	return nil
}

func (l *Label) ToGitHub() *github.Label {
	color := l.Color
	if strings.HasPrefix(l.Color, "#") {
		color = strings.TrimPrefix(color, "#")
	}

	return &github.Label{
		Name:        &l.Name,
		Color:       &color,
		Description: l.Description,
	}
}

type Labels map[string]Label
