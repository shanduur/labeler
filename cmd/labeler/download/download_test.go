package download

import (
	"context"
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/google/go-github/v53/github"
	"github.com/seborama/govcr/v13"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fixtures = "test/fixtures.json"

func TestNew(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		Name string
	}{
		{
			Name: "happy_path",
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			cli := New()
			assert.NotNil(t, cli)
		})
	}
}

func TestGetOwnerRepo(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		Name          string
		Input         string
		ExpectedOwner string
		ExpectedRepo  string
	}{
		{
			Name:          "happy_path",
			Input:         "shanduur/labeler",
			ExpectedOwner: "shanduur",
			ExpectedRepo:  "labeler",
		},
		{
			Name:          "missing_owner",
			Input:         "/labeler",
			ExpectedOwner: "",
			ExpectedRepo:  "labeler",
		},
		{
			Name:          "missing_repo",
			Input:         "shanduur/",
			ExpectedOwner: "shanduur",
			ExpectedRepo:  "",
		},
		{
			Name:          "missing_owner_and_repo",
			Input:         "/",
			ExpectedOwner: "",
			ExpectedRepo:  "",
		},
		{
			Name:          "empty_input",
			Input:         "",
			ExpectedOwner: "",
			ExpectedRepo:  "",
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			owner, repo := getOwnerRepo(tc.Input)
			assert.Equal(t, tc.ExpectedOwner, owner)
			assert.Equal(t, tc.ExpectedRepo, repo)
		})
	}
}

func TestListAll(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		Name           string
		Owner          string
		Repo           string
		Cassette       *govcr.CassetteLoader
		ExpectedLabels []*github.Label
		ExpectedError  error
	}{
		{
			Name:     "happy_path",
			Owner:    "test",
			Repo:     "test",
			Cassette: govcr.NewCassetteLoader(fixtures),
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			ctx := context.Background()

			client := github.NewClient(govcr.NewVCR(tc.Cassette).HTTPClient())

			labelList, err := listAll(ctx, client, tc.Owner, tc.Repo)
			assert.ErrorIs(t, tc.ExpectedError, err)
			if tc.ExpectedError == nil {
				assert.Equal(t, tc.ExpectedLabels, labelList)
			}
		})
	}
}

func TestToYAML(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		Name string
	}{
		{
			Name: "happy_path",
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
		})
	}
}

func TestSave(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		Name          string
		Location      string
		Permissions   fs.FileMode
		Data          []byte
		ExpectedError error
	}{
		{
			Name:          "happy_path",
			Location:      "test",
			Permissions:   0o700,
			Data:          []byte("hello_world"),
			ExpectedError: nil,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			tmpDir, err := os.MkdirTemp("", tc.Name+"*")
			require.NoError(t, err)
			tmpDir = path.Join(tmpDir, tc.Name)

			os.Mkdir(tmpDir, tc.Permissions)

			defer func() {
				err := os.RemoveAll(tmpDir)
				require.NoError(t, err)
			}()

			tmpPath := path.Join(tmpDir, tc.Location)

			err = save(tmpPath, tc.Data)
			assert.ErrorIs(t, tc.ExpectedError, err)
			if tc.ExpectedError == nil {
				assert.FileExists(t, path.Join(tmpPath, "labels.yaml"))
			}
		})
	}
}
