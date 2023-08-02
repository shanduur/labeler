package upload

import (
	"context"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/google/go-github/v53/github"
	"github.com/shanduur/labeler/labels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

func init() {
	token, ok := os.LookupEnv("LABELER_TOKEN")
	if ok {
		transport = oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)).Transport
		recorderMode = recorder.ModeRecordOnly
	}
}

var (
	recorderMode                   = recorder.ModeRecordOnce
	transport    http.RoundTripper = nil

	fixtures = "fixtures"

	testLabelName        = "test"
	testLabelColor       = "ffffff"
	testLabelDescription = "Test label. 🚧"

	testOwner = "shanduur"
	testRepo  = "labeler"
)

func TestNew(t *testing.T) {}

func TestPostLabel(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		Name          string
		Owner         string
		Repo          string
		Label         labels.Label
		ExpectedError error
	}{
		{
			Name: "happy_path",
			Label: labels.Label{
				Name:        testLabelName,
				Color:       testLabelColor,
				Description: &testLabelDescription,
			},
			Owner: testOwner,
			Repo:  testRepo,
		},
		{
			Name: "with_prefix",
			Label: labels.Label{
				Name:        testLabelName,
				Color:       "#" + testLabelColor,
				Description: &testLabelDescription,
			},
			Owner: testOwner,
			Repo:  testRepo,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			ctx := context.Background()

			r, err := recorder.NewWithOptions(&recorder.Options{
				CassetteName:  path.Join(fixtures, tc.Name),
				Mode:          recorderMode,
				RealTransport: transport,
			})
			require.NoError(t, err)

			defer r.Stop() //nolint:errcheck

			client := github.NewClient(r.GetDefaultClient())

			err = uploadLabel(ctx, client, tc.Owner, tc.Repo, tc.Label)
			assert.ErrorIs(t, err, tc.ExpectedError)
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
