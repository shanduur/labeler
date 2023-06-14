package upload

import (
	"context"
	"testing"

	"github.com/google/go-github/v53/github"
	"github.com/seborama/govcr/v13"
	"github.com/shanduur/labeler/labels"
	"github.com/stretchr/testify/assert"
)

var (
	fixtures = "test/fixtures.json"

	testLabelName        = "test"
	testLabelColor       = "ffffff"
	testLabelDescription = "Test label. 🚧"

	testOwner = "shanduur"
	testRepo  = "labelmgr"
)

func TestNew(t *testing.T) {}

func TestPostLabel(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		Name          string
		Owner         string
		Repo          string
		Label         labels.Label
		Cassette      *govcr.CassetteLoader
		ExpectedError error
	}{
		{
			Name: "happy_path",
			Label: labels.Label{
				Name:        testLabelName,
				Color:       testLabelColor,
				Description: &testLabelDescription,
			},
			Owner:    testOwner,
			Repo:     testRepo,
			Cassette: govcr.NewCassetteLoader(fixtures),
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			ctx := context.Background()

			client := github.NewClient(govcr.NewVCR(tc.Cassette).HTTPClient())

			err := uploadLabel(ctx, client, tc.Owner, tc.Repo, tc.Label)
			assert.ErrorIs(t, tc.ExpectedError, err)
		})
	}
}

func TestLabelsEqual(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		Name            string
		LabelA          *github.Label
		LabelB          *github.Label
		ExpectedOutcome bool
	}{
		{
			Name: "a_equal_b",
			LabelA: &github.Label{
				Name:        &testLabelName,
				Description: &testLabelDescription,
				Color:       &testLabelColor,
			},
			LabelB: &github.Label{
				Name:        &testLabelName,
				Description: &testLabelDescription,
				Color:       &testLabelColor,
			},
			ExpectedOutcome: true,
		},
		{
			Name: "a_equal_b_missing_color",
			LabelA: &github.Label{
				Name:        &testLabelName,
				Description: &testLabelDescription,
			},
			LabelB: &github.Label{
				Name:        &testLabelName,
				Description: &testLabelDescription,
			},
			ExpectedOutcome: true,
		},
		{
			Name: "a_equal_b_missing_description",
			LabelA: &github.Label{
				Name:  &testLabelName,
				Color: &testLabelColor,
			},
			LabelB: &github.Label{
				Name:  &testLabelName,
				Color: &testLabelColor,
			},
			ExpectedOutcome: true,
		},
		{
			Name: "a_equal_b_missing_name",
			LabelA: &github.Label{
				Description: &testLabelDescription,
				Color:       &testLabelColor,
			},
			LabelB: &github.Label{
				Description: &testLabelDescription,
				Color:       &testLabelColor,
			},
			ExpectedOutcome: true,
		},
		{
			Name:            "nil_labels",
			ExpectedOutcome: true,
		},
		{
			Name: "a_missing_name",
			LabelA: &github.Label{
				Description: &testLabelDescription,
				Color:       &testLabelColor,
			},
			LabelB: &github.Label{
				Name:        &testLabelName,
				Description: &testLabelDescription,
				Color:       &testLabelColor,
			},
			ExpectedOutcome: false,
		},
		{
			Name: "b_missing_name",
			LabelA: &github.Label{
				Name:        &testLabelName,
				Description: &testLabelDescription,
				Color:       &testLabelColor,
			},
			LabelB: &github.Label{
				Description: &testLabelDescription,
				Color:       &testLabelColor,
			},
			ExpectedOutcome: false,
		},
		{
			Name: "a_missing_description",
			LabelA: &github.Label{
				Name:  &testLabelName,
				Color: &testLabelColor,
			},
			LabelB: &github.Label{
				Name:        &testLabelName,
				Description: &testLabelDescription,
				Color:       &testLabelColor,
			},
			ExpectedOutcome: false,
		},
		{
			Name: "b_missing_description",
			LabelA: &github.Label{
				Name:        &testLabelName,
				Description: &testLabelDescription,
				Color:       &testLabelColor,
			},
			LabelB: &github.Label{
				Name:  &testLabelName,
				Color: &testLabelColor,
			},
			ExpectedOutcome: false,
		},
		{
			Name: "a_missing_color",
			LabelA: &github.Label{
				Name:        &testLabelName,
				Description: &testLabelDescription,
			},
			LabelB: &github.Label{
				Name:        &testLabelName,
				Description: &testLabelDescription,
				Color:       &testLabelColor,
			},
			ExpectedOutcome: false,
		},
		{
			Name: "b_missing_color",
			LabelA: &github.Label{
				Name:        &testLabelName,
				Description: &testLabelDescription,
				Color:       &testLabelColor,
			},
			LabelB: &github.Label{
				Name:        &testLabelName,
				Description: &testLabelDescription,
			},
			ExpectedOutcome: false,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			ok := labelsEqual(tc.LabelA, tc.LabelB)
			assert.Equal(t, tc.ExpectedOutcome, ok)
		})
	}
}
