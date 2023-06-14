package labels

import (
	"testing"

	"github.com/google/go-github/v53/github"
	"github.com/stretchr/testify/assert"
)

var (
	testName        = "test"
	testColor       = "000000"
	testDescription = "Test label. 🚧"

	testLabel = Label{
		Name:        testName,
		Color:       testColor,
		Description: &testDescription,
	}
)

func TestLabelToGithub(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		Name           string
		Label          Label
		ExpectedOutput *github.Label
	}{
		{
			Name:  "happy_path",
			Label: testLabel,
			ExpectedOutput: &github.Label{
				Name:        &testName,
				Color:       &testColor,
				Description: &testDescription,
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			tc := tc
			t.Parallel()
			gh := tc.Label.ToGitHub()
			assert.Equal(t, tc.ExpectedOutput.Name, gh.Name)
			assert.Equal(t, tc.ExpectedOutput.Color, gh.Color)
			assert.Equal(t, tc.ExpectedOutput.Description, gh.Description)
		})
	}

}
