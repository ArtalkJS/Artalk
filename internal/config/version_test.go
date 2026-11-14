package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReleaseVersionIncludesBuildCommit(t *testing.T) {
	previous := buildCommitHash
	t.Cleanup(func() { buildCommitHash = previous })
	for _, tc := range []struct{ input, want string }{{"abc1234def5678", "abc1234"}, {"abc1234", "abc1234"}, {"abc", "abc"}} {
		buildCommitHash = tc.input
		require.Equal(t, tc.want, CommitHash())
		require.Equal(t, Version+"/"+tc.want, VersionString())
	}
}
