package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/pkged"
	"github.com/stretchr/testify/require"
)

func TestCLISetupAndDataTransfer(t *testing.T) {
	repo, err := filepath.Abs("..")
	require.NoError(t, err)
	oldFS := pkged.FS()
	t.Cleanup(func() { pkged.SetFS(oldFS) })
	pkged.SetFS(os.DirFS(repo))
	t.Setenv("ATK_LOCALE", "en")

	dir := t.TempDir()
	t.Chdir(dir)
	configPath := filepath.Join(dir, "artalk.yml")
	previousArgs := os.Args
	t.Cleanup(func() { os.Args = previousArgs })

	run := func(check func(*ArtalkCmd), args ...string) {
		t.Helper()
		os.Args = append([]string{"artalk", "-w", dir, "-c", configPath}, args...)
		app := New()
		app.mountCommands()
		app.RootCmd.SetArgs(os.Args[1:])
		require.NoError(t, app.RootCmd.Execute())
		if app.App != nil {
			defer func() { require.NoError(t, app.App.ResetBootstrapState()) }()
		}
		if check != nil {
			check(app)
		}
	}

	t.Run("generate configuration", func(t *testing.T) {
		run(func(app *ArtalkCmd) {
			require.Nil(t, app.App, "config generation can run without a database")
		}, "gen", "config", dir)
		conf, err := config.NewFromFile(configPath)
		require.NoError(t, err)
		require.NotEmpty(t, conf.AppKey)
	})

	// Use a disposable SQLite database for the full-boot commands.
	const cfg = "app_key: cli-test\ntimezone: Local\nsite_default: CLI Site\ndb:\n  type: sqlite\n  file: cli.db\n"
	require.NoError(t, os.WriteFile(configPath, []byte(cfg), 0600))

	t.Run("create and rotate administrator password", func(t *testing.T) {
		var id uint
		admin := []string{"admin", "--name", "Operator", "--email", "operator@example.org"}
		run(func(app *ArtalkCmd) {
			user := app.Dao().FindUser("Operator", "operator@example.org")
			require.True(t, user.IsAdmin)
			require.True(t, user.CheckPassword("original-password"))
			id = user.ID
		}, append(admin, "--password", "original-password")...)
		run(func(app *ArtalkCmd) {
			user := app.Dao().FindUser("Operator", "operator@example.org")
			require.Equal(t, id, user.ID)
			require.True(t, user.CheckPassword("new-password"))
			require.False(t, user.CheckPassword("original-password"))
		}, append(admin, "--password", "new-password")...)
	})

	t.Run("import and export", func(t *testing.T) {
		input := filepath.Join(dir, "comments.artrans")
		const comments = `[{
			"id":"1","rid":"0","content":"Imported reply",
			"nick":"Reader","email":"reader@example.org",
			"page_key":"/imported","site_name":"CLI Site"
		}]`
		require.NoError(t, os.WriteFile(input, []byte(comments), 0600))
		run(func(app *ArtalkCmd) {
			require.NotZero(t, app.Dao().FindPage("/imported", "CLI Site").ID)
		}, "import", "-y", input)

		output := t.TempDir()
		run(nil, "export", output)
		files, err := filepath.Glob(filepath.Join(output, "backup-*.artrans"))
		require.NoError(t, err)
		require.Len(t, files, 1)
		body, err := os.ReadFile(files[0])
		require.NoError(t, err)
		var exported []entity.Artran
		require.NoError(t, json.Unmarshal(body, &exported))
		require.Len(t, exported, 1)
		require.Equal(t, "Imported reply", exported[0].Content)
		require.Equal(t, "/imported", exported[0].PageKey)
		require.Equal(t, "CLI Site", exported[0].SiteName)
	})
}
