package labels

import (
	"github.com/google/go-github/v53/github"
)

type Label struct {
	Name        string  `yaml:"name"`
	Color       string  `yaml:"color"`
	Description *string `yaml:"description,omitempty"`
}

func (l *Label) ToGitHub() *github.Label {
	return &github.Label{
		Name:        &l.Name,
		Color:       &l.Color,
		Description: l.Description,
	}
}

type Labels map[string]Label
