package files

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGitHash(t *testing.T) {
	t.Run("known hash", func(t *testing.T) {
		require.Equal(t,
			"257cc5642cb1a054f08cc83f2d943e56fd3ebe99",
			GitHash([]byte("foo\n")),
		)
	})
}
